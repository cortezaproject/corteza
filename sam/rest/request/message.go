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

// Message getReplies request parameters
type MessageGetReplies struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageGetReplies() *MessageGetReplies {
	return &MessageGetReplies{}
}

func (m *MessageGetReplies) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewMessageGetReplies()

// Message createReply request parameters
type MessageCreateReply struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
	Message   string
}

func NewMessageCreateReply() *MessageCreateReply {
	return &MessageCreateReply{}
}

func (m *MessageCreateReply) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewMessageCreateReply()

// Message pin request parameters
type MessagePin struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessagePin() *MessagePin {
	return &MessagePin{}
}

func (m *MessagePin) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewMessagePin()

// Message removePin request parameters
type MessageRemovePin struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageRemovePin() *MessageRemovePin {
	return &MessageRemovePin{}
}

func (m *MessageRemovePin) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewMessageRemovePin()

// Message bookmark request parameters
type MessageBookmark struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageBookmark() *MessageBookmark {
	return &MessageBookmark{}
}

func (m *MessageBookmark) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewMessageBookmark()

// Message removeBookmark request parameters
type MessageRemoveBookmark struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessageRemoveBookmark() *MessageRemoveBookmark {
	return &MessageRemoveBookmark{}
}

func (m *MessageRemoveBookmark) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewMessageRemoveBookmark()

// Message react request parameters
type MessageReact struct {
	MessageID uint64 `json:",string"`
	Reaction  string
	ChannelID uint64 `json:",string"`
}

func NewMessageReact() *MessageReact {
	return &MessageReact{}
}

func (m *MessageReact) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewMessageReact()

// Message removeReaction request parameters
type MessageRemoveReaction struct {
	MessageID uint64 `json:",string"`
	Reaction  string
	ChannelID uint64 `json:",string"`
}

func NewMessageRemoveReaction() *MessageRemoveReaction {
	return &MessageRemoveReaction{}
}

func (m *MessageRemoveReaction) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewMessageRemoveReaction()
