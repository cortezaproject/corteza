package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/errors"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/label"

	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/agentic/tcl"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	agent struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        agentAccessController
		llm       agentLLMValidator
	}

	agentAccessController interface {
		CanCreateAgent(ctx context.Context) bool
		CanSearchAgents(ctx context.Context) bool
		CanReadAgent(ctx context.Context, a *types.Agent) bool
		CanUpdateAgent(ctx context.Context, a *types.Agent) bool
		CanDeleteAgent(ctx context.Context, a *types.Agent) bool
	}

	agentLLMValidator interface {
		ValidateTemperature(ctx context.Context, providerID uint64, model string, temperature *float64) error
	}
)

func Agent() *agent {
	return &agent{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *agent) WithLLMValidator(v agentLLMValidator) *agent {
	svc.llm = v
	return svc
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

		if err = label.Load(ctx, svc.store, a); err != nil {
			return err
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

		if new.Execution.Model.Temperature != nil && svc.llm != nil {
			if err = svc.llm.ValidateTemperature(ctx, new.Execution.Model.LLMProviderID, new.Execution.Model.Model, new.Execution.Model.Temperature); err != nil {
				return
			}
		}

		prepareTCL(&new.Behavior)

		if err = store.CreateAgent(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

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

		if upd.Execution.Model.Temperature != nil && svc.llm != nil {
			if err = svc.llm.ValidateTemperature(ctx, upd.Execution.Model.LLMProviderID, upd.Execution.Model.Model, upd.Execution.Model.Temperature); err != nil {
				return
			}
		}

		prepareTCL(&upd.Behavior)

		if err = store.UpdateAgent(ctx, svc.store, upd); err != nil {
			return
		}

		if err = label.Update(ctx, svc.store, upd); err != nil {
			return
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

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Agent{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchAgents(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledAgents(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, err
}

func toLabeledAgents(set types.AgentSet) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
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

func loadAgent(ctx context.Context, s store.Agents, ID uint64) (res *types.Agent, err error) {
	if ID == 0 {
		return nil, AgentErrInvalidID()
	}

	if res, err = store.LookupAgentByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, AgentErrNotFound()
	}

	return
}
