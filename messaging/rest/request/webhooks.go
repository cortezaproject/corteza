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

// WebhooksList request parameters
type WebhooksList struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewWebhooksList request
func NewWebhooksList() *WebhooksList {
	return &WebhooksList{}
}

// Auditable returns all auditable/loggable parameters
func (r WebhooksList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
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
		r.hasChannelID = true
		r.rawChannelID = val
		r.ChannelID = parseUInt64(val)
	}
	if val, ok := get["userID"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewWebhooksList()

// WebhooksCreate request parameters
type WebhooksCreate struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasKind bool
	rawKind string
	Kind    types.WebhookKind

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`

	hasTrigger bool
	rawTrigger string
	Trigger    string

	hasUrl bool
	rawUrl string
	Url    string

	hasUsername bool
	rawUsername string
	Username    string

	hasAvatar bool
	rawAvatar string
	Avatar    *multipart.FileHeader

	hasAvatarURL bool
	rawAvatarURL string
	AvatarURL    string
}

// NewWebhooksCreate request
func NewWebhooksCreate() *WebhooksCreate {
	return &WebhooksCreate{}
}

// Auditable returns all auditable/loggable parameters
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

// Fill processes request and fills internal variables
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
		r.hasChannelID = true
		r.rawChannelID = val
		r.ChannelID = parseUInt64(val)
	}
	if val, ok := post["kind"]; ok {
		r.hasKind = true
		r.rawKind = val
		r.Kind = types.WebhookKind(val)
	}
	if val, ok := post["userID"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseUInt64(val)
	}
	if val, ok := post["trigger"]; ok {
		r.hasTrigger = true
		r.rawTrigger = val
		r.Trigger = val
	}
	if val, ok := post["url"]; ok {
		r.hasUrl = true
		r.rawUrl = val
		r.Url = val
	}
	if val, ok := post["username"]; ok {
		r.hasUsername = true
		r.rawUsername = val
		r.Username = val
	}
	if _, r.Avatar, err = req.FormFile("avatar"); err != nil {
		return errors.Wrap(err, "error processing uploaded file")
	}

	if val, ok := post["avatarURL"]; ok {
		r.hasAvatarURL = true
		r.rawAvatarURL = val
		r.AvatarURL = val
	}

	return err
}

var _ RequestFiller = NewWebhooksCreate()

// WebhooksUpdate request parameters
type WebhooksUpdate struct {
	hasWebhookID bool
	rawWebhookID string
	WebhookID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasKind bool
	rawKind string
	Kind    types.WebhookKind

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`

	hasTrigger bool
	rawTrigger string
	Trigger    string

	hasUrl bool
	rawUrl string
	Url    string

	hasUsername bool
	rawUsername string
	Username    string

	hasAvatar bool
	rawAvatar string
	Avatar    *multipart.FileHeader

	hasAvatarURL bool
	rawAvatarURL string
	AvatarURL    string
}

// NewWebhooksUpdate request
func NewWebhooksUpdate() *WebhooksUpdate {
	return &WebhooksUpdate{}
}

// Auditable returns all auditable/loggable parameters
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

// Fill processes request and fills internal variables
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

	r.hasWebhookID = true
	r.rawWebhookID = chi.URLParam(req, "webhookID")
	r.WebhookID = parseUInt64(chi.URLParam(req, "webhookID"))
	if val, ok := post["channelID"]; ok {
		r.hasChannelID = true
		r.rawChannelID = val
		r.ChannelID = parseUInt64(val)
	}
	if val, ok := post["kind"]; ok {
		r.hasKind = true
		r.rawKind = val
		r.Kind = types.WebhookKind(val)
	}
	if val, ok := post["userID"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseUInt64(val)
	}
	if val, ok := post["trigger"]; ok {
		r.hasTrigger = true
		r.rawTrigger = val
		r.Trigger = val
	}
	if val, ok := post["url"]; ok {
		r.hasUrl = true
		r.rawUrl = val
		r.Url = val
	}
	if val, ok := post["username"]; ok {
		r.hasUsername = true
		r.rawUsername = val
		r.Username = val
	}
	if _, r.Avatar, err = req.FormFile("avatar"); err != nil {
		return errors.Wrap(err, "error processing uploaded file")
	}

	if val, ok := post["avatarURL"]; ok {
		r.hasAvatarURL = true
		r.rawAvatarURL = val
		r.AvatarURL = val
	}

	return err
}

