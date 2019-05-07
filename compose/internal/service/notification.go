package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gomail "gopkg.in/mail.v2"

	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/mail"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	notification struct {
		ctx    context.Context
		logger *zap.Logger

		userSvc systemService.UserService
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

		userSvc: systemService.DefaultUser,
	}).With(context.Background())
}

func (svc notification) With(ctx context.Context) NotificationService {
	return &notification{
		ctx:    ctx,
		logger: svc.logger,

		userSvc: systemService.User(ctx),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc notification) log(fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
}

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

	if recipients, err = svc.expandUserRefs(recipients); err != nil {
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

// Expands references to users (strings as numeric uint64)
//
// This func is extracted to make testing/mocking mocking
func (svc notification) expandUserRefs(recipients []string) ([]string, error) {
	for r, rcpt := range recipients {
		// First, get userID off the table
		if userID, _ := strconv.ParseUint(rcpt, 10, 64); userID > 0 {
			if user, err := svc.userSvc.FindByID(userID); err != nil {
				return nil, errors.Wrapf(err, "invalid recipient %d", userID)
			} else {
				recipients[r] = user.Email + " " + user.Name
			}
		}
	}

	return recipients, nil
}
