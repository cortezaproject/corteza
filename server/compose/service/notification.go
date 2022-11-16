package service

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"

	gomail "gopkg.in/mail.v2"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	httpClient "github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/cortezaproject/corteza/server/pkg/mail"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	// notification is a service that relays notification to receipients
	// currently, we only support email notifications
	//
	// @todo due to initial architectural decisions, this service landed under compose
	//       but should be moved to system
	//       Warning: API endpoints on compose should be kept so that we do not break backward compatibility)
	notification struct {
		actionlog actionlog.Recorder
		users     userFinder
	}

	notificationUserFinder interface {
		FindByID(uint64) (*systemTypes.User, error)
	}
)

func Notification(uf userFinder) *notification {
	return &notification{
		actionlog: DefaultActionlog,
		users:     uf,
	}
}

// SendEmail sends email notification
func (svc notification) SendEmail(ctx context.Context, n *types.EmailNotification) (err error) {
	var (
		aProps = &notificationActionProps{mail: n}
	)

	err = func() error {
		msg := mail.New()

		if len(n.To) == 0 {
			return NotificationErrNoRecipients()
		}

		if err = svc.procEmailRecipients(ctx, msg, "To", n.To...); err != nil {
			return err
		}

		if err = svc.procEmailRecipients(ctx, msg, "Cc", n.Cc...); err != nil {
			return err
		}

		if len(n.ReplyTo) > 0 {
			// extra check for length because ReplyTo can only hold 1 address!
			if err = svc.procEmailRecipients(ctx, msg, "ReplyTo", n.ReplyTo); err != nil {
				return err
			}
		}

		msg.SetHeader("Subject", n.Subject)

		if len(n.ContentHTML) > 0 {
			msg.SetBody("text/html", n.ContentHTML)

		}

		if len(n.ContentPlain) > 0 || len(n.ContentHTML) == 0 {
			// Make sure plain body is always set, even if empty
			msg.SetBody("text/plain", n.ContentPlain)
		}

		if err = svc.procEmailAttachments(ctx, msg, n.RemoteAttachments...); err != nil {

			return err
		}

		return mail.Send(msg)
	}()

	return svc.recordAction(ctx, aProps, NotificationActionSend, err)
}

// procEmailRecipients validates, resolves, formats and attaches set of recipients to message
//
// Supports 3 input formats:
//  - <valid email>
//  - <valid email><space><name...>
//  - <userID>
// Last one is then translated into valid email + name (when/if possible)
func (svc notification) procEmailRecipients(ctx context.Context, m *gomail.Message, field string, rr ...string) (err error) {
	var (
		email string
		name  string
	)

	if len(rr) == 0 {
		return
	}

	for r, rcpt := range rr {
		aProps := &notificationActionProps{recipient: rcpt}

		name, email = "", ""
		rcpt = strings.TrimSpace(rcpt)

		if userID, err := strconv.ParseUint(rcpt, 10, 64); err == nil && userID > 0 {
			// proc <user ID>
			if user, err := svc.users.FindByID(ctx, userID); err != nil {
				return NotificationErrFailedToLoadUser(aProps).Wrap(err)
			} else {
				email = user.Email
				name = user.Name
			}

		} else if spaceAt := strings.Index(rcpt, " "); spaceAt > -1 {
			// proc <email> <name> ("foo@bar.baz foo baz")
			email, name = rcpt[:spaceAt], strings.TrimSpace(rcpt[spaceAt+1:])
		} else {
			// proc <email>
			email = rcpt
		}

		// Validate email here
		if !mail.IsValidAddress(email) {
			return NotificationErrInvalidReceipientFormat(aProps)
		}

		rr[r] = m.FormatAddress(email, name)
	}

	m.SetHeader(field, rr...)
	return nil
}

// procEmailAttachments treats given strings (URLs) as remote objects, downloads and attaches them to
// the message
//
// This could/should be easily extended to support data URLs as well
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URIs
func (svc notification) procEmailAttachments(ctx context.Context, message *gomail.Message, aa ...string) error {
	var (
		// threading safely (when multiple objects are to be attached)
		l  = &sync.Mutex{}
		wg = &sync.WaitGroup{}

		client, err = httpClient.New(&httpClient.Config{
			Timeout: 10,
		})
	)

	if err != nil {
		return err
	}

	get := func(url string) error {
		aProps := &notificationActionProps{attachmentURL: url}

		req, err := client.Get(url)
		if err != nil {
			return err
		}

		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			return NotificationErrFailedToDownloadAttachment(aProps).Wrap(err)
		}

		aProps.setAttachmentType(resp.Header.Get("Content-Type"))
		aProps.setAttachmentSize(resp.ContentLength)

		if resp.StatusCode != http.StatusOK {
			return NotificationErrFailedToDownloadAttachment(aProps).
				Wrap(fmt.Errorf("unexpected HTTP status: %s", resp.Status))
		}

		l.Lock()
		defer l.Unlock()
		message.AttachReader(path.Base(req.URL.Path), resp.Body)
		return nil
	}

	for _, url := range aa {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			_ = svc.recordAction(
				ctx,
				&notificationActionProps{attachmentURL: url},
				NotificationActionAttachmentDownload,
				get(url),
			)
		}(url)
	}

	wg.Wait()
	return nil
}
