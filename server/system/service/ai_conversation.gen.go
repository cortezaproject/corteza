package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *aiConversation) FindByID(ctx context.Context, ID uint64) (res *types.AiConversation, err error) {
	var (
		aProps = &aiConversationActionProps{aiConversation: &types.AiConversation{ID: ID}}
	)

	err = func() error {
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
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, AiConversationActionUpdate, err)
}

func (svc *aiConversation) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &aiConversationActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, AiConversationActionDelete, err)
}

func (svc *aiConversation) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &aiConversationActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, AiConversationActionUndelete, err)
}
