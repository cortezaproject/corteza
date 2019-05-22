package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `webhooks_public.go`, `webhooks_public.util.go` or `webhooks_public_test.go` to
	implement your API calls, helper functions and tests. The file `webhooks_public.go`
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
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// WebhooksPublic delete request parameters
type WebhooksPublicDelete struct {
	WebhookID    uint64 `json:",string"`
	WebhookToken string
}

func NewWebhooksPublicDelete() *WebhooksPublicDelete {
	return &WebhooksPublicDelete{}
}

func (r WebhooksPublicDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["webhookID"] = r.WebhookID
	out["webhookToken"] = r.WebhookToken

	return out
}

func (r *WebhooksPublicDelete) Fill(req *http.Request) (err error) {
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

	r.WebhookID = parseUInt64(chi.URLParam(req, "webhookID"))
	r.WebhookToken = chi.URLParam(req, "webhookToken")

	return err
}

var _ RequestFiller = NewWebhooksPublicDelete()

// WebhooksPublic create request parameters
type WebhooksPublicCreate struct {
	Username     string
	AvatarURL    string
	Content      string
	WebhookID    uint64 `json:",string"`
	WebhookToken string
}

func NewWebhooksPublicCreate() *WebhooksPublicCreate {
	return &WebhooksPublicCreate{}
}

func (r WebhooksPublicCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["username"] = r.Username
	out["avatarURL"] = r.AvatarURL
	out["content"] = r.Content
	out["webhookID"] = r.WebhookID
	out["webhookToken"] = r.WebhookToken

	return out
}

func (r *WebhooksPublicCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := get["username"]; ok {
		r.Username = val
	}
	if val, ok := get["avatarURL"]; ok {
		r.AvatarURL = val
	}
	if val, ok := get["content"]; ok {
		r.Content = val
	}
	r.WebhookID = parseUInt64(chi.URLParam(req, "webhookID"))
	r.WebhookToken = chi.URLParam(req, "webhookToken")

	return err
}

var _ RequestFiller = NewWebhooksPublicCreate()
