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

func (nReq *NotificationEmailSend) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(nReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := post["replyTo"]; ok {

		nReq.ReplyTo = val
	}
	if val, ok := post["subject "]; ok {

		nReq.Subject = val
	}
	if val, ok := post["content"]; ok {

		if nReq.Content, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewNotificationEmailSend()
