package rest

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

// Message edit request parameters
type MessageEditRequest struct {
	ID         uint64
	Channel_id uint64
	Contents   string
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

	m.ID = parseUInt64(post["id"])

	m.Channel_id = parseUInt64(post["channel_id"])

	m.Contents = post["contents"]
	return nil
}

var _ RequestFiller = MessageEditRequest{}.new()

// Message attach request parameters
type MessageAttachRequest struct {
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
	return nil
}

var _ RequestFiller = MessageAttachRequest{}.new()

// Message remove request parameters
type MessageRemoveRequest struct {
	ID uint64
}

func (MessageRemoveRequest) new() *MessageRemoveRequest {
	return &MessageRemoveRequest{}
}

func (m *MessageRemoveRequest) Fill(r *http.Request) error {
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

	m.ID = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = MessageRemoveRequest{}.new()

// Message read request parameters
type MessageReadRequest struct {
	Channel_id uint64
}

func (MessageReadRequest) new() *MessageReadRequest {
	return &MessageReadRequest{}
}

func (m *MessageReadRequest) Fill(r *http.Request) error {
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

	m.Channel_id = parseUInt64(post["channel_id"])
	return nil
}

var _ RequestFiller = MessageReadRequest{}.new()

// Message search request parameters
type MessageSearchRequest struct {
	Query        string
	Message_type string
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
	return nil
}

var _ RequestFiller = MessageSearchRequest{}.new()

// Message pin request parameters
type MessagePinRequest struct {
	ID uint64
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

	m.ID = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = MessagePinRequest{}.new()

// Message flag request parameters
type MessageFlagRequest struct {
	ID uint64
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

	m.ID = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = MessageFlagRequest{}.new()
