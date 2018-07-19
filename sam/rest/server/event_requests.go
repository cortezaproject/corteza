package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `event.go`, `event.util.go` or `event_test.go` to
	implement your API calls, helper functions and tests. The file `event.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Event edit request parameters
type EventEditRequest struct {
	ID         uint64
	Channel_id uint64
	Contents   string
}

func (EventEditRequest) new() *EventEditRequest {
	return &EventEditRequest{}
}

func (e *EventEditRequest) Fill(r *http.Request) error {
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

	e.ID = parseUInt64(post["id"])

	e.Channel_id = parseUInt64(post["channel_id"])

	e.Contents = post["contents"]
	return nil
}

var _ RequestFiller = EventEditRequest{}.new()

// Event attach request parameters
type EventAttachRequest struct {
}

func (EventAttachRequest) new() *EventAttachRequest {
	return &EventAttachRequest{}
}

func (e *EventAttachRequest) Fill(r *http.Request) error {
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

var _ RequestFiller = EventAttachRequest{}.new()

// Event remove request parameters
type EventRemoveRequest struct {
	ID uint64
}

func (EventRemoveRequest) new() *EventRemoveRequest {
	return &EventRemoveRequest{}
}

func (e *EventRemoveRequest) Fill(r *http.Request) error {
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

	e.ID = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = EventRemoveRequest{}.new()

// Event read request parameters
type EventReadRequest struct {
	Channel_id uint64
}

func (EventReadRequest) new() *EventReadRequest {
	return &EventReadRequest{}
}

func (e *EventReadRequest) Fill(r *http.Request) error {
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

	e.Channel_id = parseUInt64(post["channel_id"])
	return nil
}

var _ RequestFiller = EventReadRequest{}.new()

// Event search request parameters
type EventSearchRequest struct {
	Query        string
	Message_type string
}

func (EventSearchRequest) new() *EventSearchRequest {
	return &EventSearchRequest{}
}

func (e *EventSearchRequest) Fill(r *http.Request) error {
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

	e.Query = get["query"]

	e.Message_type = get["message_type"]
	return nil
}

var _ RequestFiller = EventSearchRequest{}.new()

// Event pin request parameters
type EventPinRequest struct {
	ID uint64
}

func (EventPinRequest) new() *EventPinRequest {
	return &EventPinRequest{}
}

func (e *EventPinRequest) Fill(r *http.Request) error {
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

	e.ID = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = EventPinRequest{}.new()

// Event flag request parameters
type EventFlagRequest struct {
	ID uint64
}

func (EventFlagRequest) new() *EventFlagRequest {
	return &EventFlagRequest{}
}

func (e *EventFlagRequest) Fill(r *http.Request) error {
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

	e.ID = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = EventFlagRequest{}.new()
