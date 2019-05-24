package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `webhooks.go`, `webhooks.util.go` or `webhooks_test.go` to
	implement your API calls, helper functions and tests. The file `webhooks.go`
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

	"github.com/cortezaproject/corteza-server/messaging/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Webhooks list request parameters
type WebhooksList struct {
	ChannelID uint64 `json:",string"`
	UserID    uint64 `json:",string"`
}

func NewWebhooksList() *WebhooksList {
	return &WebhooksList{}
}

func (r WebhooksList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["userID"] = r.UserID

	return out
}

func (r *WebhooksList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["channelID"]; ok {
		r.ChannelID = parseUInt64(val)
	}
	if val, ok := get["userID"]; ok {
		r.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewWebhooksList()

// Webhooks create request parameters
type WebhooksCreate struct {
	ChannelID uint64 `json:",string"`
	Kind      types.WebhookKind
	UserID    uint64 `json:",string"`
	Trigger   string
	Url       string
	Username  string
	Avatar    *multipart.FileHeader
	AvatarURL string
}

func NewWebhooksCreate() *WebhooksCreate {
	return &WebhooksCreate{}
}

func (r WebhooksCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["kind"] = r.Kind
	out["userID"] = r.UserID
	out["trigger"] = r.Trigger
	out["url"] = r.Url
	out["username"] = r.Username
	out["avatar.size"] = r.Avatar.Size
	out["avatar.filename"] = r.Avatar.Filename

	out["avatarURL"] = r.AvatarURL

	return out
}

func (r *WebhooksCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseMultipartForm(32 << 20); err != nil {
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

	if val, ok := post["channelID"]; ok {
		r.ChannelID = parseUInt64(val)
	}
	if val, ok := post["kind"]; ok {
		r.Kind = types.WebhookKind(val)
	}
	if val, ok := post["userID"]; ok {
		r.UserID = parseUInt64(val)
	}
	if val, ok := post["trigger"]; ok {
		r.Trigger = val
	}
	if val, ok := post["url"]; ok {
		r.Url = val
	}
	if val, ok := post["username"]; ok {
		r.Username = val
	}
	if _, r.Avatar, err = req.FormFile("avatar"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	if val, ok := post["avatarURL"]; ok {
		r.AvatarURL = val
	}

	return err
}

var _ RequestFiller = NewWebhooksCreate()

// Webhooks update request parameters
type WebhooksUpdate struct {
	WebhookID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
	Kind      types.WebhookKind
	UserID    uint64 `json:",string"`
	Trigger   string
	Url       string
	Username  string
	Avatar    *multipart.FileHeader
	AvatarURL string
}

func NewWebhooksUpdate() *WebhooksUpdate {
	return &WebhooksUpdate{}
}

func (r WebhooksUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["webhookID"] = r.WebhookID
	out["channelID"] = r.ChannelID
	out["kind"] = r.Kind
	out["userID"] = r.UserID
	out["trigger"] = r.Trigger
	out["url"] = r.Url
	out["username"] = r.Username
	out["avatar.size"] = r.Avatar.Size
	out["avatar.filename"] = r.Avatar.Filename

	out["avatarURL"] = r.AvatarURL

	return out
}

func (r *WebhooksUpdate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseMultipartForm(32 << 20); err != nil {
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
	if val, ok := post["channelID"]; ok {
		r.ChannelID = parseUInt64(val)
	}
	if val, ok := post["kind"]; ok {
		r.Kind = types.WebhookKind(val)
	}
	if val, ok := post["userID"]; ok {
		r.UserID = parseUInt64(val)
	}
	if val, ok := post["trigger"]; ok {
		r.Trigger = val
	}
	if val, ok := post["url"]; ok {
		r.Url = val
	}
	if val, ok := post["username"]; ok {
		r.Username = val
	}
	if _, r.Avatar, err = req.FormFile("avatar"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	if val, ok := post["avatarURL"]; ok {
		r.AvatarURL = val
	}

	return err
}

var _ RequestFiller = NewWebhooksUpdate()

// Webhooks get request parameters
type WebhooksGet struct {
	WebhookID uint64 `json:",string"`
}

func NewWebhooksGet() *WebhooksGet {
	return &WebhooksGet{}
}

func (r WebhooksGet) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["webhookID"] = r.WebhookID

	return out
}

func (r *WebhooksGet) Fill(req *http.Request) (err error) {
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

	return err
}

var _ RequestFiller = NewWebhooksGet()

// Webhooks delete request parameters
type WebhooksDelete struct {
	WebhookID uint64 `json:",string"`
}

func NewWebhooksDelete() *WebhooksDelete {
	return &WebhooksDelete{}
}

func (r WebhooksDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["webhookID"] = r.WebhookID

	return out
}

func (r *WebhooksDelete) Fill(req *http.Request) (err error) {
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

	return err
}

var _ RequestFiller = NewWebhooksDelete()
