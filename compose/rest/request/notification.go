package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `notification.go`, `notification.util.go` or `notification_test.go` to
	implement your API calls, helper functions and tests. The file `notification.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"io"
	"strings"

	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	sqlxTypes "github.com/jmoiron/sqlx/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// NotificationEmailSend request parameters
type NotificationEmailSend struct {
	hasTo bool
	rawTo []string
	To    []string

	hasCc bool
	rawCc []string
	Cc    []string

	hasReplyTo bool
	rawReplyTo string
	ReplyTo    string

	hasSubject bool
	rawSubject string
	Subject    string

	hasContent bool
	rawContent string
	Content    sqlxTypes.JSONText

	hasRemoteAttachments bool
	rawRemoteAttachments []string
	RemoteAttachments    []string
}

// NewNotificationEmailSend request
func NewNotificationEmailSend() *NotificationEmailSend {
	return &NotificationEmailSend{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationEmailSend) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["to"] = r.To
	out["cc"] = r.Cc
	out["replyTo"] = r.ReplyTo
	out["subject"] = r.Subject
	out["content"] = r.Content
	out["remoteAttachments"] = r.RemoteAttachments

	return out
}

// Fill processes request and fills internal variables
func (r *NotificationEmailSend) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := req.Form["to"]; ok {
		r.hasTo = true
		r.rawTo = val
		r.To = parseStrings(val)
	}

	if val, ok := req.Form["cc"]; ok {
		r.hasCc = true
		r.rawCc = val
		r.Cc = parseStrings(val)
	}

	if val, ok := post["replyTo"]; ok {
		r.hasReplyTo = true
		r.rawReplyTo = val
		r.ReplyTo = val
	}
	if val, ok := post["subject"]; ok {
		r.hasSubject = true
		r.rawSubject = val
		r.Subject = val
	}
	if val, ok := post["content"]; ok {
		r.hasContent = true
		r.rawContent = val

		if r.Content, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	if val, ok := req.Form["remoteAttachments"]; ok {
		r.hasRemoteAttachments = true
		r.rawRemoteAttachments = val
		r.RemoteAttachments = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewNotificationEmailSend()

// HasTo returns true if to was set
func (r *NotificationEmailSend) HasTo() bool {
	return r.hasTo
}

// RawTo returns raw value of to parameter
func (r *NotificationEmailSend) RawTo() []string {
	return r.rawTo
}

// GetTo returns casted value of  to parameter
func (r *NotificationEmailSend) GetTo() []string {
	return r.To
}

// HasCc returns true if cc was set
func (r *NotificationEmailSend) HasCc() bool {
	return r.hasCc
}

// RawCc returns raw value of cc parameter
func (r *NotificationEmailSend) RawCc() []string {
	return r.rawCc
}

// GetCc returns casted value of  cc parameter
func (r *NotificationEmailSend) GetCc() []string {
	return r.Cc
}

// HasReplyTo returns true if replyTo was set
func (r *NotificationEmailSend) HasReplyTo() bool {
	return r.hasReplyTo
}

// RawReplyTo returns raw value of replyTo parameter
func (r *NotificationEmailSend) RawReplyTo() string {
	return r.rawReplyTo
}

// GetReplyTo returns casted value of  replyTo parameter
func (r *NotificationEmailSend) GetReplyTo() string {
	return r.ReplyTo
}

// HasSubject returns true if subject was set
func (r *NotificationEmailSend) HasSubject() bool {
	return r.hasSubject
}

// RawSubject returns raw value of subject parameter
func (r *NotificationEmailSend) RawSubject() string {
	return r.rawSubject
}

// GetSubject returns casted value of  subject parameter
func (r *NotificationEmailSend) GetSubject() string {
	return r.Subject
}

// HasContent returns true if content was set
func (r *NotificationEmailSend) HasContent() bool {
	return r.hasContent
}

// RawContent returns raw value of content parameter
func (r *NotificationEmailSend) RawContent() string {
	return r.rawContent
}

// GetContent returns casted value of  content parameter
func (r *NotificationEmailSend) GetContent() sqlxTypes.JSONText {
	return r.Content
}

// HasRemoteAttachments returns true if remoteAttachments was set
func (r *NotificationEmailSend) HasRemoteAttachments() bool {
	return r.hasRemoteAttachments
}

// RawRemoteAttachments returns raw value of remoteAttachments parameter
func (r *NotificationEmailSend) RawRemoteAttachments() []string {
	return r.rawRemoteAttachments
}

// GetRemoteAttachments returns casted value of  remoteAttachments parameter
func (r *NotificationEmailSend) GetRemoteAttachments() []string {
	return r.RemoteAttachments
}
