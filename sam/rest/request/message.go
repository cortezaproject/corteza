package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `message.go`, `message.util.go` or `message_test.go` to
	implement your API calls, helper functions and tests. The file `message.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"encoding/json"
	"github.com/crusttech/crust/internal/rbac"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}
var _ = multipart.FileHeader{}
var _ = rbac.Operation{}

// Message create request parameters
type MessageCreate struct {
	Message   string
	ChannelID uint64 `json:",string"`
}

func NewMessageCreate() *MessageCreate {
	return &MessageCreate{}
}

func (m *MessageCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	if val, ok := post["message"]; ok {

		m.Message = val
	}
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageCreate()

// Message history request parameters
type MessageHistory struct {
	LastMessageID uint64 `json:",string"`
	ChannelID     uint64 `json:",string"`
}

func NewMessageHistory() *MessageHistory {
	return &MessageHistory{}
}

func (m *MessageHistory) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	if val, ok := get["lastMessageID"]; ok {

		m.LastMessageID = parseUInt64(val)
	}
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageHistory()

// Message edit request parameters
type MessageEdit struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
	Message   string
}

func NewMessageEdit() *MessageEdit {
	return &MessageEdit{}
}

func (m *MessageEdit) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["message"]; ok {

		m.Message = val
	}

	return err
}

var _ RequestFiller = NewMessageEdit()

// Message delete request parameters
type MessageDelete struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageDelete() *MessageDelete {
	return &MessageDelete{}
}

func (m *MessageDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageDelete()

// Message replyGet request parameters
type MessageReplyGet struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageReplyGet() *MessageReplyGet {
	return &MessageReplyGet{}
}

func (m *MessageReplyGet) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageReplyGet()

// Message replyCreate request parameters
type MessageReplyCreate struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
	Message   string
}

func NewMessageReplyCreate() *MessageReplyCreate {
	return &MessageReplyCreate{}
}

func (m *MessageReplyCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["message"]; ok {

		m.Message = val
	}

	return err
}

var _ RequestFiller = NewMessageReplyCreate()

// Message markAsUnread request parameters
type MessageMarkAsUnread struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageMarkAsUnread() *MessageMarkAsUnread {
	return &MessageMarkAsUnread{}
}

func (m *MessageMarkAsUnread) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageMarkAsUnread()

// Message pinCreate request parameters
type MessagePinCreate struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessagePinCreate() *MessagePinCreate {
	return &MessagePinCreate{}
}

func (m *MessagePinCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessagePinCreate()

// Message pinRemove request parameters
type MessagePinRemove struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessagePinRemove() *MessagePinRemove {
	return &MessagePinRemove{}
}

func (m *MessagePinRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessagePinRemove()

// Message bookmarkCreate request parameters
type MessageBookmarkCreate struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageBookmarkCreate() *MessageBookmarkCreate {
	return &MessageBookmarkCreate{}
}

func (m *MessageBookmarkCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageBookmarkCreate()

// Message bookmarkRemove request parameters
type MessageBookmarkRemove struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageBookmarkRemove() *MessageBookmarkRemove {
	return &MessageBookmarkRemove{}
}

func (m *MessageBookmarkRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageBookmarkRemove()

// Message reactionCreate request parameters
type MessageReactionCreate struct {
	MessageID uint64 `json:",string"`
	Reaction  string
	ChannelID uint64 `json:",string"`
}

func NewMessageReactionCreate() *MessageReactionCreate {
	return &MessageReactionCreate{}
}

func (m *MessageReactionCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.Reaction = chi.URLParam(r, "reaction")
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageReactionCreate()

// Message reactionRemove request parameters
type MessageReactionRemove struct {
	MessageID uint64 `json:",string"`
	Reaction  string
	ChannelID uint64 `json:",string"`
}

func NewMessageReactionRemove() *MessageReactionRemove {
	return &MessageReactionRemove{}
}

func (m *MessageReactionRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.Reaction = chi.URLParam(r, "reaction")
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageReactionRemove()
