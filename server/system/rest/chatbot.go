package rest

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Chatbot struct {
		svc        chatbotService
		attachment service.AttachmentService
		ac         chatbotAccessController
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
		svc:        service.DefaultChatbot,
		attachment: service.DefaultAttachment,
		ac:         service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *Chatbot) makeFilter(ctx context.Context, r *request.ChatbotList) (types.ChatbotFilter, error) {
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
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params (including Labels).
func (ctrl *Chatbot) beforeCreate(ctx context.Context, res *types.Chatbot, r *request.ChatbotCreate) error {
	res.AllowedOrigins = r.AllowedOrigins
	res.Handoff = r.Handoff
	res.Styling = r.Styling
	res.Scenarios = r.Scenarios
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (including Labels) and ID/UpdatedAt.
func (ctrl *Chatbot) beforeUpdate(ctx context.Context, res *types.Chatbot, r *request.ChatbotUpdate) error {
	res.AllowedOrigins = r.AllowedOrigins
	res.Handoff = r.Handoff
	res.Styling = r.Styling
	res.Scenarios = r.Scenarios
	return nil
}

func (ctrl *Chatbot) RegenerateWidgetKey(ctx context.Context, r *request.ChatbotRegenerateWidgetKey) (interface{}, error) {
	c, err := ctrl.svc.RegenerateWidgetKey(ctx, r.ChatbotID)
	return ctrl.makePayload(ctx, c, err)
}

func (ctrl *Chatbot) UploadAsset(ctx context.Context, r *request.ChatbotUploadAsset) (interface{}, error) {
	if r.Upload == nil {
		return nil, fmt.Errorf("missing upload")
	}

	// Permission follows chatbot edit.
	cb, err := ctrl.svc.FindByID(ctx, r.ChatbotID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanUpdateChatbot(ctx, cb) {
		return nil, service.ChatbotErrNotAllowedToUpdate()
	}

	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	att, err := ctrl.attachment.CreateChatbotAttachment(ctx, cb.ID, r.Upload.Filename, r.Upload.Size, file)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"attachmentID": strconv.FormatUint(att.ID, 10),
		"url":          chatbotAssetURL(cb.WidgetKey, att.ID, att.Name),
	}, nil
}

// chatbotAssetURL builds the widget-scoped read URL. The frontend stores only
// the attachmentID; the full URL is a convenience for immediate preview.
func chatbotAssetURL(widgetKey string, attachmentID uint64, name string) string {
	return fmt.Sprintf(
		"/api/widget/v1/asset/%d/%s?widgetKey=%s",
		attachmentID,
		url.PathEscape(name),
		url.QueryEscape(widgetKey),
	)
}

func (ctrl *Chatbot) makePayload(ctx context.Context, c *types.Chatbot, err error) (*chatbotPayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	resolveChatbotAssetURLs(c)

	p := &chatbotPayload{
		Chatbot: c,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateChatbot: ctrl.ac.CanUpdateChatbot(ctx, c),
		CanDeleteChatbot: ctrl.ac.CanDeleteChatbot(ctx, c),
	}

	return p, nil
}

// resolveChatbotAssetURLs fills LogoURL / IconURL from the corresponding
// attachment IDs so callers that consume URLs directly (embed snippet, widget
// preview) keep working without separate resolution.
func resolveChatbotAssetURLs(c *types.Chatbot) {
	if c.Styling.LogoAttachmentID != 0 {
		c.Styling.LogoURL = chatbotAssetURL(c.WidgetKey, c.Styling.LogoAttachmentID, "logo")
	}
	if c.Styling.Launcher.IconAttachmentID != 0 {
		c.Styling.Launcher.IconURL = chatbotAssetURL(c.WidgetKey, c.Styling.Launcher.IconAttachmentID, "icon")
	}
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
