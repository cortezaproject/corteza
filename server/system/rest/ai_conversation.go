package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	AiConversation struct {
		svc aiConversationService
		ac  aiConversationAccessController
	}

	aiConversationService interface {
		FindByID(ctx context.Context, ID uint64) (*types.AiConversation, error)
		Create(ctx context.Context, new *types.AiConversation) (*types.AiConversation, error)
		Update(ctx context.Context, upd *types.AiConversation) (*types.AiConversation, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.AiConversationFilter) (types.AiConversationSet, types.AiConversationFilter, error)
	}

	aiConversationAccessController interface {
		CanCreateAiConversation(context.Context) bool
		CanSearchAiConversations(context.Context) bool
	}
)

func (AiConversation) New() *AiConversation {
	return &AiConversation{
		svc: service.DefaultAiConversation,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl *AiConversation) List(ctx context.Context, r *request.AiConversationList) (interface{}, error) {
	var (
		err error
		f   = types.AiConversationFilter{
			AgentID: r.AgentID,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	if r.IncTotal {
		f.IncTotal = true
	}

	set, f, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl *AiConversation) Read(ctx context.Context, r *request.AiConversationRead) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.AiConversationID)
}

func (ctrl *AiConversation) Delete(ctx context.Context, r *request.AiConversationDelete) (interface{}, error) {
	return nil, ctrl.svc.DeleteByID(ctx, r.AiConversationID)
}

func (ctrl *AiConversation) Undelete(ctx context.Context, r *request.AiConversationUndelete) (interface{}, error) {
	return nil, ctrl.svc.UndeleteByID(ctx, r.AiConversationID)
}

func (ctrl *AiConversation) makeFilterPayload(_ context.Context, nn types.AiConversationSet, f types.AiConversationFilter, err error) (*aiConversationSetPayload, error) {
	if err != nil {
		return nil, err
	}

	return &aiConversationSetPayload{Filter: f, Set: nn}, nil
}

type aiConversationSetPayload struct {
	Filter types.AiConversationFilter `json:"filter"`
	Set    types.AiConversationSet    `json:"set"`
}
