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

// MessageCreate request parameters
type MessageCreate struct {
	hasMessage bool
	rawMessage string
	Message    string

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessageCreate request
func NewMessageCreate() *MessageCreate {
	return &MessageCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["message"] = "*masked*sensitive*data*"

	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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
		r.hasMessage = true
		r.rawMessage = val
		r.Message = val
	}
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessageCreate()

// MessageExecuteCommand request parameters
type MessageExecuteCommand struct {
	hasCommand bool
	rawCommand string
	Command    string

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasInput bool
	rawInput string
	Input    string

	hasParams bool
	rawParams []string
	Params    []string
}

// NewMessageExecuteCommand request
func NewMessageExecuteCommand() *MessageExecuteCommand {
	return &MessageExecuteCommand{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageExecuteCommand) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["command"] = r.Command
	out["channelID"] = r.ChannelID
	out["input"] = r.Input
	out["params"] = r.Params

	return out
}

// Fill processes request and fills internal variables
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

	r.hasCommand = true
	r.rawCommand = chi.URLParam(req, "command")
	r.Command = chi.URLParam(req, "command")
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["input"]; ok {
		r.hasInput = true
		r.rawInput = val
		r.Input = val
	}

	if val, ok := req.Form["params"]; ok {
		r.hasParams = true
		r.rawParams = val
		r.Params = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewMessageExecuteCommand()

// MessageMarkAsRead request parameters
type MessageMarkAsRead struct {
	hasThreadID bool
	rawThreadID string
	ThreadID    uint64 `json:",string"`

	hasLastReadMessageID bool
	rawLastReadMessageID string
	LastReadMessageID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessageMarkAsRead request
func NewMessageMarkAsRead() *MessageMarkAsRead {
	return &MessageMarkAsRead{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageMarkAsRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["threadID"] = r.ThreadID
	out["lastReadMessageID"] = r.LastReadMessageID
	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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
		r.hasThreadID = true
		r.rawThreadID = val
		r.ThreadID = parseUInt64(val)
	}
	if val, ok := get["lastReadMessageID"]; ok {
		r.hasLastReadMessageID = true
		r.rawLastReadMessageID = val
		r.LastReadMessageID = parseUInt64(val)
	}
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessageMarkAsRead()

// MessageEdit request parameters
type MessageEdit struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasMessage bool
	rawMessage string
	Message    string
}

// NewMessageEdit request
func NewMessageEdit() *MessageEdit {
	return &MessageEdit{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageEdit) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["channelID"] = r.ChannelID
	out["message"] = "*masked*sensitive*data*"

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["message"]; ok {
		r.hasMessage = true
		r.rawMessage = val
		r.Message = val
	}

	return err
}

var _ RequestFiller = NewMessageEdit()

// MessageDelete request parameters
type MessageDelete struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessageDelete request
func NewMessageDelete() *MessageDelete {
	return &MessageDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessageDelete()

// MessageReplyCreate request parameters
type MessageReplyCreate struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasMessage bool
	rawMessage string
	Message    string
}

// NewMessageReplyCreate request
func NewMessageReplyCreate() *MessageReplyCreate {
	return &MessageReplyCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReplyCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["channelID"] = r.ChannelID
	out["message"] = "*masked*sensitive*data*"

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["message"]; ok {
		r.hasMessage = true
		r.rawMessage = val
		r.Message = val
	}

	return err
}

var _ RequestFiller = NewMessageReplyCreate()

// MessagePinCreate request parameters
type MessagePinCreate struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessagePinCreate request
func NewMessagePinCreate() *MessagePinCreate {
	return &MessagePinCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessagePinCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessagePinCreate()

// MessagePinRemove request parameters
type MessagePinRemove struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessagePinRemove request
func NewMessagePinRemove() *MessagePinRemove {
	return &MessagePinRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r MessagePinRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessagePinRemove()

// MessageBookmarkCreate request parameters
type MessageBookmarkCreate struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessageBookmarkCreate request
func NewMessageBookmarkCreate() *MessageBookmarkCreate {
	return &MessageBookmarkCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageBookmarkCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessageBookmarkCreate()

// MessageBookmarkRemove request parameters
type MessageBookmarkRemove struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessageBookmarkRemove request
func NewMessageBookmarkRemove() *MessageBookmarkRemove {
	return &MessageBookmarkRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageBookmarkRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessageBookmarkRemove()

// MessageReactionCreate request parameters
type MessageReactionCreate struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasReaction bool
	rawReaction string
	Reaction    string

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessageReactionCreate request
func NewMessageReactionCreate() *MessageReactionCreate {
	return &MessageReactionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["reaction"] = r.Reaction
	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasReaction = true
	r.rawReaction = chi.URLParam(req, "reaction")
	r.Reaction = chi.URLParam(req, "reaction")
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessageReactionCreate()

// MessageReactionRemove request parameters
type MessageReactionRemove struct {
	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasReaction bool
	rawReaction string
	Reaction    string

	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewMessageReactionRemove request
func NewMessageReactionRemove() *MessageReactionRemove {
	return &MessageReactionRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["messageID"] = r.MessageID
	out["reaction"] = r.Reaction
	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasMessageID = true
	r.rawMessageID = chi.URLParam(req, "messageID")
	r.MessageID = parseUInt64(chi.URLParam(req, "messageID"))
	r.hasReaction = true
	r.rawReaction = chi.URLParam(req, "reaction")
	r.Reaction = chi.URLParam(req, "reaction")
	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewMessageReactionRemove()

// HasMessage returns true if message was set
func (r *MessageCreate) HasMessage() bool {
	return r.hasMessage
}

// RawMessage returns raw value of message parameter
func (r *MessageCreate) RawMessage() string {
	return r.rawMessage
}

// GetMessage returns casted value of  message parameter
func (r *MessageCreate) GetMessage() string {
	return r.Message
}

// HasChannelID returns true if channelID was set
func (r *MessageCreate) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageCreate) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// HasCommand returns true if command was set
func (r *MessageExecuteCommand) HasCommand() bool {
	return r.hasCommand
}

// RawCommand returns raw value of command parameter
func (r *MessageExecuteCommand) RawCommand() string {
	return r.rawCommand
}

// GetCommand returns casted value of  command parameter
func (r *MessageExecuteCommand) GetCommand() string {
	return r.Command
}

// HasChannelID returns true if channelID was set
func (r *MessageExecuteCommand) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageExecuteCommand) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageExecuteCommand) GetChannelID() uint64 {
	return r.ChannelID
}

// HasInput returns true if input was set
func (r *MessageExecuteCommand) HasInput() bool {
	return r.hasInput
}

// RawInput returns raw value of input parameter
func (r *MessageExecuteCommand) RawInput() string {
	return r.rawInput
}

// GetInput returns casted value of  input parameter
func (r *MessageExecuteCommand) GetInput() string {
	return r.Input
}

// HasParams returns true if params was set
func (r *MessageExecuteCommand) HasParams() bool {
	return r.hasParams
}

// RawParams returns raw value of params parameter
func (r *MessageExecuteCommand) RawParams() []string {
	return r.rawParams
}

// GetParams returns casted value of  params parameter
func (r *MessageExecuteCommand) GetParams() []string {
	return r.Params
}

// HasThreadID returns true if threadID was set
func (r *MessageMarkAsRead) HasThreadID() bool {
	return r.hasThreadID
}

// RawThreadID returns raw value of threadID parameter
func (r *MessageMarkAsRead) RawThreadID() string {
	return r.rawThreadID
}

// GetThreadID returns casted value of  threadID parameter
func (r *MessageMarkAsRead) GetThreadID() uint64 {
	return r.ThreadID
}

// HasLastReadMessageID returns true if lastReadMessageID was set
func (r *MessageMarkAsRead) HasLastReadMessageID() bool {
	return r.hasLastReadMessageID
}

// RawLastReadMessageID returns raw value of lastReadMessageID parameter
func (r *MessageMarkAsRead) RawLastReadMessageID() string {
	return r.rawLastReadMessageID
}

// GetLastReadMessageID returns casted value of  lastReadMessageID parameter
func (r *MessageMarkAsRead) GetLastReadMessageID() uint64 {
	return r.LastReadMessageID
}

// HasChannelID returns true if channelID was set
func (r *MessageMarkAsRead) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageMarkAsRead) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageMarkAsRead) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessageID returns true if messageID was set
func (r *MessageEdit) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessageEdit) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessageEdit) GetMessageID() uint64 {
	return r.MessageID
}

// HasChannelID returns true if channelID was set
func (r *MessageEdit) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageEdit) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageEdit) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessage returns true if message was set
func (r *MessageEdit) HasMessage() bool {
	return r.hasMessage
}

// RawMessage returns raw value of message parameter
func (r *MessageEdit) RawMessage() string {
	return r.rawMessage
}

// GetMessage returns casted value of  message parameter
func (r *MessageEdit) GetMessage() string {
	return r.Message
}

// HasMessageID returns true if messageID was set
func (r *MessageDelete) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessageDelete) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessageDelete) GetMessageID() uint64 {
	return r.MessageID
}

