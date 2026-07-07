package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type aiConversationAccessController interface {
	CanCreateAiConversation(context.Context) bool
	CanSearchAiConversations(context.Context) bool
	CanReadAiConversation(context.Context, *types.AiConversation) bool
	CanUpdateAiConversation(context.Context, *types.AiConversation) bool
	CanDeleteAiConversation(context.Context, *types.AiConversation) bool
}

type aiConversation struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        aiConversationAccessController
}

func AiConversation() *aiConversation {
	return &aiConversation{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

type aiConversationServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *aiConversation) FindByID(ctx context.Context, ID uint64) (res *types.AiConversation, err error) {
	var (
		aProps = &aiConversationActionProps{aiConversation: &types.AiConversation{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, AiConversationActionLookup, err)
}

func (svc *aiConversation) Search(ctx context.Context, filter types.AiConversationFilter) (set types.AiConversationSet, f types.AiConversationFilter, err error) {
	var (
		aProps = &aiConversationActionProps{filter: &filter}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, AiConversationActionSearch, err)
}

func (svc *aiConversation) Create(ctx context.Context, new *types.AiConversation) (res *types.AiConversation, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &aiConversationActionProps{aiConversation: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateAiConversation(ctx) {
			return AiConversationErrNotAllowedToCreate()
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, AiConversationActionCreate, err)
}

func (svc *aiConversation) Update(ctx context.Context, upd *types.AiConversation) (res *types.AiConversation, err error) {
	var (
		aProps = &aiConversationActionProps{update: upd}
		old    *types.AiConversation
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadAiConversation(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setAiConversation(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return AiConversationErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.AgentID = upd.AgentID
		res.TokenCount = upd.TokenCount
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateAiConversation(ctx, s, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, AiConversationActionUpdate, err, old, res)
}

func (svc *aiConversation) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &aiConversationActionProps{}
		res    *types.AiConversation
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadAiConversation(ctx, s, ID); err != nil {
			return
		}

		aProps.setAiConversation(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, AiConversationActionDelete, err)
}

func (svc *aiConversation) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &aiConversationActionProps{}
		res    *types.AiConversation
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadAiConversation(ctx, s, ID); err != nil {
			return
		}

		aProps.setAiConversation(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, AiConversationActionUndelete, err)
}

func loadAiConversation(ctx context.Context, s store.AiConversations, ID uint64) (res *types.AiConversation, err error) {
	if ID == 0 {
		return nil, AiConversationErrInvalidID()
	}

	if res, err = store.LookupAiConversationByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, AiConversationErrNotFound()
	}

	return
}

func (svc *aiConversation) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *aiConversation) scopeServices(ctx context.Context) *aiConversationServices {
	return &aiConversationServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
