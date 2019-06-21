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

// Notification email/send request parameters
type NotificationEmailSend struct {
	To      []string
	Cc      []string
	ReplyTo string
	Subject string
	Content sqlxTypes.JSONText
}

func NewNotificationEmailSend() *NotificationEmailSend {
	return &NotificationEmailSend{}
}

func (r NotificationEmailSend) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["to"] = r.To
	out["cc"] = r.Cc
	out["replyTo"] = r.ReplyTo
	out["subject "] = r.Subject
	out["content"] = r.Content

	return out
}

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
		r.To = parseStrings(val)
	}

	if val, ok := req.Form["cc"]; ok {
		r.Cc = parseStrings(val)
	}

	if val, ok := post["replyTo"]; ok {
		r.ReplyTo = val
	}
	if val, ok := post["subject "]; ok {
		r.Subject = val
	}
	if val, ok := post["content"]; ok {

		if r.Content, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewNotificationEmailSend()
