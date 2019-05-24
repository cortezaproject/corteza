package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/internal/store"
	"github.com/cortezaproject/corteza-server/messaging/internal/service"
	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

var _ = errors.Wrap

type Webhooks struct {
	webhook service.WebhookService
}

func (Webhooks) New() *Webhooks {
	return &Webhooks{}
}

func (ctrl *Webhooks) Get(ctx context.Context, r *request.WebhooksGet) (interface{}, error) {
	return ctrl.webhook.With(ctx).Get(r.WebhookID)
}

func (ctrl *Webhooks) Delete(ctx context.Context, r *request.WebhooksDelete) (interface{}, error) {
	return nil, ctrl.webhook.With(ctx).Delete(r.WebhookID)
}

func (ctrl *Webhooks) List(ctx context.Context, r *request.WebhooksList) (interface{}, error) {
	return ctrl.webhook.With(ctx).Find(&types.WebhookFilter{
		ChannelID:   r.ChannelID,
		OwnerUserID: r.UserID,
	})
}

func (ctrl *Webhooks) Create(ctx context.Context, r *request.WebhooksCreate) (interface{}, error) {
	avatar, err := store.FromAny(r.Avatar, r.AvatarURL)
	if err != nil {
		return nil, err
	}
	defer avatar.Close()
	// Webhook request parameters
	parameters := types.WebhookRequest{
		r.Username,
		r.UserID,
		avatar,
		r.Trigger,
		r.Url,
	}
	return ctrl.webhook.With(ctx).Create(r.Kind, r.ChannelID, parameters)
}

func (ctrl *Webhooks) Update(ctx context.Context, r *request.WebhooksUpdate) (interface{}, error) {
	avatar, err := store.FromAny(r.Avatar, r.AvatarURL)
	if err != nil {
		return nil, err
	}
	defer avatar.Close()
	// Webhook request parameters
	parameters := types.WebhookRequest{
		r.Username,
		r.UserID,
		avatar,
		r.Trigger,
		r.Url,
	}
	return ctrl.webhook.With(ctx).Update(r.WebhookID, r.Kind, r.ChannelID, parameters)
}
