package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	WebhookRepository interface {
		With(ctx context.Context, db *factory.DB) WebhookRepository

		Create(*types.Webhook) (*types.Webhook, error)
		Update(*types.Webhook) (*types.Webhook, error)

		Get(webhookID uint64) (*types.Webhook, error)
		GetByToken(webhookID uint64, webhookToken string) (*types.Webhook, error)

		Find(filter *types.WebhookFilter) (types.WebhookSet, error)

		Delete(webhookID uint64) error
		DeleteByToken(webhookID uint64, webhookToken string) error
	}

	webhook struct {
		webhook string

		*repository
	}
)

func Webhook(ctx context.Context, db *factory.DB) WebhookRepository {
	return (&webhook{}).With(ctx, db)
}

func (r *webhook) With(ctx context.Context, db *factory.DB) WebhookRepository {
	return &webhook{
		webhook:    "messaging_webhook",
		repository: r.repository.With(ctx, db),
	}
}

func (r *webhook) Create(webhook *types.Webhook) (*types.Webhook, error) {
	webhook.ID = factory.Sonyflake.NextID()
	webhook.CreatedAt = time.Now()

	return webhook, errors.WithStack(r.db().Insert(r.webhook, webhook))
}

func (r *webhook) Update(webhook *types.Webhook) (*types.Webhook, error) {
	webhook.UpdatedAt = timeNowPtr()

	return webhook, errors.WithStack(r.db().Replace(r.webhook, webhook))
}

func (r *webhook) Get(webhookID uint64) (*types.Webhook, error) {
	hook := &types.Webhook{}
	if err := r.db().Get(hook, "select * from "+r.webhook+" where id=?", webhookID); err != nil {
		return nil, errors.WithStack(err)
	}
	return hook, nil
}

func (r *webhook) GetByToken(webhookID uint64, webhookToken string) (*types.Webhook, error) {
	webhook, err := r.Get(webhookID)
	switch {
	case err != nil:
		return nil, errors.WithStack(err)
	case webhook.AuthToken == webhookToken:
		return webhook, nil
	default:
		return nil, errors.New("Invalid Webhook Token")
	}
}

// Find webhooks based on filter
//
// If ChannelID > 0 it returns webhooks created on a specific channel
// If OwnerUserID > 0 it returns webhooks owned by a specific user
func (r *webhook) Find(filter *types.WebhookFilter) (types.WebhookSet, error) {
	params := make([]interface{}, 0)
	vv := types.WebhookSet{}
	sql := "select * from messaging_webhook where 1=1"

	if filter != nil {
		if filter.OwnerUserID > 0 {
			// scope: only channel we have access to
			sql += " AND rel_owner=?"
			params = append(params, filter.OwnerUserID)
		}
		if filter.OutgoingTrigger != "" {
			// scope: only channel we have access to
			sql += " AND outgoing_trigger=?"
			params = append(params, filter.OutgoingTrigger)
		}
		if filter.ChannelID > 0 {
			// scope: only channel we have access to
			sql += " AND rel_channel=?"
			params = append(params, filter.ChannelID)
		}
	}

	return vv, errors.WithStack(r.db().Select(&vv, sql, params...))
}

func (r *webhook) Delete(webhookID uint64) error {
	if _, err := r.Get(webhookID); err != nil {
		return err
	}
	_, err := r.db().Exec("delete from "+r.webhook+" where id=?", webhookID)
	return errors.WithStack(err)
}

func (r *webhook) DeleteByToken(webhookID uint64, webhookToken string) error {
	if _, err := r.GetByToken(webhookID, webhookToken); err != nil {
		return err
	}
	_, err := r.db().Exec("delete from "+r.webhook+" where id=?", webhookID)
	return errors.WithStack(err)
}
