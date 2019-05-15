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

// Message create request parameters
type MessageCreate struct {
	Message   string
	ChannelID uint64 `json:",string"`
}

func NewMessageCreate() *MessageCreate {
	return &MessageCreate{}
}

func (r MessageCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["message"] = r.Message

	out["channelID"] = r.ChannelID

	return out
}

func (mReq *MessageCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

		mReq.Message = val
	}
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageCreate()

// Message executeCommand request parameters
type MessageExecuteCommand struct {
	Command   string
	ChannelID uint64 `json:",string"`
	Input     string
	Params    []string
}

func NewMessageExecuteCommand() *MessageExecuteCommand {
	return &MessageExecuteCommand{}
}

func (r MessageExecuteCommand) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["command"] = r.Command

	out["channelID"] = r.ChannelID

	out["input"] = r.Input

	out["params"] = r.Params

	return out
}

func (mReq *MessageExecuteCommand) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.Command = chi.URLParam(r, "command")
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["input"]; ok {

		mReq.Input = val
	}

	return err
}

var _ RequestFiller = NewMessageExecuteCommand()

// Message markAsRead request parameters
type MessageMarkAsRead struct {
	ChannelID         uint64 `json:",string"`
	ThreadID          uint64 `json:",string"`
	LastReadMessageID uint64 `json:",string"`
}

func NewMessageMarkAsRead() *MessageMarkAsRead {
	return &MessageMarkAsRead{}
}

func (r MessageMarkAsRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID

	out["threadID"] = r.ThreadID

	out["lastReadMessageID"] = r.LastReadMessageID

	return out
}

func (mReq *MessageMarkAsRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["threadID"]; ok {

		mReq.ThreadID = parseUInt64(val)
	}
	if val, ok := post["lastReadMessageID"]; ok {

		mReq.LastReadMessageID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewMessageMarkAsRead()

// Message edit request parameters
type MessageEdit struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
	Message   string
}

func NewMessageEdit() *MessageEdit {
	return &MessageEdit{}
}

func (r MessageEdit) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["channelID"] = r.ChannelID

	out["message"] = r.Message

	return out
}

func (mReq *MessageEdit) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["message"]; ok {

		mReq.Message = val
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

func (r MessageDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["channelID"] = r.ChannelID

	return out
}

func (mReq *MessageDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageDelete()

// Message replyCreate request parameters
type MessageReplyCreate struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
	Message   string
}

func NewMessageReplyCreate() *MessageReplyCreate {
	return &MessageReplyCreate{}
}

func (r MessageReplyCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["channelID"] = r.ChannelID

	out["message"] = r.Message

	return out
}

func (mReq *MessageReplyCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["message"]; ok {

		mReq.Message = val
	}

	return err
}

var _ RequestFiller = NewMessageReplyCreate()

// Message pinCreate request parameters
type MessagePinCreate struct {
	MessageID uint64 `json:",string"`
	ChannelID uint64 `json:",string"`
}

func NewMessagePinCreate() *MessagePinCreate {
	return &MessagePinCreate{}
}

func (r MessagePinCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["channelID"] = r.ChannelID

	return out
}

func (mReq *MessagePinCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

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

func (r MessagePinRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["channelID"] = r.ChannelID

	return out
}

func (mReq *MessagePinRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

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

func (r MessageBookmarkCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["channelID"] = r.ChannelID

	return out
}

func (mReq *MessageBookmarkCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

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

func (r MessageBookmarkRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["channelID"] = r.ChannelID

	return out
}

func (mReq *MessageBookmarkRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

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

func (r MessageReactionCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["reaction"] = r.Reaction

	out["channelID"] = r.ChannelID

	return out
}

func (mReq *MessageReactionCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.Reaction = chi.URLParam(r, "reaction")
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

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

func (r MessageReactionRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID

	out["reaction"] = r.Reaction

	out["channelID"] = r.ChannelID

	return out
}

func (mReq *MessageReactionRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	mReq.Reaction = chi.URLParam(r, "reaction")
	mReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageReactionRemove()
