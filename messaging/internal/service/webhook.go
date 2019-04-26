package service

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/http"
	"github.com/crusttech/crust/internal/store"
	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
	systemTypes "github.com/crusttech/crust/system/types"
)

type (
	webhook struct {
		db     db
		ctx    context.Context
		client *http.Client

		users   systemService.UserService
		webhook repository.WebhookRepository
		perms   PermissionsService
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
		client: client,
	}).With(ctx)
}

func (svc *webhook) With(ctx context.Context) WebhookService {
	db := repository.DB(ctx)
	return &webhook{
		db:     db,
		ctx:    ctx,
		client: svc.client,

		users:   systemService.User(ctx),
		webhook: repository.Webhook(ctx, db),
		perms:   Permissions(ctx),
	}
}

func (svc *webhook) Create(kind types.WebhookKind, channelID uint64, params types.WebhookRequest) (*types.Webhook, error) {
	var userID = repository.Identity(svc.ctx)

	webhook := &types.Webhook{
		Kind:            kind,
		OwnerUserID:     userID,
		ChannelID:       channelID,
		OutgoingTrigger: params.OutgoingTrigger,
		OutgoingURL:     params.OutgoingURL,
	}

	if !svc.perms.CanManageWebhooks(webhook) && !svc.perms.CanManageOwnWebhooks(webhook) {
		return nil, errors.WithStack(ErrNoPermissions)
	}

	botUser := &systemTypes.User{
		Username:      params.Username,
		Name:          params.Username,
		Handle:        params.Username,
		Kind:          systemTypes.BotUser,
		RelatedUserID: userID,
	}

	user, err := svc.users.CreateWithAvatar(botUser, params.Avatar)
	if err != nil {
		return nil, err
	}
	webhook.UserID = user.ID

	wh, err := svc.webhook.Create(webhook)
	if err != nil {
		// cross service rollback (delete user)
		svc.users.Delete(user.ID)
		return nil, err
	}
	return wh, err
}

func (svc *webhook) Update(webhookID uint64, kind types.WebhookKind, channelID uint64, params types.WebhookRequest) (*types.Webhook, error) {
	var userID = repository.Identity(svc.ctx)
	webhook, err := svc.Get(webhookID)
	if err != nil {
		return nil, err
	}

	if !svc.perms.CanManageWebhooks(webhook) && !(webhook.OwnerUserID == userID && svc.perms.CanManageOwnWebhooks(webhook)) {
		return nil, errors.WithStack(ErrNoPermissions)
	}

	botUser, err := svc.users.FindByID(webhook.UserID)
	if err != nil {
		return nil, errors.Wrapf(err, "Error when looking for User ID %d", webhook.UserID)
	}

	if _, err := svc.users.UpdateWithAvatar(botUser, params.Avatar); err != nil {
		return nil, err
	}

	webhook.Kind = kind
	webhook.ChannelID = channelID
	webhook.OutgoingTrigger = params.OutgoingTrigger
	webhook.OutgoingURL = params.OutgoingURL

	return svc.webhook.Update(webhook)
}

func (svc *webhook) Get(webhookID uint64) (*types.Webhook, error) {
	return svc.webhook.Get(webhookID)
}

func (svc *webhook) Find(filter *types.WebhookFilter) (types.WebhookSet, error) {
	return svc.webhook.Find(filter)
}

func (svc *webhook) Delete(webhookID uint64) error {
	webhook, err := svc.Get(webhookID)
	if err != nil {
		return err
	}
	if svc.perms.CanManageWebhooks(webhook) {
		return svc.webhook.Delete(webhookID)
	}
	var userID = repository.Identity(svc.ctx)
	if webhook.OwnerUserID == userID && svc.perms.CanManageOwnWebhooks(webhook) {
		return svc.webhook.Delete(webhookID)
	}
	return errors.WithStack(ErrNoPermissions)
}

func (svc *webhook) DeleteByToken(webhookID uint64, webhookToken string) error {
	return svc.webhook.DeleteByToken(webhookID, webhookToken)
}

func (svc *webhook) Message(webhookID uint64, webhookToken string, username string, avatar io.Reader, message string) (interface{}, error) {
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
func (svc *webhook) Do(webhook *types.Webhook, message string) (*types.Message, error) {
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

func (svc *webhook) sendMessage(webhook *types.Webhook, msg *types.Message, avatar io.Reader) (*types.Message, error) {
	// We need a webhook user context for message service
	ctx := auth.SetIdentityToContext(svc.ctx, &systemTypes.User{ID: webhook.UserID})
	messageSvc := Message(ctx)

	msg.ChannelID = webhook.ChannelID
	msg.UserID = webhook.UserID

	return messageSvc.CreateWithAvatar(msg, avatar)
}
