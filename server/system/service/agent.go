package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/label"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/agentic/tcl"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Search, Create, Update, DeleteByID and the
// toLabeledAgents helper) is generated in agent.gen.go from system/agent.cue.
//
// This file owns the struct, access-controller interface, constructor, the
// LLM-validator opt-in, the on<Op> custom bodies the generated FindByID /
// Search / Create / Update delegate to, the hand-written UndeleteByID (the
// generated undelete is disabled because it uses the non-standard
// CanDeleteAgent + AgentErrNotAllowedToDelete pairing), and the
// resource-specific helpers (Get, prepareTCL). The loadAgent helper is
// generated in agent.gen.go.

type (
	agent struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        agentAccessController
		llm       agentLLMValidator
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

// onLookup is the custom body for the generated FindByID. The generated method
// owns the action-log scaffold + recordAction; the load, read access check and
// label load (which the standard template does not emit) live here.
func (svc *agent) onLookup(ctx context.Context, ID uint64, aProps *agentActionProps) (a *types.Agent, err error) {
	if a, err = loadAgent(ctx, svc.store, ID); err != nil {
		return nil, err
	}

	aProps.setAgent(a)

	if !svc.ac.CanReadAgent(ctx, a) {
		return nil, AgentErrNotAllowedToRead()
	}

	if err = label.Load(ctx, svc.store, a); err != nil {
		return nil, err
	}

	return a, nil
}

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold + recordAction + the CanCreateAgent check; the
// Revision / Status defaults, optional temperature validation, prepareTCL and
// persistence live here.
func (svc *agent) onCreate(ctx context.Context, new *types.Agent) (err error) {
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

	return nil
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction; the update access check (on the
// incoming resource), the stale guard, the Revision bump, optional temperature
// validation, prepareTCL and the whole-record persistence live here.
func (svc *agent) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Agent, aProps *agentActionProps, _ func() error, _ func() error) error {
	if !svc.ac.CanUpdateAgent(ctx, upd) {
		return AgentErrNotAllowedToUpdate()
	}

	upd.Revision = res.Revision + 1
	upd.UpdatedAt = now()
	upd.CreatedAt = res.CreatedAt
	upd.DeletedAt = res.DeletedAt

	if upd.Execution.Model.Temperature != nil && svc.llm != nil {
		if err := svc.llm.ValidateTemperature(ctx, upd.Execution.Model.LLMProviderID, upd.Execution.Model.Model, upd.Execution.Model.Temperature); err != nil {
			return err
		}
	}

	prepareTCL(&upd.Behavior)

	if err := store.UpdateAgent(ctx, s, upd); err != nil {
		return err
	}

	// copy fields back so the caller (generated Update) returns the updated record
	*res = *upd

	return label.Update(ctx, s, upd)
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

// onSearch is the custom body for the generated Search. The generated method
// owns the action-log scaffold + recordAction + the CanSearchAgents check; the
// Check predicate, label filtering and label load live here.
func (svc *agent) onSearch(ctx context.Context, filter types.AgentFilter, aProps *agentActionProps) (set types.AgentSet, f types.AgentFilter, err error) {
	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Agent) (bool, error) {
		if !svc.ac.CanReadAgent(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	if len(filter.Labels) > 0 {
		filter.LabeledIDs, err = label.Search(
			ctx,
			svc.store,
			types.Agent{}.LabelResourceKind(),
			filter.Labels,
		)

		if err != nil {
			return set, f, err
		}

		// labels specified but no labeled resources found
		if len(filter.LabeledIDs) == 0 {
			return set, f, nil
		}
	}

	if set, f, err = store.SearchAgents(ctx, svc.store, filter); err != nil {
		return set, f, err
	}

	if err = label.Load(ctx, svc.store, toLabeledAgents(set)...); err != nil {
		return set, f, err
	}

	return set, f, nil
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
