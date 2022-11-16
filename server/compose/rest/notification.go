package rest

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
)

type (
	Notification struct {
		svc notificationService
	}

	contentPayload struct {
		Plain string `json:"plain"`
		HTML  string `json:"html"`
	}

	notificationService interface {
		SendEmail(context.Context, *types.EmailNotification) error
	}
)

func (Notification) New() *Notification {
	return &Notification{
		svc: service.DefaultNotification,
	}
}

// EmailSend assembles Email Message and pushes message to notification service
func (ctrl *Notification) EmailSend(ctx context.Context, r *request.NotificationEmailSend) (interface{}, error) {
	ntf := &types.EmailNotification{
		To:                r.To,
		Cc:                r.Cc,
		ReplyTo:           r.ReplyTo,
		Subject:           r.Subject,
		RemoteAttachments: r.RemoteAttachments,
	}

	var cp = contentPayload{}
	if err := r.Content.Unmarshal(&cp); err != nil {
		return false, fmt.Errorf("could not unmarshal content: %w", err)
	} else {
		if len(cp.HTML) > 0 {
			ntf.ContentHTML = cp.HTML

		}

		if len(cp.Plain) > 0 {
			ntf.ContentPlain = cp.Plain
		}
	}

	if err := ctrl.svc.SendEmail(ctx, ntf); err != nil {
		return false, err
	} else {
		return true, nil
	}
}
