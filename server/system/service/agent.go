package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/errors"

	automationService "github.com/cortezaproject/corteza/server/automation/service"
	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/agentic/tcl"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	agent struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        agentAccessController
	}

	agentAccessController interface {
		CanCreateAgent(ctx context.Context) bool
		CanSearchAgents(ctx context.Context) bool
		CanReadAgent(ctx context.Context, a *types.Agent) bool
		CanUpdateAgent(ctx context.Context, a *types.Agent) bool
		CanDeleteAgent(ctx context.Context, a *types.Agent) bool
	}
)

func Agent() *agent {
	return &agent{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *agent) Get(ctx context.Context, ID uint64) (*types.Agent, error) {
	return svc.FindByID(ctx, ID)
}

func (svc *agent) FindByID(ctx context.Context, ID uint64) (a *types.Agent, err error) {
	err = func() error {
		if a, err = loadAgent(ctx, svc.store, ID); err != nil {
			return err
		}

		if !svc.ac.CanReadAgent(ctx, a) {
			return AgentErrNotAllowedToRead()
		}

		return nil
	}()

	return a, err
}

func (svc *agent) Create(ctx context.Context, new *types.Agent) (a *types.Agent, err error) {
	err = func() (err error) {
		if !svc.ac.CanCreateAgent(ctx) {
			return AgentErrNotAllowedToCreate()
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.Revision = 1

		if new.Status == "" {
			new.Status = "active"
		}

		prepareTCL(&new.Behavior)

		if err = store.CreateAgent(ctx, svc.store, new); err != nil {
			return
		}

		// Auto-create a linked TriggerDefinition using service-user context
		svcCtx := intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())
		_, _ = automationService.DefaultTriggerDefinition.Create(svcCtx, agentToTriggerDef(new))

		a = new
		return nil
	}()

	return a, err
}

func (svc *agent) Update(ctx context.Context, upd *types.Agent) (a *types.Agent, err error) {
	err = func() (err error) {
		if !svc.ac.CanUpdateAgent(ctx, upd) {
			return AgentErrNotAllowedToUpdate()
		}

		var existing *types.Agent
		if existing, err = store.LookupAgentByID(ctx, svc.store, upd.ID); err != nil {
			return AgentErrNotFound()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, existing.UpdatedAt, existing.CreatedAt) {
			return AgentErrStaleData()
		}

		upd.Revision = existing.Revision + 1
		upd.UpdatedAt = now()
		upd.CreatedAt = existing.CreatedAt
		upd.DeletedAt = existing.DeletedAt

		prepareTCL(&upd.Behavior)

		if err = store.UpdateAgent(ctx, svc.store, upd); err != nil {
			return
		}

		// Sync linked TriggerDefinition meta using service-user context
		svcCtx := intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())
		if td, tdErr := automationService.DefaultTriggerDefinition.LookupByHandle(svcCtx, "agent-"+upd.Handle); tdErr == nil && td != nil {
			td.Meta = &automationTypes.TriggerDefinitionMeta{
				Short:       upd.Meta.Short,
				Description: upd.Meta.Description,
			}
			_, _ = automationService.DefaultTriggerDefinition.Update(svcCtx, td)
		}

		a = upd
		return nil
	}()

	return a, err
}

func (svc *agent) DeleteByID(ctx context.Context, ID uint64) (err error) {
	err = func() (err error) {
		var a *types.Agent
		if a, err = loadAgent(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteAgent(ctx, a) {
			return AgentErrNotAllowedToDelete()
		}

		a.DeletedAt = now()
		if err = store.UpdateAgent(ctx, svc.store, a); err != nil {
			return
		}

		// Soft-delete linked TriggerDefinition using service-user context
		svcCtx := intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())
		if td, tdErr := automationService.DefaultTriggerDefinition.LookupByHandle(svcCtx, "agent-"+a.Handle); tdErr == nil && td != nil {
			_ = automationService.DefaultTriggerDefinition.DeleteByID(svcCtx, td.ID)
		}

		return nil
	}()

	return err
}

func (svc *agent) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	err = func() (err error) {
		var a *types.Agent
		if a, err = loadAgent(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteAgent(ctx, a) {
			return AgentErrNotAllowedToDelete()
		}

		a.DeletedAt = nil
		if err = store.UpdateAgent(ctx, svc.store, a); err != nil {
			return
		}

		// Restore linked TriggerDefinition using service-user context
		svcCtx := intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())
		if td, tdErr := automationService.DefaultTriggerDefinition.LookupByHandle(svcCtx, "agent-"+a.Handle); tdErr == nil && td != nil {
			_ = automationService.DefaultTriggerDefinition.UndeleteByID(svcCtx, td.ID)
		}

		return nil
	}()

	return err
}

func (svc *agent) Search(ctx context.Context, filter types.AgentFilter) (set types.AgentSet, f types.AgentFilter, err error) {
	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Agent) (bool, error) {
		if !svc.ac.CanReadAgent(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchAgents(ctx) {
			return AgentErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchAgents(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, err
}


// prepareTCL ensures TCL is properly initialized on the agent behavior:
// - if enabled and no articles selected, populate defaults
// - always merge hardwired articles back in (user cannot remove them)
func prepareTCL(b *types.AgentBehavior) {
	if b.TreatyCLEnabled != nil && !*b.TreatyCLEnabled {
		return
	}
	if len(b.TreatyCLArticles) == 0 {
		b.TreatyCLArticles = tcl.DefaultArticleIDs()
	} else {
		b.TreatyCLArticles = tcl.MergeWithHardwired(b.TreatyCLArticles)
	}
}

// agentToTriggerDef builds a TriggerDefinition linked to the given agent.
func agentToTriggerDef(a *types.Agent) *automationTypes.TriggerDefinition {
	return &automationTypes.TriggerDefinition{
		AgentID:      a.ID,
		Handle:       "agent-" + a.Handle,
		SkipEventBus: true,
		Meta: &automationTypes.TriggerDefinitionMeta{
			Short:       a.Meta.Short,
			Description: a.Meta.Description,
		},
	}
}

func loadAgent(ctx context.Context, s store.Agents, ID uint64) (res *types.Agent, err error) {
	if ID == 0 {
		return nil, AgentErrInvalidID()
	}

	if res, err = store.LookupAgentByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, AgentErrNotFound()
	}

	return
}
