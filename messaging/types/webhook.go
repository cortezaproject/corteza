package types

import (
	"io"
	"time"
)

type (
	Webhook struct {
		ID uint64 `json:"id" db:"id"`

		Kind      WebhookKind `json:"kind" db:"kind"`
		AuthToken string      `json:"-" db:"token"`

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
		UserID   uint64

		Avatar io.Reader

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
	OutgoingWebhook WebhookKind = "outgoing"
)