// HasChannelID returns true if channelID was set
func (r *MessageDelete) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageDelete) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageDelete) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessageID returns true if messageID was set
func (r *MessageReplyCreate) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessageReplyCreate) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessageReplyCreate) GetMessageID() uint64 {
	return r.MessageID
}

// HasChannelID returns true if channelID was set
func (r *MessageReplyCreate) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageReplyCreate) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageReplyCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessage returns true if message was set
func (r *MessageReplyCreate) HasMessage() bool {
	return r.hasMessage
}

// RawMessage returns raw value of message parameter
func (r *MessageReplyCreate) RawMessage() string {
	return r.rawMessage
}

// GetMessage returns casted value of  message parameter
func (r *MessageReplyCreate) GetMessage() string {
	return r.Message
}

// HasMessageID returns true if messageID was set
func (r *MessagePinCreate) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessagePinCreate) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessagePinCreate) GetMessageID() uint64 {
	return r.MessageID
}

// HasChannelID returns true if channelID was set
func (r *MessagePinCreate) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessagePinCreate) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessagePinCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessageID returns true if messageID was set
func (r *MessagePinRemove) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessagePinRemove) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessagePinRemove) GetMessageID() uint64 {
	return r.MessageID
}

// HasChannelID returns true if channelID was set
func (r *MessagePinRemove) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessagePinRemove) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessagePinRemove) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessageID returns true if messageID was set
func (r *MessageBookmarkCreate) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessageBookmarkCreate) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessageBookmarkCreate) GetMessageID() uint64 {
	return r.MessageID
}

