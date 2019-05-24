package service

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/http"
	"github.com/cortezaproject/corteza-server/internal/store"
	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	webhook struct {
		db     db
		ctx    context.Context
		logger *zap.Logger

		client *http.Client

		webhook repository.WebhookRepository
		ac      webhookAccessController
	}

	webhookAccessController interface {
		CanCreateWebhook(context.Context) bool
		CanManageWebhooks(context.Context) bool
		CanManageOwnWebhooks(context.Context, *types.Webhook) bool
	}

	WebhookService interface {
		With(ctx context.Context) WebhookService

		Get(webhookID uint64) (*types.Webhook, error)

		Find(*types.WebhookFilter) (types.WebhookSet, error)

		Delete(webhookID uint64) error
		DeleteByToken(webhookID uint64, webhookToken string) error

		Message(webhookID uint64, webhookToken string, username string, avatar io.Reader, message string) (interface{}, error)

		Create(kind types.WebhookKind, channelID uint64, params types.WebhookRequest) (*types.Webhook, error)
		Update(webhookID uint64, kind types.WebhookKind, channelID uint64, params types.WebhookRequest) (*types.Webhook, error)

		Do(webhook *types.Webhook, message string) (*types.Message, error)
	}
)

func Webhook(ctx context.Context, client *http.Client) WebhookService {
	return (&webhook{
		logger: DefaultLogger.Named("webhook"),

		client: client,
	}).With(ctx)
}

func (svc webhook) With(ctx context.Context) WebhookService {
	db := repository.DB(ctx)
	return &webhook{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		client: svc.client,

		webhook: repository.Webhook(ctx, db),
		ac:      DefaultAccessControl,
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc webhook) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc webhook) Create(kind types.WebhookKind, channelID uint64, params types.WebhookRequest) (*types.Webhook, error) {
	var userID = repository.Identity(svc.ctx)

	// @todo: params.Avatar (io.Reader)

	webhook := &types.Webhook{
		Kind:            kind,
		UserID:          params.UserID,
		OwnerUserID:     userID,
		ChannelID:       channelID,
		OutgoingTrigger: params.OutgoingTrigger,
		OutgoingURL:     params.OutgoingURL,
	}

	if !svc.ac.CanCreateWebhook(svc.ctx) {
		return nil, ErrNoPermissions.withStack()
	}

	return svc.webhook.Create(webhook)
}

func (svc webhook) Update(webhookID uint64, kind types.WebhookKind, channelID uint64, params types.WebhookRequest) (*types.Webhook, error) {
	if webhookID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	webhook, err := svc.Get(webhookID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanManageOwnWebhooks(svc.ctx, webhook) || !svc.ac.CanManageWebhooks(svc.ctx) {
		return nil, ErrNoPermissions.withStack()
	}

	// @todo: params.Avatar (io.Reader)

	webhook.Kind = kind
	webhook.ChannelID = channelID
	webhook.OutgoingTrigger = params.OutgoingTrigger
	webhook.OutgoingURL = params.OutgoingURL

	return svc.webhook.Update(webhook)
}

func (svc webhook) Get(webhookID uint64) (*types.Webhook, error) {
	return svc.webhook.Get(webhookID)
}

func (svc webhook) Find(filter *types.WebhookFilter) (types.WebhookSet, error) {
	return svc.webhook.Find(filter)
}

func (svc webhook) Delete(webhookID uint64) error {
	webhook, err := svc.Get(webhookID)
	if err != nil {
		return err
	}
	if !svc.ac.CanManageOwnWebhooks(svc.ctx, webhook) || !svc.ac.CanManageWebhooks(svc.ctx) {
		return svc.webhook.Delete(webhookID)
	}
	var userID = repository.Identity(svc.ctx)
	if webhook.OwnerUserID == userID && svc.ac.CanManageOwnWebhooks(svc.ctx, webhook) {
		return svc.webhook.Delete(webhookID)
	}
	return ErrNoPermissions.withStack()
}

func (svc webhook) DeleteByToken(webhookID uint64, webhookToken string) error {
	return svc.webhook.DeleteByToken(webhookID, webhookToken)
}

func (svc webhook) Message(webhookID uint64, webhookToken string, username string, avatar io.Reader, message string) (interface{}, error) {
	if webhook, err := svc.webhook.GetByToken(webhookID, webhookToken); err != nil {
		return nil, err
	} else {
		msg := &types.Message{
			Message: message,
			Meta: &types.MessageMeta{
				Username: username,
			},
		}
		return svc.sendMessage(webhook, msg, avatar)
	}
}

// Do executes an outgoing HTTP webhook request
func (svc webhook) Do(webhook *types.Webhook, message string) (*types.Message, error) {
	if webhook.Kind != types.OutgoingWebhook {
		return nil, errors.Errorf("Unsupported webhook type: %s", webhook.Kind)
	}

	// replace url query %s with message
	url := strings.Replace(webhook.OutgoingURL, "%s", url.QueryEscape(message), -1)

	// post body contains only `text`
	requestBody := types.WebhookBody{
		Text: message,
	}
	req, err := svc.client.Post(url, requestBody)
	if err != nil {
		return nil, err
	}

	// execute outgoing webhook
	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse response body
	responseBody := types.WebhookBody{}
	contentType := resp.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "text/plain"):
		// keep plain/text as-is
		if b, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, errors.WithStack(err)
		} else {
			responseBody.Text = string(b)
		}
	default:
		switch resp.StatusCode {
		case 200:
			// assume the response is an expected json structure
			if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
				return nil, errors.WithStack(err)
			}
			if responseBody.Text == "" {
				return nil, errors.New("Empty webhook response")
			}
		default:
			return nil, http.ToError(resp)
		}
	}

	msg := &types.Message{
		Message: responseBody.Text,
		Meta: &types.MessageMeta{
			Username: responseBody.Username,
		},
	}

	avatar, err := store.FromURL(responseBody.Avatar)
	if err != nil {
		return svc.sendMessage(webhook, msg, nil)
	}
	defer avatar.Close()

	return svc.sendMessage(webhook, msg, avatar)
}

func (svc webhook) sendMessage(webhook *types.Webhook, msg *types.Message, avatar io.Reader) (*types.Message, error) {
	// We need a webhook user context for message service
	ctx := auth.SetIdentityToContext(svc.ctx, auth.NewIdentity(webhook.UserID))
	messageSvc := Message(ctx)

	msg.ChannelID = webhook.ChannelID
	msg.UserID = webhook.UserID

	return messageSvc.CreateWithAvatar(msg, avatar)
}
