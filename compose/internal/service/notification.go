package service

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"

	"github.com/cortezaproject/corteza-server/internal/mail"
)

type (
	notification struct {
		ctx    context.Context
		logger *zap.Logger
	}

	NotificationService interface {
		With(ctx context.Context) NotificationService

		SendEmail(message *gomail.Message) error
		AttachEmailRecipients(message *gomail.Message, field string, recipients ...string) error
	}
)

func Notification() NotificationService {
	return (&notification{
		logger: DefaultLogger.Named("notification"),
	}).With(context.Background())
}

func (svc notification) With(ctx context.Context) NotificationService {
	return &notification{
		ctx:    ctx,
		logger: svc.logger,
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc notification) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc notification) SendEmail(message *gomail.Message) error {
	return mail.Send(message)
}

// AttachEmailRecipients validates, resolves, formats andd attaches set of recipients to message
//
// Supports 3 input formats:
//  - <valid email>
//  - <valid email><space><name...>
//  - <userID>
// Last one is then translated into valid email + name (when/if possible)
func (svc notification) AttachEmailRecipients(message *gomail.Message, field string, recipients ...string) (err error) {
	var (
		email string
		name  string
	)

	if len(recipients) == 0 {
		return
	}

	for r, rcpt := range recipients {
		name, email = "", ""
		rcpt = strings.TrimSpace(rcpt)

		// First, get userID off the table
		if spaceAt := strings.Index(rcpt, " "); spaceAt > -1 {
			email, name = rcpt[:spaceAt], strings.TrimSpace(rcpt[spaceAt+1:])
		} else {
			email = rcpt
		}

		// Validate email here
		if !mail.IsValidAddress(email) {
			return errors.New("Invalid recipient email format")
		}

		recipients[r] = message.FormatAddress(email, name)
	}

	message.SetHeader(field, recipients...)
	return
}