var _ RequestFiller = NewWebhooksUpdate()

// WebhooksGet request parameters
type WebhooksGet struct {
	hasWebhookID bool
	rawWebhookID string
	WebhookID    uint64 `json:",string"`
}

// NewWebhooksGet request
func NewWebhooksGet() *WebhooksGet {
	return &WebhooksGet{}
}

// Auditable returns all auditable/loggable parameters
func (r WebhooksGet) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["webhookID"] = r.WebhookID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasWebhookID = true
	r.rawWebhookID = chi.URLParam(req, "webhookID")
	r.WebhookID = parseUInt64(chi.URLParam(req, "webhookID"))

	return err
}

var _ RequestFiller = NewWebhooksGet()

// WebhooksDelete request parameters
type WebhooksDelete struct {
	hasWebhookID bool
	rawWebhookID string
	WebhookID    uint64 `json:",string"`
}

// NewWebhooksDelete request
func NewWebhooksDelete() *WebhooksDelete {
	return &WebhooksDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r WebhooksDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["webhookID"] = r.WebhookID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasWebhookID = true
	r.rawWebhookID = chi.URLParam(req, "webhookID")
	r.WebhookID = parseUInt64(chi.URLParam(req, "webhookID"))

	return err
}

var _ RequestFiller = NewWebhooksDelete()

// HasChannelID returns true if channelID was set
func (r *WebhooksList) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *WebhooksList) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *WebhooksList) GetChannelID() uint64 {
	return r.ChannelID
}

// HasUserID returns true if userID was set
func (r *WebhooksList) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *WebhooksList) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *WebhooksList) GetUserID() uint64 {
	return r.UserID
}

// HasChannelID returns true if channelID was set
func (r *WebhooksCreate) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *WebhooksCreate) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *WebhooksCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// HasKind returns true if kind was set
func (r *WebhooksCreate) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *WebhooksCreate) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *WebhooksCreate) GetKind() types.WebhookKind {
	return r.Kind
}

// HasUserID returns true if userID was set
func (r *WebhooksCreate) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *WebhooksCreate) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *WebhooksCreate) GetUserID() uint64 {
	return r.UserID
}

// HasTrigger returns true if trigger was set
func (r *WebhooksCreate) HasTrigger() bool {
	return r.hasTrigger
}

// RawTrigger returns raw value of trigger parameter
func (r *WebhooksCreate) RawTrigger() string {
	return r.rawTrigger
}

// GetTrigger returns casted value of  trigger parameter
func (r *WebhooksCreate) GetTrigger() string {
	return r.Trigger
}

// HasUrl returns true if url was set
func (r *WebhooksCreate) HasUrl() bool {
	return r.hasUrl
}

// RawUrl returns raw value of url parameter
func (r *WebhooksCreate) RawUrl() string {
	return r.rawUrl
}

// GetUrl returns casted value of  url parameter
func (r *WebhooksCreate) GetUrl() string {
	return r.Url
}

// HasUsername returns true if username was set
func (r *WebhooksCreate) HasUsername() bool {
	return r.hasUsername
}

// RawUsername returns raw value of username parameter
func (r *WebhooksCreate) RawUsername() string {
	return r.rawUsername
}

// GetUsername returns casted value of  username parameter
func (r *WebhooksCreate) GetUsername() string {
	return r.Username
}

// HasAvatar returns true if avatar was set
func (r *WebhooksCreate) HasAvatar() bool {
	return r.hasAvatar
}

// RawAvatar returns raw value of avatar parameter
func (r *WebhooksCreate) RawAvatar() string {
	return r.rawAvatar
}

// GetAvatar returns casted value of  avatar parameter
func (r *WebhooksCreate) GetAvatar() *multipart.FileHeader {
	return r.Avatar
}

// HasAvatarURL returns true if avatarURL was set
func (r *WebhooksCreate) HasAvatarURL() bool {
	return r.hasAvatarURL
}

// RawAvatarURL returns raw value of avatarURL parameter
func (r *WebhooksCreate) RawAvatarURL() string {
	return r.rawAvatarURL
}

