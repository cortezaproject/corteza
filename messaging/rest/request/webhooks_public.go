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

// WebhooksPublicDelete request parameters
type WebhooksPublicDelete struct {
	hasWebhookID bool
	rawWebhookID string
	WebhookID    uint64 `json:",string"`

	hasWebhookToken bool
	rawWebhookToken string
	WebhookToken    string
}

// NewWebhooksPublicDelete request
func NewWebhooksPublicDelete() *WebhooksPublicDelete {
	return &WebhooksPublicDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r WebhooksPublicDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["webhookID"] = r.WebhookID
	out["webhookToken"] = r.WebhookToken

	return out
}

// Fill processes request and fills internal variables
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

	r.hasWebhookID = true
	r.rawWebhookID = chi.URLParam(req, "webhookID")
	r.WebhookID = parseUInt64(chi.URLParam(req, "webhookID"))
	r.hasWebhookToken = true
	r.rawWebhookToken = chi.URLParam(req, "webhookToken")
	r.WebhookToken = chi.URLParam(req, "webhookToken")

	return err
}

var _ RequestFiller = NewWebhooksPublicDelete()

// WebhooksPublicCreate request parameters
type WebhooksPublicCreate struct {
	hasUsername bool
	rawUsername string
	Username    string

	hasAvatarURL bool
	rawAvatarURL string
	AvatarURL    string

	hasContent bool
	rawContent string
	Content    string

	hasWebhookID bool
	rawWebhookID string
	WebhookID    uint64 `json:",string"`

	hasWebhookToken bool
	rawWebhookToken string
	WebhookToken    string
}

// NewWebhooksPublicCreate request
func NewWebhooksPublicCreate() *WebhooksPublicCreate {
	return &WebhooksPublicCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r WebhooksPublicCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["username"] = r.Username
	out["avatarURL"] = r.AvatarURL
	out["content"] = r.Content
	out["webhookID"] = r.WebhookID
	out["webhookToken"] = r.WebhookToken

	return out
}

// Fill processes request and fills internal variables
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
	r.hasWebhookID = true
	r.rawWebhookID = chi.URLParam(req, "webhookID")
	r.WebhookID = parseUInt64(chi.URLParam(req, "webhookID"))
	r.hasWebhookToken = true
	r.rawWebhookToken = chi.URLParam(req, "webhookToken")
	r.WebhookToken = chi.URLParam(req, "webhookToken")

	return err
}

var _ RequestFiller = NewWebhooksPublicCreate()

// HasWebhookID returns true if webhookID was set
func (r *WebhooksPublicDelete) HasWebhookID() bool {
	return r.hasWebhookID
}

// RawWebhookID returns raw value of webhookID parameter
func (r *WebhooksPublicDelete) RawWebhookID() string {
	return r.rawWebhookID
}

// GetWebhookID returns casted value of  webhookID parameter
func (r *WebhooksPublicDelete) GetWebhookID() uint64 {
	return r.WebhookID
}

// HasWebhookToken returns true if webhookToken was set
func (r *WebhooksPublicDelete) HasWebhookToken() bool {
	return r.hasWebhookToken
}

// RawWebhookToken returns raw value of webhookToken parameter
func (r *WebhooksPublicDelete) RawWebhookToken() string {
	return r.rawWebhookToken
}

// GetWebhookToken returns casted value of  webhookToken parameter
func (r *WebhooksPublicDelete) GetWebhookToken() string {
	return r.WebhookToken
}

// HasUsername returns true if username was set
func (r *WebhooksPublicCreate) HasUsername() bool {
	return r.hasUsername
}

// RawUsername returns raw value of username parameter
func (r *WebhooksPublicCreate) RawUsername() string {
	return r.rawUsername
}

// GetUsername returns casted value of  username parameter
func (r *WebhooksPublicCreate) GetUsername() string {
	return r.Username
}

// HasAvatarURL returns true if avatarURL was set
func (r *WebhooksPublicCreate) HasAvatarURL() bool {
	return r.hasAvatarURL
}

// RawAvatarURL returns raw value of avatarURL parameter
func (r *WebhooksPublicCreate) RawAvatarURL() string {
	return r.rawAvatarURL
}

// GetAvatarURL returns casted value of  avatarURL parameter
func (r *WebhooksPublicCreate) GetAvatarURL() string {
	return r.AvatarURL
}

// HasContent returns true if content was set
func (r *WebhooksPublicCreate) HasContent() bool {
	return r.hasContent
}

// RawContent returns raw value of content parameter
func (r *WebhooksPublicCreate) RawContent() string {
	return r.rawContent
}

// GetContent returns casted value of  content parameter
func (r *WebhooksPublicCreate) GetContent() string {
	return r.Content
}

// HasWebhookID returns true if webhookID was set
func (r *WebhooksPublicCreate) HasWebhookID() bool {
	return r.hasWebhookID
}

// RawWebhookID returns raw value of webhookID parameter
func (r *WebhooksPublicCreate) RawWebhookID() string {
	return r.rawWebhookID
}

// GetWebhookID returns casted value of  webhookID parameter
func (r *WebhooksPublicCreate) GetWebhookID() uint64 {
	return r.WebhookID
}

// HasWebhookToken returns true if webhookToken was set
func (r *WebhooksPublicCreate) HasWebhookToken() bool {
	return r.hasWebhookToken
}

// RawWebhookToken returns raw value of webhookToken parameter
func (r *WebhooksPublicCreate) RawWebhookToken() string {
	return r.rawWebhookToken
}

// GetWebhookToken returns casted value of  webhookToken parameter
func (r *WebhooksPublicCreate) GetWebhookToken() string {
	return r.WebhookToken
}
