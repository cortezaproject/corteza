package server

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
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Message create request parameters
type MessageCreateRequest struct {
	Message   string
	ChannelId uint64
}

func (MessageCreateRequest) new() *MessageCreateRequest {
	return &MessageCreateRequest{}
}

func (m *MessageCreateRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.Message = post["message"]

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageCreateRequest{}.new()

// Message history request parameters
type MessageHistoryRequest struct {
	LastMessageId uint64
	ChannelId     uint64
}

func (MessageHistoryRequest) new() *MessageHistoryRequest {
	return &MessageHistoryRequest{}
}

func (m *MessageHistoryRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.LastMessageId = parseUInt64(get["lastMessageId"])

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageHistoryRequest{}.new()

// Message edit request parameters
type MessageEditRequest struct {
	MessageId uint64
	ChannelId uint64
	Message   string
}

func (MessageEditRequest) new() *MessageEditRequest {
	return &MessageEditRequest{}
}

func (m *MessageEditRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.MessageId = parseUInt64(chi.URLParam(r, "messageId"))

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))

	m.Message = post["message"]
	return nil
}

var _ RequestFiller = MessageEditRequest{}.new()

// Message delete request parameters
type MessageDeleteRequest struct {
	MessageId uint64
	ChannelId uint64
}

func (MessageDeleteRequest) new() *MessageDeleteRequest {
	return &MessageDeleteRequest{}
}

func (m *MessageDeleteRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.MessageId = parseUInt64(chi.URLParam(r, "messageId"))

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageDeleteRequest{}.new()

// Message attach request parameters
type MessageAttachRequest struct {
	ChannelId uint64
}

func (MessageAttachRequest) new() *MessageAttachRequest {
	return &MessageAttachRequest{}
}

func (m *MessageAttachRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageAttachRequest{}.new()

// Message search request parameters
type MessageSearchRequest struct {
	Query        string
	Message_type string
	ChannelId    uint64
}

func (MessageSearchRequest) new() *MessageSearchRequest {
	return &MessageSearchRequest{}
}

func (m *MessageSearchRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.Query = get["query"]

	m.Message_type = get["message_type"]

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageSearchRequest{}.new()

// Message pin request parameters
type MessagePinRequest struct {
	MessageId uint64
	ChannelId uint64
}

func (MessagePinRequest) new() *MessagePinRequest {
	return &MessagePinRequest{}
}

func (m *MessagePinRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.MessageId = parseUInt64(chi.URLParam(r, "messageId"))

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessagePinRequest{}.new()

// Message unpin request parameters
type MessageUnpinRequest struct {
	MessageId uint64
	ChannelId uint64
}

func (MessageUnpinRequest) new() *MessageUnpinRequest {
	return &MessageUnpinRequest{}
}

func (m *MessageUnpinRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.MessageId = parseUInt64(chi.URLParam(r, "messageId"))

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageUnpinRequest{}.new()

// Message flag request parameters
type MessageFlagRequest struct {
	MessageId uint64
	ChannelId uint64
}

func (MessageFlagRequest) new() *MessageFlagRequest {
	return &MessageFlagRequest{}
}

func (m *MessageFlagRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.MessageId = parseUInt64(chi.URLParam(r, "messageId"))

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageFlagRequest{}.new()

// Message unflag request parameters
type MessageUnflagRequest struct {
	MessageId uint64
	ChannelId uint64
}

func (MessageUnflagRequest) new() *MessageUnflagRequest {
	return &MessageUnflagRequest{}
}

func (m *MessageUnflagRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.MessageId = parseUInt64(chi.URLParam(r, "messageId"))

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageUnflagRequest{}.new()

// Message react request parameters
type MessageReactRequest struct {
	MessageId uint64
	Reaction  string
	ChannelId uint64
}

func (MessageReactRequest) new() *MessageReactRequest {
	return &MessageReactRequest{}
}

func (m *MessageReactRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.MessageId = parseUInt64(chi.URLParam(r, "messageId"))

	m.Reaction = chi.URLParam(r, "reaction")

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageReactRequest{}.new()

// Message unreact request parameters
type MessageUnreactRequest struct {
	MessageId uint64
	Reaction  string
	ChannelId uint64
}

func (MessageUnreactRequest) new() *MessageUnreactRequest {
	return &MessageUnreactRequest{}
}

func (m *MessageUnreactRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	m.MessageId = parseUInt64(chi.URLParam(r, "messageId"))

	m.Reaction = chi.URLParam(r, "reaction")

	m.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = MessageUnreactRequest{}.new()
