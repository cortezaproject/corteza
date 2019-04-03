package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"

	"github.com/crusttech/crust/internal/mail"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	notification struct {
		ctx context.Context

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
		userSvc: systemService.DefaultUser,
	}).With(context.Background())
}

func (s *notification) With(ctx context.Context) NotificationService {
	return &notification{
		ctx:     ctx,
		userSvc: systemService.User(ctx),
	}
}

func (s *notification) SendEmail(message *gomail.Message) error {
	return mail.Send(message)
}

// AttachEmailRecipients validates, resolves, formats andd attaches set of recipients to message
//
// Supports 3 input formats:
//  - <valid email>
//  - <valid email><space><name...>
//  - <userID>
// Last one is then translated into valid email + name (when/if possible)
func (s *notification) AttachEmailRecipients(message *gomail.Message, field string, recipients ...string) (err error) {
	var (
		email string
		name  string
	)

	if len(recipients) == 0 {
		return
	}

	if recipients, err = s.expandUserRefs(recipients); err != nil {
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
func (s *notification) expandUserRefs(recipients []string) ([]string, error) {
	for r, rcpt := range recipients {
		// First, get userID off the table
		if userID, _ := strconv.ParseUint(rcpt, 10, 64); userID > 0 {
			if user, err := s.userSvc.FindByID(userID); err != nil {
				return nil, errors.Wrapf(err, "invalid recipient %d", userID)
			} else {
				recipients[r] = user.Email + " " + user.Name
			}
		}
	}

	return recipients, nil
}
