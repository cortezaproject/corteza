package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	aiConversation struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        aiConversationAccessController
	}

	aiConversationAccessController interface {
		CanCreateAiConversation(context.Context) bool
		CanSearchAiConversations(context.Context) bool
		CanReadAiConversation(context.Context, *types.AiConversation) bool
		CanUpdateAiConversation(context.Context, *types.AiConversation) bool
		CanDeleteAiConversation(context.Context, *types.AiConversation) bool
	}
)

func AiConversation() *aiConversation {
	return &aiConversation{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc *aiConversation) FindByID(ctx context.Context, ID uint64) (conv *types.AiConversation, err error) {
	var (
		aaProps = &aiConversationActionProps{}
	)

	err = func() error {
		if !svc.ac.CanReadAiConversation(ctx, &types.AiConversation{ID: ID}) {
			return AiConversationErrNotAllowedToRead(aaProps)
		}

		conv, err = store.LookupAiConversationByID(ctx, svc.store, ID)
		if err != nil {
			return err
		}

		aaProps.setAiConversation(conv)
		return nil
	}()

	return conv, svc.recordAction(ctx, aaProps, AiConversationActionLookup, err)
}

func (svc *aiConversation) Create(ctx context.Context, new *types.AiConversation) (conv *types.AiConversation, err error) {
	var (
		aaProps = &aiConversationActionProps{new: new}
	)

	err = func() error {
		if !svc.ac.CanCreateAiConversation(ctx) {
			return AiConversationErrNotAllowedToCreate(aaProps)
		}

		new.ID = id.Next()
		new.CreatedAt = *now()

		if err = store.CreateAiConversation(ctx, svc.store, new); err != nil {
			return err
		}

		conv = new
		aaProps.setAiConversation(conv)
		return nil
	}()

	return conv, svc.recordAction(ctx, aaProps, AiConversationActionCreate, err)
}

func (svc *aiConversation) Update(ctx context.Context, upd *types.AiConversation) (conv *types.AiConversation, err error) {
	var (
		aaProps = &aiConversationActionProps{update: upd}
	)

	err = func() error {
		if !svc.ac.CanUpdateAiConversation(ctx, upd) {
			return AiConversationErrNotAllowedToUpdate(aaProps)
		}

		conv, err = store.LookupAiConversationByID(ctx, svc.store, upd.ID)
		if err != nil {
			return AiConversationErrNotFound(aaProps)
		}

		if isStale(upd.UpdatedAt, conv.UpdatedAt, conv.CreatedAt) {
			return AiConversationErrStaleData(aaProps)
		}

		conv.AgentID = upd.AgentID
		conv.Messages = upd.Messages
		conv.TokenCount = upd.TokenCount
		conv.UpdatedAt = now()

		if err = store.UpdateAiConversation(ctx, svc.store, conv); err != nil {
			return err
		}

		aaProps.setAiConversation(conv)
		return nil
	}()

	return conv, svc.recordAction(ctx, aaProps, AiConversationActionUpdate, err)
}

func (svc *aiConversation) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &aiConversationActionProps{}
		conv    *types.AiConversation
	)

	err = func() error {
		conv, err = store.LookupAiConversationByID(ctx, svc.store, ID)
		if err != nil {
			return err
		}

		if !svc.ac.CanDeleteAiConversation(ctx, conv) {
			return AiConversationErrNotAllowedToDelete(aaProps)
		}

		conv.DeletedAt = now()
		if err = store.UpdateAiConversation(ctx, svc.store, conv); err != nil {
			return err
		}

		aaProps.setAiConversation(conv)
		return nil
	}()

	return svc.recordAction(ctx, aaProps, AiConversationActionDelete, err)
}

func (svc *aiConversation) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &aiConversationActionProps{}
		conv    *types.AiConversation
	)

	err = func() error {
		conv, err = store.LookupAiConversationByID(ctx, svc.store, ID)
		if err != nil {
			return err
		}

		if !svc.ac.CanDeleteAiConversation(ctx, conv) {
			return AiConversationErrNotAllowedToDelete(aaProps)
		}

		conv.DeletedAt = nil
		if err = store.UpdateAiConversation(ctx, svc.store, conv); err != nil {
			return err
		}

		aaProps.setAiConversation(conv)
		return nil
	}()

	return svc.recordAction(ctx, aaProps, AiConversationActionUndelete, err)
}

func (svc *aiConversation) Search(ctx context.Context, filter types.AiConversationFilter) (set types.AiConversationSet, f types.AiConversationFilter, err error) {
	var (
		aaProps = &aiConversationActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanSearchAiConversations(ctx) {
			return AiConversationErrNotAllowedToSearch(aaProps)
		}

		set, f, err = store.SearchAiConversations(ctx, svc.store, filter)
		return err
	}()

	return set, f, svc.recordAction(ctx, aaProps, AiConversationActionSearch, err)
}

func loadAiConversation(ctx context.Context, s store.Agents, ID uint64) (res *types.Agent, err error) {
	if ID == 0 {
		return nil, AgentErrInvalidID()
	}

	if res, err = store.LookupAgentByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, AgentErrNotFound()
	}

	return
}
