package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/internal/mail"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Notification struct {
		notification service.NotificationService
	}

	contentPayload struct {
		Plain string `json:"plain"`
		Html  string `json:"html"`
	}
)

func (Notification) New() *Notification {
	return &Notification{
		notification: service.DefaultNotification,
	}
}

// EmailSend assembles Email Message and pushes message to notification service
func (ctrl *Notification) EmailSend(ctx context.Context, r *request.NotificationEmailSend) (interface{}, error) {
	ntf := ctrl.notification.With(ctx)

	msg := mail.New()
	if err := ntf.AttachEmailRecipients(msg, "To", r.To...); err != nil {
		return false, err
	}

	if err := ntf.AttachEmailRecipients(msg, "Cc", r.Cc...); err != nil {
		return false, err
	}

	var cp = contentPayload{}
	if err := r.Content.Unmarshal(&cp); err != nil {
		return false, errors.Wrap(err, "Could not unmarshal content")
	} else {
		if len(cp.Html) > 0 {
			msg.SetBody("text/html", cp.Html)

		}

		if len(cp.Plain) > 0 {
			msg.SetBody("text/plain", cp.Plain)
		}
	}

	msg.SetHeader("Subject", r.Subject)

	if err := ctrl.notification.With(ctx).SendEmail(msg); err != nil {
		return false, err
	} else {
		return true, nil
	}
}
