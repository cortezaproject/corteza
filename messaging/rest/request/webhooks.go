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

	"github.com/crusttech/crust/messaging/types"
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

func (wReq *WebhooksList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

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

	if val, ok := get["channelID"]; ok {

		wReq.ChannelID = parseUInt64(val)
	}
	if val, ok := get["userID"]; ok {

		wReq.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewWebhooksList()

// Webhooks create request parameters
type WebhooksCreate struct {
	ChannelID uint64 `json:",string"`
	Kind      types.WebhookKind
	Trigger   string
	Url       string
	Username  string
	Avatar    *multipart.FileHeader
	AvatarURL string
}

func NewWebhooksCreate() *WebhooksCreate {
	return &WebhooksCreate{}
}

func (wReq *WebhooksCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseMultipartForm(32 << 20); err != nil {
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

	if val, ok := post["channelID"]; ok {

		wReq.ChannelID = parseUInt64(val)
	}
	if val, ok := post["kind"]; ok {

		wReq.Kind = types.WebhookKind(val)
	}
	if val, ok := post["trigger"]; ok {

		wReq.Trigger = val
	}
	if val, ok := post["url"]; ok {

		wReq.Url = val
	}
	if val, ok := post["username"]; ok {

		wReq.Username = val
	}
	if _, wReq.Avatar, err = r.FormFile("avatar"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	if val, ok := post["avatarURL"]; ok {

		wReq.AvatarURL = val
	}

	return err
}

var _ RequestFiller = NewWebhooksCreate()

// Webhooks update request parameters
type WebhooksUpdate struct {
	WebhookID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
	Kind      types.WebhookKind
	Trigger   string
	Url       string
	Username  string
	Avatar    *multipart.FileHeader
	AvatarURL string
}

func NewWebhooksUpdate() *WebhooksUpdate {
	return &WebhooksUpdate{}
}

func (wReq *WebhooksUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseMultipartForm(32 << 20); err != nil {
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

	wReq.WebhookID = parseUInt64(chi.URLParam(r, "webhookID"))
	if val, ok := post["channelID"]; ok {

		wReq.ChannelID = parseUInt64(val)
	}
	if val, ok := post["kind"]; ok {

		wReq.Kind = types.WebhookKind(val)
	}
	if val, ok := post["trigger"]; ok {

		wReq.Trigger = val
	}
	if val, ok := post["url"]; ok {

		wReq.Url = val
	}
	if val, ok := post["username"]; ok {

		wReq.Username = val
	}
	if _, wReq.Avatar, err = r.FormFile("avatar"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	if val, ok := post["avatarURL"]; ok {

		wReq.AvatarURL = val
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

func (wReq *WebhooksGet) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

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

	wReq.WebhookID = parseUInt64(chi.URLParam(r, "webhookID"))

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

func (wReq *WebhooksDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(wReq)

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

	wReq.WebhookID = parseUInt64(chi.URLParam(r, "webhookID"))

	return err
}

var _ RequestFiller = NewWebhooksDelete()
