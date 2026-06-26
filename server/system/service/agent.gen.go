package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type agentAccessController interface {
	CanCreateAgent(context.Context) bool
	CanSearchAgents(context.Context) bool
	CanReadAgent(context.Context, *types.Agent) bool
	CanUpdateAgent(context.Context, *types.Agent) bool
	CanDeleteAgent(context.Context, *types.Agent) bool
}

func (svc *agent) FindByID(ctx context.Context, ID uint64) (res *types.Agent, err error) {
	var (
		aProps = &agentActionProps{agent: &types.Agent{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, AgentActionLookup, err)
}

func (svc *agent) Search(ctx context.Context, filter types.AgentFilter) (set types.AgentSet, f types.AgentFilter, err error) {
	var (
		aProps = &agentActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanSearchAgents(ctx) {
			return AgentErrNotAllowedToSearch()
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, AgentActionSearch, err)
}

func (svc *agent) Create(ctx context.Context, new *types.Agent) (res *types.Agent, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &agentActionProps{agent: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateAgent(ctx) {
			return AgentErrNotAllowedToCreate()
		}
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, AgentActionCreate, err)
}

func (svc *agent) Update(ctx context.Context, upd *types.Agent) (res *types.Agent, err error) {
	var (
		aProps = &agentActionProps{update: upd}
		old    *types.Agent
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadAgent(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setAgent(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return AgentErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return AgentErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Handle = upd.Handle
		res.Status = upd.Status
		res.Revision = upd.Revision
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateAgent(ctx, s, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, AgentActionUpdate, err, old, res)
}

func (svc *agent) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &agentActionProps{}
		res    *types.Agent
	)
	err = func() (err error) {
		if res, err = loadAgent(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setAgent(res)

		if !svc.ac.CanDeleteAgent(ctx, res) {
			return AgentErrNotAllowedToDelete()
		}

		res.DeletedAt = now()
		if err = store.UpdateAgent(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, AgentActionDelete, err)
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

// toLabeledAgents converts to []label.LabeledResource
func toLabeledAgents(set []*types.Agent) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
