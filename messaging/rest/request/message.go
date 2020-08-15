package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
)

type (
	// Internal API interface
	MessageCreate struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// Message POST parameter
		//
		// Message contents (markdown)
		Message string
	}

	MessageExecuteCommand struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// Command PATH parameter
		//
		// Command to be executed
		Command string

		// Input POST parameter
		//
		// Arbitrary command input
		Input string

		// Params POST parameter
		//
		// Command parameters
		Params []string
	}

	MessageMarkAsRead struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// ThreadID GET parameter
		//
		// ID of thread (messageID)
		ThreadID uint64 `json:",string"`

		// LastReadMessageID GET parameter
		//
		// ID of the last read message
		LastReadMessageID uint64 `json:",string"`
	}

	MessageEdit struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`

		// Message POST parameter
		//
		// Message contents (markdown)
		Message string
	}

	MessageDelete struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`
	}

	MessageReplyCreate struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`

		// Message POST parameter
		//
		// Message contents (markdown)
		Message string
	}

	MessagePinCreate struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`
	}

	MessagePinRemove struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`
	}

	MessageBookmarkCreate struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`
	}

	MessageBookmarkRemove struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`
	}

	MessageReactionCreate struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`

		// Reaction PATH parameter
		//
		// Reaction
		Reaction string
	}

	MessageReactionRemove struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// MessageID PATH parameter
		//
		// Message ID
		MessageID uint64 `json:",string"`

		// Reaction PATH parameter
		//
		// Reaction
		Reaction string
	}
)

// NewMessageCreate request
func NewMessageCreate() *MessageCreate {
	return &MessageCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"message":   r.Message,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageCreate) GetMessage() string {
	return r.Message
}

// Fill processes request and fills internal variables
func (r *MessageCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["message"]; ok && len(val) > 0 {
			r.Message, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageExecuteCommand request
func NewMessageExecuteCommand() *MessageExecuteCommand {
	return &MessageExecuteCommand{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageExecuteCommand) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"command":   r.Command,
		"input":     r.Input,
		"params":    r.Params,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageExecuteCommand) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageExecuteCommand) GetCommand() string {
	return r.Command
}

// Auditable returns all auditable/loggable parameters
func (r MessageExecuteCommand) GetInput() string {
	return r.Input
}

// Auditable returns all auditable/loggable parameters
func (r MessageExecuteCommand) GetParams() []string {
	return r.Params
}

// Fill processes request and fills internal variables
func (r *MessageExecuteCommand) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["input"]; ok && len(val) > 0 {
			r.Input, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["params[]"]; ok && len(val) > 0  {
		//    r.Params, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "command")
		r.Command, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageMarkAsRead request
func NewMessageMarkAsRead() *MessageMarkAsRead {
	return &MessageMarkAsRead{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageMarkAsRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID":         r.ChannelID,
		"threadID":          r.ThreadID,
		"lastReadMessageID": r.LastReadMessageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageMarkAsRead) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageMarkAsRead) GetThreadID() uint64 {
	return r.ThreadID
}

// Auditable returns all auditable/loggable parameters
func (r MessageMarkAsRead) GetLastReadMessageID() uint64 {
	return r.LastReadMessageID
}

// Fill processes request and fills internal variables
func (r *MessageMarkAsRead) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["threadID"]; ok && len(val) > 0 {
			r.ThreadID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["lastReadMessageID"]; ok && len(val) > 0 {
			r.LastReadMessageID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageEdit request
func NewMessageEdit() *MessageEdit {
	return &MessageEdit{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageEdit) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
		"message":   r.Message,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageEdit) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageEdit) GetMessageID() uint64 {
	return r.MessageID
}

// Auditable returns all auditable/loggable parameters
func (r MessageEdit) GetMessage() string {
	return r.Message
}

// Fill processes request and fills internal variables
func (r *MessageEdit) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["message"]; ok && len(val) > 0 {
			r.Message, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageDelete request
func NewMessageDelete() *MessageDelete {
	return &MessageDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageDelete) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageDelete) GetMessageID() uint64 {
	return r.MessageID
}

// Fill processes request and fills internal variables
func (r *MessageDelete) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageReplyCreate request
func NewMessageReplyCreate() *MessageReplyCreate {
	return &MessageReplyCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReplyCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
		"message":   r.Message,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReplyCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageReplyCreate) GetMessageID() uint64 {
	return r.MessageID
}

// Auditable returns all auditable/loggable parameters
func (r MessageReplyCreate) GetMessage() string {
	return r.Message
}

// Fill processes request and fills internal variables
func (r *MessageReplyCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["message"]; ok && len(val) > 0 {
			r.Message, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessagePinCreate request
func NewMessagePinCreate() *MessagePinCreate {
	return &MessagePinCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessagePinCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessagePinCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessagePinCreate) GetMessageID() uint64 {
	return r.MessageID
}

// Fill processes request and fills internal variables
func (r *MessagePinCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessagePinRemove request
func NewMessagePinRemove() *MessagePinRemove {
	return &MessagePinRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r MessagePinRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessagePinRemove) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessagePinRemove) GetMessageID() uint64 {
	return r.MessageID
}

// Fill processes request and fills internal variables
func (r *MessagePinRemove) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageBookmarkCreate request
func NewMessageBookmarkCreate() *MessageBookmarkCreate {
	return &MessageBookmarkCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageBookmarkCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageBookmarkCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageBookmarkCreate) GetMessageID() uint64 {
	return r.MessageID
}

// Fill processes request and fills internal variables
func (r *MessageBookmarkCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageBookmarkRemove request
func NewMessageBookmarkRemove() *MessageBookmarkRemove {
	return &MessageBookmarkRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageBookmarkRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageBookmarkRemove) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageBookmarkRemove) GetMessageID() uint64 {
	return r.MessageID
}

// Fill processes request and fills internal variables
func (r *MessageBookmarkRemove) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageReactionCreate request
func NewMessageReactionCreate() *MessageReactionCreate {
	return &MessageReactionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
		"reaction":  r.Reaction,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionCreate) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionCreate) GetMessageID() uint64 {
	return r.MessageID
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionCreate) GetReaction() string {
	return r.Reaction
}

// Fill processes request and fills internal variables
func (r *MessageReactionCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "reaction")
		r.Reaction, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewMessageReactionRemove request
func NewMessageReactionRemove() *MessageReactionRemove {
	return &MessageReactionRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
		"reaction":  r.Reaction,
	}
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionRemove) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionRemove) GetMessageID() uint64 {
	return r.MessageID
}

// Auditable returns all auditable/loggable parameters
func (r MessageReactionRemove) GetReaction() string {
	return r.Reaction
}

// Fill processes request and fills internal variables
func (r *MessageReactionRemove) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "messageID")
		r.MessageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "reaction")
		r.Reaction, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