// GetAvatarURL returns casted value of  avatarURL parameter
func (r *WebhooksCreate) GetAvatarURL() string {
	return r.AvatarURL
}

// HasWebhookID returns true if webhookID was set
func (r *WebhooksUpdate) HasWebhookID() bool {
	return r.hasWebhookID
}

// RawWebhookID returns raw value of webhookID parameter
func (r *WebhooksUpdate) RawWebhookID() string {
	return r.rawWebhookID
}

// GetWebhookID returns casted value of  webhookID parameter
func (r *WebhooksUpdate) GetWebhookID() uint64 {
	return r.WebhookID
}

// HasChannelID returns true if channelID was set
func (r *WebhooksUpdate) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *WebhooksUpdate) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *WebhooksUpdate) GetChannelID() uint64 {
	return r.ChannelID
}

// HasKind returns true if kind was set
func (r *WebhooksUpdate) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *WebhooksUpdate) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *WebhooksUpdate) GetKind() types.WebhookKind {
	return r.Kind
}

// HasUserID returns true if userID was set
func (r *WebhooksUpdate) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *WebhooksUpdate) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *WebhooksUpdate) GetUserID() uint64 {
	return r.UserID
}

// HasTrigger returns true if trigger was set
func (r *WebhooksUpdate) HasTrigger() bool {
	return r.hasTrigger
}

// RawTrigger returns raw value of trigger parameter
func (r *WebhooksUpdate) RawTrigger() string {
	return r.rawTrigger
}

// GetTrigger returns casted value of  trigger parameter
func (r *WebhooksUpdate) GetTrigger() string {
	return r.Trigger
}

// HasUrl returns true if url was set
func (r *WebhooksUpdate) HasUrl() bool {
	return r.hasUrl
}

// RawUrl returns raw value of url parameter
func (r *WebhooksUpdate) RawUrl() string {
	return r.rawUrl
}

// GetUrl returns casted value of  url parameter
func (r *WebhooksUpdate) GetUrl() string {
	return r.Url
}

// HasUsername returns true if username was set
func (r *WebhooksUpdate) HasUsername() bool {
	return r.hasUsername
}

// RawUsername returns raw value of username parameter
func (r *WebhooksUpdate) RawUsername() string {
	return r.rawUsername
}

// GetUsername returns casted value of  username parameter
func (r *WebhooksUpdate) GetUsername() string {
	return r.Username
}

// HasAvatar returns true if avatar was set
func (r *WebhooksUpdate) HasAvatar() bool {
	return r.hasAvatar
}

// RawAvatar returns raw value of avatar parameter
func (r *WebhooksUpdate) RawAvatar() string {
	return r.rawAvatar
}

// GetAvatar returns casted value of  avatar parameter
func (r *WebhooksUpdate) GetAvatar() *multipart.FileHeader {
	return r.Avatar
}

// HasAvatarURL returns true if avatarURL was set
func (r *WebhooksUpdate) HasAvatarURL() bool {
	return r.hasAvatarURL
}

// RawAvatarURL returns raw value of avatarURL parameter
func (r *WebhooksUpdate) RawAvatarURL() string {
	return r.rawAvatarURL
}

// GetAvatarURL returns casted value of  avatarURL parameter
func (r *WebhooksUpdate) GetAvatarURL() string {
	return r.AvatarURL
}

// HasWebhookID returns true if webhookID was set
func (r *WebhooksGet) HasWebhookID() bool {
	return r.hasWebhookID
}

// RawWebhookID returns raw value of webhookID parameter
func (r *WebhooksGet) RawWebhookID() string {
	return r.rawWebhookID
}

// GetWebhookID returns casted value of  webhookID parameter
func (r *WebhooksGet) GetWebhookID() uint64 {
	return r.WebhookID
}

// HasWebhookID returns true if webhookID was set
func (r *WebhooksDelete) HasWebhookID() bool {
	return r.hasWebhookID
}

// RawWebhookID returns raw value of webhookID parameter
func (r *WebhooksDelete) RawWebhookID() string {
	return r.rawWebhookID
}

// GetWebhookID returns casted value of  webhookID parameter
func (r *WebhooksDelete) GetWebhookID() uint64 {
	return r.WebhookID
}
