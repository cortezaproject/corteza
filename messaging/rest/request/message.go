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

	out["message"] = "*masked*sensitive*data*"

	out["channelID"] = r.ChannelID

	return out
}

func (r *MessageCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["message"]; ok {
		r.Message = val
	}
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

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

func (r *MessageExecuteCommand) Fill(req *http.Request) (err error) {
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

	r.Command = chi.URLParam(req, "command")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["input"]; ok {
		r.Input = val
	}

	if val, ok := req.Form["params"]; ok {
		r.Params = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewMessageExecuteCommand()

// Message markAsRead request parameters
type MessageMarkAsRead struct {
	ThreadID          uint64 `json:",string"`
	LastReadMessageID uint64 `json:",string"`
	ChannelID         uint64 `json:",string"`
}

func NewMessageMarkAsRead() *MessageMarkAsRead {
	return &MessageMarkAsRead{}
}

func (r MessageMarkAsRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["threadID"] = r.ThreadID
	out["lastReadMessageID"] = r.LastReadMessageID
	out["channelID"] = r.ChannelID

	return out
}

func (r *MessageMarkAsRead) Fill(req *http.Request) (err error) {
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

	if val, ok := get["threadID"]; ok {
		r.ThreadID = parseUInt64(val)
	}
	if val, ok := get["lastReadMessageID"]; ok {
		r.LastReadMessageID = parseUInt64(val)
	}
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

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
	out["message"] = "*masked*sensitive*data*"

	return out
}

func (r *MessageEdit) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["message"]; ok {
		r.Message = val
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

func (r *MessageDelete) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

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
	out["message"] = "*masked*sensitive*data*"

	return out
}

func (r *MessageReplyCreate) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["message"]; ok {
		r.Message = val
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

func (r *MessagePinCreate) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

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

func (r *MessagePinRemove) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

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

func (r *MessageBookmarkCreate) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

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

func (r *MessageBookmarkRemove) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

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

func (r *MessageReactionCreate) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.Reaction = chi.URLParam(req, "reaction")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

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

func (r *MessageReactionRemove) Fill(req *http.Request) (err error) {
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

	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.Reaction = chi.URLParam(req, "reaction")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessageReactionRemove()
