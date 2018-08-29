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
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Message create request parameters
type MessageCreate struct {
	Message   string
	ChannelID uint64
}

func NewMessageCreate() *MessageCreate {
	return &MessageCreate{}
}

func (m *MessageCreate) Fill(r *http.Request) error {
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
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageCreate()

// Message history request parameters
type MessageHistory struct {
	LastMessageID uint64
	ChannelID     uint64
}

func NewMessageHistory() *MessageHistory {
	return &MessageHistory{}
}

func (m *MessageHistory) Fill(r *http.Request) error {
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

	m.LastMessageID = parseUInt64(get["lastMessageID"])
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageHistory()

// Message edit request parameters
type MessageEdit struct {
	MessageID uint64
	ChannelID uint64
	Message   string
}

func NewMessageEdit() *MessageEdit {
	return &MessageEdit{}
}

func (m *MessageEdit) Fill(r *http.Request) error {
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	m.Message = post["message"]

	return nil
}

var _ RequestFiller = NewMessageEdit()

// Message delete request parameters
type MessageDelete struct {
	MessageID uint64
	ChannelID uint64
}

func NewMessageDelete() *MessageDelete {
	return &MessageDelete{}
}

func (m *MessageDelete) Fill(r *http.Request) error {
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageDelete()

// Message attach request parameters
type MessageAttach struct {
	ChannelID uint64
}

func NewMessageAttach() *MessageAttach {
	return &MessageAttach{}
}

func (m *MessageAttach) Fill(r *http.Request) error {
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

	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageAttach()

// Message search request parameters
type MessageSearch struct {
	Query        string
	Message_type string
	ChannelID    uint64
}

func NewMessageSearch() *MessageSearch {
	return &MessageSearch{}
}

func (m *MessageSearch) Fill(r *http.Request) error {
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
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageSearch()

// Message pin request parameters
type MessagePin struct {
	MessageID uint64
	ChannelID uint64
}

func NewMessagePin() *MessagePin {
	return &MessagePin{}
}

func (m *MessagePin) Fill(r *http.Request) error {
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessagePin()

// Message unpin request parameters
type MessageUnpin struct {
	MessageID uint64
	ChannelID uint64
}

func NewMessageUnpin() *MessageUnpin {
	return &MessageUnpin{}
}

func (m *MessageUnpin) Fill(r *http.Request) error {
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageUnpin()

// Message flag request parameters
type MessageFlag struct {
	MessageID uint64
	ChannelID uint64
}

func NewMessageFlag() *MessageFlag {
	return &MessageFlag{}
}

func (m *MessageFlag) Fill(r *http.Request) error {
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageFlag()

// Message unflag request parameters
type MessageUnflag struct {
	MessageID uint64
	ChannelID uint64
}

func NewMessageUnflag() *MessageUnflag {
	return &MessageUnflag{}
}

func (m *MessageUnflag) Fill(r *http.Request) error {
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageUnflag()

// Message react request parameters
type MessageReact struct {
	MessageID uint64
	Reaction  string
	ChannelID uint64
}

func NewMessageReact() *MessageReact {
	return &MessageReact{}
}

func (m *MessageReact) Fill(r *http.Request) error {
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.Reaction = chi.URLParam(r, "reaction")
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageReact()

// Message unreact request parameters
type MessageUnreact struct {
	MessageID uint64
	Reaction  string
	ChannelID uint64
}

func NewMessageUnreact() *MessageUnreact {
	return &MessageUnreact{}
}

func (m *MessageUnreact) Fill(r *http.Request) error {
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.Reaction = chi.URLParam(r, "reaction")
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewMessageUnreact()
