package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	NotificationEmailSend struct {
		// To POST parameter
		//
		// Email addresses
		To []string

		// Cc POST parameter
		//
		// Email addresses
		Cc []string

		// ReplyTo POST parameter
		//
		// Email address in reply-to field
		ReplyTo string

		// Subject POST parameter
		//
		// Email subject
		Subject string

		// Content POST parameter
		//
		// Message content
		Content sqlxTypes.JSONText

		// RemoteAttachments POST parameter
		//
		// Remote files to attach to the email
		RemoteAttachments []string
	}
)

// NewNotificationEmailSend request
func NewNotificationEmailSend() *NotificationEmailSend {
	return &NotificationEmailSend{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationEmailSend) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"to":                r.To,
		"cc":                r.Cc,
		"replyTo":           r.ReplyTo,
		"subject":           r.Subject,
		"content":           r.Content,
		"remoteAttachments": r.RemoteAttachments,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationEmailSend) GetTo() []string {
	return r.To
}

// Auditable returns all auditable/loggable parameters
func (r NotificationEmailSend) GetCc() []string {
	return r.Cc
}

// Auditable returns all auditable/loggable parameters
func (r NotificationEmailSend) GetReplyTo() string {
	return r.ReplyTo
}

// Auditable returns all auditable/loggable parameters
func (r NotificationEmailSend) GetSubject() string {
	return r.Subject
}

// Auditable returns all auditable/loggable parameters
func (r NotificationEmailSend) GetContent() sqlxTypes.JSONText {
	return r.Content
}

// Auditable returns all auditable/loggable parameters
func (r NotificationEmailSend) GetRemoteAttachments() []string {
	return r.RemoteAttachments
}

// Fill processes request and fills internal variables
func (r *NotificationEmailSend) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["replyTo"]; ok && len(val) > 0 {
				r.ReplyTo, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["subject"]; ok && len(val) > 0 {
				r.Subject, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["content"]; ok && len(val) > 0 {
				r.Content, err = payload.ParseJSONTextWithErr(val[0])
				if err != nil {
					return err
				}
			}

		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		//if val, ok := req.Form["to[]"]; ok && len(val) > 0  {
		//    r.To, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["cc[]"]; ok && len(val) > 0  {
		//    r.Cc, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["replyTo"]; ok && len(val) > 0 {
			r.ReplyTo, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["subject"]; ok && len(val) > 0 {
			r.Subject, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["content"]; ok && len(val) > 0 {
			r.Content, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["remoteAttachments[]"]; ok && len(val) > 0  {
		//    r.RemoteAttachments, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	return err
}
