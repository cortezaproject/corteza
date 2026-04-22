package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Chatbot struct {
		svc chatbotService
		ac  chatbotAccessController
	}

	chatbotPayload struct {
		*types.Chatbot

		CanGrant         bool `json:"canGrant"`
		CanUpdateChatbot bool `json:"canUpdateChatbot"`
		CanDeleteChatbot bool `json:"canDeleteChatbot"`
	}

	chatbotSetPayload struct {
		Filter types.ChatbotFilter `json:"filter"`
		Set    []*chatbotPayload   `json:"set"`
	}

	chatbotService interface {
		FindByID(ctx context.Context, ID uint64) (c *types.Chatbot, err error)
		Create(ctx context.Context, new *types.Chatbot) (c *types.Chatbot, err error)
		Update(ctx context.Context, upd *types.Chatbot) (c *types.Chatbot, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		UndeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.ChatbotFilter) (set types.ChatbotSet, f types.ChatbotFilter, err error)
		RegenerateWidgetKey(ctx context.Context, ID uint64) (c *types.Chatbot, err error)
	}

	chatbotAccessController interface {
		CanGrant(context.Context) bool

		CanCreateChatbot(context.Context) bool
		CanUpdateChatbot(context.Context, *types.Chatbot) bool
		CanDeleteChatbot(context.Context, *types.Chatbot) bool
	}
)

func (Chatbot) New() *Chatbot {
	return &Chatbot{
		svc: service.DefaultChatbot,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl *Chatbot) List(ctx context.Context, r *request.ChatbotList) (interface{}, error) {
	var (
		err error
		f   = types.ChatbotFilter{
			Query:   r.Query,
			Handle:  r.Handle,
			Labels:  r.Labels,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Chatbot) Create(ctx context.Context, r *request.ChatbotCreate) (interface{}, error) {
	var (
		err error
		c   = &types.Chatbot{
			Handle:         r.Handle,
			Name:           r.Name,
			Enabled:        r.Enabled,
			SessionTTL:     r.SessionTTL,
			AllowedOrigins: r.AllowedOrigins,
			Handoff:        r.Handoff,
			Styling:        r.Styling,
			Scenarios:      r.Scenarios,
			Labels:         r.Labels,
		}
	)

	c, err = ctrl.svc.Create(ctx, c)
	return ctrl.makePayload(ctx, c, err)
}

func (ctrl *Chatbot) Read(ctx context.Context, r *request.ChatbotRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ChatbotID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Chatbot) Update(ctx context.Context, r *request.ChatbotUpdate) (interface{}, error) {
	var (
		err error
		c   = &types.Chatbot{
			ID:             r.ChatbotID,
			Handle:         r.Handle,
			Name:           r.Name,
			Enabled:        r.Enabled,
			SessionTTL:     r.SessionTTL,
			AllowedOrigins: r.AllowedOrigins,
			Handoff:        r.Handoff,
			Styling:        r.Styling,
			Scenarios:      r.Scenarios,
			UpdatedAt:      r.UpdatedAt,
			Labels:         r.Labels,
		}
	)

	c, err = ctrl.svc.Update(ctx, c)
	return ctrl.makePayload(ctx, c, err)
}

func (ctrl *Chatbot) Delete(ctx context.Context, r *request.ChatbotDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ChatbotID)
}

func (ctrl *Chatbot) Undelete(ctx context.Context, r *request.ChatbotUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ChatbotID)
}

func (ctrl *Chatbot) RegenerateWidgetKey(ctx context.Context, r *request.ChatbotRegenerateWidgetKey) (interface{}, error) {
	c, err := ctrl.svc.RegenerateWidgetKey(ctx, r.ChatbotID)
	return ctrl.makePayload(ctx, c, err)
}

func (ctrl *Chatbot) makePayload(ctx context.Context, c *types.Chatbot, err error) (*chatbotPayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	p := &chatbotPayload{
		Chatbot: c,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateChatbot: ctrl.ac.CanUpdateChatbot(ctx, c),
		CanDeleteChatbot: ctrl.ac.CanDeleteChatbot(ctx, c),
	}

	return p, nil
}

func (ctrl *Chatbot) makeFilterPayload(ctx context.Context, nn types.ChatbotSet, f types.ChatbotFilter, err error) (*chatbotSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &chatbotSetPayload{Filter: f, Set: make([]*chatbotPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
