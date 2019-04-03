package types

import (
	"time"

	"mime/multipart"

	"github.com/crusttech/crust/internal/rules"
)

type (
	Webhook struct {
		ID uint64 `json:"id" db:"id"`

		Kind      WebhookKind `json:"kind" db:"webhook_kind"`
		AuthToken string      `json:"-" db:"webhook_token"`

		OwnerUserID uint64 `json:"userId" db:"rel_owner"`

		// Created bot User ID
		UserID    uint64 `json:"userId" db:"rel_user"`
		ChannelID uint64 `json:"channelId" db:"rel_channel"`

		// Outgoing webhook details
		OutgoingTrigger string `json:"trigger" db:"outgoing_trigger"`
		OutgoingURL     string `json:"url" db:"outgoing_url"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	WebhookRequest struct {
		Username string

		Avatar    *multipart.FileHeader
		AvatarURL string

		OutgoingTrigger string
		OutgoingURL     string
	}

	WebhookFilter struct {
		ChannelID       uint64
		OwnerUserID     uint64
		OutgoingTrigger string
	}

	WebhookBody struct {
		Text     string `json:"text"`
		Avatar   string `json:"avatar,omitempty"`
		Username string `json:"username,omitempty"`
	}

	WebhookKind string
)

const (
	IncomingWebhook WebhookKind = "incoming"
	OutgoingWebhook             = "outgoing"
)

// Resource returns a system resource ID for this type
func (wh Webhook) PermissionResource() rules.Resource {
	return WebhookPermissionResource.AppendID(wh.ID)
}
