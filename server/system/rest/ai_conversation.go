package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
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

// makeFilter builds the search filter for the generated List controller.
func (ctrl *AiConversation) makeFilter(ctx context.Context, r *request.AiConversationList) (types.AiConversationFilter, error) {
	var (
		err error
		f   = types.AiConversationFilter{
			AgentID: r.AgentID,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	if r.IncTotal {
		f.IncTotal = true
	}

	return f, nil
}

// makePayload wraps a single resource for the generated Read controller.
// The original Read returned the resource (or error) directly, so this is
// a straight passthrough that preserves that behavior.
func (ctrl *AiConversation) makePayload(_ context.Context, m *types.AiConversation, err error) (*types.AiConversation, error) {
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (ctrl *AiConversation) Continue(ctx context.Context, r *request.AiConversationContinue) (interface{}, error) {
	conv, err := ctrl.svc.FindByID(ctx, r.AiConversationID)
	if err != nil {
		return nil, err
	}

	return service.DefaultAgenticRuntime.Run(ctx, &runtime.AgentRequest{
		AgentID:        conv.AgentID,
		ConversationID: conv.ID,
		Input:          r.Input,
	})
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
