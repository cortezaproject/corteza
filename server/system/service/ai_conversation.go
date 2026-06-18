package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

func (svc *aiConversation) onLookup(ctx context.Context, ID uint64, aProps *aiConversationActionProps) (conv *types.AiConversation, err error) {
	if !svc.ac.CanReadAiConversation(ctx, &types.AiConversation{ID: ID}) {
		return nil, AiConversationErrNotAllowedToRead(aProps)
	}

	conv, err = store.LookupAiConversationByID(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	aProps.setAiConversation(conv)
	return conv, nil
}

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold + the standard CanCreateAiConversation check.
func (svc *aiConversation) onCreate(ctx context.Context, new *types.AiConversation) (err error) {
	new.ID = id.Next()
	new.CreatedAt = *now()

	if err = store.CreateAiConversation(ctx, svc.store, new); err != nil {
		return err
	}

	return nil
}

func (svc *aiConversation) onUpdate(ctx context.Context, upd *types.AiConversation, aProps *aiConversationActionProps) (conv *types.AiConversation, err error) {
	if !svc.ac.CanUpdateAiConversation(ctx, upd) {
		return nil, AiConversationErrNotAllowedToUpdate(aProps)
	}

	conv, err = store.LookupAiConversationByID(ctx, svc.store, upd.ID)
	if err != nil {
		return nil, AiConversationErrNotFound(aProps)
	}

	if isStale(upd.UpdatedAt, conv.UpdatedAt, conv.CreatedAt) {
		return nil, AiConversationErrStaleData(aProps)
	}

	conv.AgentID = upd.AgentID
	conv.Messages = upd.Messages
	conv.TokenCount = upd.TokenCount
	conv.UpdatedAt = now()

	if err = store.UpdateAiConversation(ctx, svc.store, conv); err != nil {
		return nil, err
	}

	aProps.setAiConversation(conv)
	return conv, nil
}

func (svc *aiConversation) onDelete(ctx context.Context, ID uint64, aProps *aiConversationActionProps) (err error) {
	conv, err := store.LookupAiConversationByID(ctx, svc.store, ID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteAiConversation(ctx, conv) {
		return AiConversationErrNotAllowedToDelete(aProps)
	}

	conv.DeletedAt = now()
	if err = store.UpdateAiConversation(ctx, svc.store, conv); err != nil {
		return err
	}

	aProps.setAiConversation(conv)
	return nil
}

func (svc *aiConversation) onUndelete(ctx context.Context, ID uint64, aProps *aiConversationActionProps) (err error) {
	conv, err := store.LookupAiConversationByID(ctx, svc.store, ID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteAiConversation(ctx, conv) {
		return AiConversationErrNotAllowedToDelete(aProps)
	}

	conv.DeletedAt = nil
	if err = store.UpdateAiConversation(ctx, svc.store, conv); err != nil {
		return err
	}

	aProps.setAiConversation(conv)
	return nil
}

func (svc *aiConversation) onSearch(ctx context.Context, filter types.AiConversationFilter, aProps *aiConversationActionProps) (set types.AiConversationSet, f types.AiConversationFilter, err error) {
	if !svc.ac.CanSearchAiConversations(ctx) {
		return nil, filter, AiConversationErrNotAllowedToSearch(aProps)
	}

	set, f, err = store.SearchAiConversations(ctx, svc.store, filter)
	return set, f, err
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
