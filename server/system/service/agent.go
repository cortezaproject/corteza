package service

import (
	"context"
	"strings"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/label"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Search, Create, Update, DeleteByID and the
// toLabeledAgents helper) is generated in agent.gen.go from system/agent.cue.
//
// This file owns the constructor, the on<Op> custom bodies the generated FindByID /
// Search / Create / Update delegate to, the hand-written UndeleteByID (the
// generated undelete is disabled because it uses the non-standard
// CanDeleteAgent + AgentErrNotAllowedToDelete pairing), and the
// resource-specific helpers (Get). The loadAgent helper is
// generated in agent.gen.go.

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

// onLookup is the custom body for the generated FindByID. The generated method
// owns the action-log scaffold + recordAction; the load, read access check and
// label load (which the standard template does not emit) live here.
func (svc *agent) guard(ctx context.Context, res *types.Agent) error {
	return guardProjectWritable(ctx, svc.store, res.ProjectID)
}

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
// name check, Revision / Status defaults and persistence live here.
func (svc *agent) onCreate(ctx context.Context, new *types.Agent) (err error) {
	if err = validateAgent(new); err != nil {
		return
	}

	new.ID = nextID()
	new.CreatedAt = *now()
	// An agent carries a tool allow-list and reads data on somebody's behalf,
	// so it is the resource where authorship matters most; it was the one
	// resource recording none.
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	new.Revision = 1

	if new.Status == "" {
		new.Status = "active"
	}

	defaultInvocation(&new.Invocation)

	if err = store.CreateAgent(ctx, svc.store, new); err != nil {
		return
	}

	if err = label.Create(ctx, svc.store, new); err != nil {
		return
	}

	return nil
}

// validateAgent refuses an agent with no name, the one field storing it
// requires. Provider and model are left to the editor and to run time.
func validateAgent(res *types.Agent) error {
	if strings.TrimSpace(res.Meta.Short) == "" {
		return AgentErrMissingName()
	}
	return nil
}

// defaultInvocation turns user invocation on for an agent that named no way of
// being invoked at all.
//
// An agent with neither user nor system invocation is one nothing can ever run,
// which is never what someone meant to build; the webapp's own Agent model
// starts with user invocation enabled, so an agent created through the API
// otherwise behaves differently from the identical one created in the editor
// — and says so only at exec time, as "not available for user invocation".
func defaultInvocation(inv *types.AgentInvocation) {
	if inv.User.Enabled || inv.System.Enabled {
		return
	}
	inv.User.Enabled = true
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction; the update access check (on the
// incoming resource), the name check, the stale guard, the Revision bump and
// the whole-record persistence live here.
func (svc *agent) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Agent, aProps *agentActionProps, _ func() error, _ func() error) error {
	if !svc.ac.CanUpdateAgent(ctx, upd) {
		return AgentErrNotAllowedToUpdate()
	}

	if err := validateAgent(upd); err != nil {
		return err
	}

	upd.Revision = res.Revision + 1
	upd.UpdatedAt = now()
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	upd.CreatedAt = res.CreatedAt
	upd.CreatedBy = res.CreatedBy
	upd.DeletedAt = res.DeletedAt
	// ProjectID is set at create and isn't part of the update payload; preserve
	// the existing owning project so an update never unlinks the agent from it
	// (otherwise it would vanish from the project-scoped resource graph).
	upd.ProjectID = res.ProjectID

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