// HasChannelID returns true if channelID was set
func (r *MessageBookmarkCreate) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageBookmarkCreate) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageBookmarkCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessageID returns true if messageID was set
func (r *MessageBookmarkRemove) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessageBookmarkRemove) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessageBookmarkRemove) GetMessageID() uint64 {
	return r.MessageID
}

// HasChannelID returns true if channelID was set
func (r *MessageBookmarkRemove) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageBookmarkRemove) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageBookmarkRemove) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessageID returns true if messageID was set
func (r *MessageReactionCreate) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessageReactionCreate) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessageReactionCreate) GetMessageID() uint64 {
	return r.MessageID
}

// HasReaction returns true if reaction was set
func (r *MessageReactionCreate) HasReaction() bool {
	return r.hasReaction
}

// RawReaction returns raw value of reaction parameter
func (r *MessageReactionCreate) RawReaction() string {
	return r.rawReaction
}

// GetReaction returns casted value of  reaction parameter
func (r *MessageReactionCreate) GetReaction() string {
	return r.Reaction
}

// HasChannelID returns true if channelID was set
func (r *MessageReactionCreate) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageReactionCreate) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageReactionCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessageID returns true if messageID was set
func (r *MessageReactionRemove) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *MessageReactionRemove) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *MessageReactionRemove) GetMessageID() uint64 {
	return r.MessageID
}

// HasReaction returns true if reaction was set
func (r *MessageReactionRemove) HasReaction() bool {
	return r.hasReaction
}

// RawReaction returns raw value of reaction parameter
func (r *MessageReactionRemove) RawReaction() string {
	return r.rawReaction
}

// GetReaction returns casted value of  reaction parameter
func (r *MessageReactionRemove) GetReaction() string {
	return r.Reaction
}

// HasChannelID returns true if channelID was set
func (r *MessageReactionRemove) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *MessageReactionRemove) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *MessageReactionRemove) GetChannelID() uint64 {
	return r.ChannelID
}
