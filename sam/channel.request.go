package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `channel.go`, `channel.util.go` or `channel_test.go` to
	implement your API calls, helper functions and tests. The file `channel.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Channel list request parameters
type channelListRequest struct {
	query string
}

func (channelListRequest) new() *channelListRequest {
	return &channelListRequest{}
}

func (c *channelListRequest) Fill(r *http.Request) error {
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

	c.query = get["query"]
	return nil
}

var _ RequestFiller = channelListRequest{}.new()

// Channel create request parameters
type channelCreateRequest struct {
	name  string
	topic string
}

func (channelCreateRequest) new() *channelCreateRequest {
	return &channelCreateRequest{}
}

func (c *channelCreateRequest) Fill(r *http.Request) error {
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

	c.name = post["name"]

	c.topic = post["topic"]
	return nil
}

var _ RequestFiller = channelCreateRequest{}.new()

// Channel edit request parameters
type channelEditRequest struct {
	id             uint64
	name           string
	topic          string
	archive        bool
	organisationId uint64
}

func (channelEditRequest) new() *channelEditRequest {
	return &channelEditRequest{}
}

func (c *channelEditRequest) Fill(r *http.Request) error {
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

	c.id = parseUInt64(post["id"])

	c.name = post["name"]

	c.topic = post["topic"]

	c.archive = parseBool(post["archive"])

	c.organisationId = parseUInt64(post["organisationId"])
	return nil
}

var _ RequestFiller = channelEditRequest{}.new()

// Channel read request parameters
type channelReadRequest struct {
	id uint64
}

func (channelReadRequest) new() *channelReadRequest {
	return &channelReadRequest{}
}

func (c *channelReadRequest) Fill(r *http.Request) error {
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

	c.id = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = channelReadRequest{}.new()

// Channel delete request parameters
type channelDeleteRequest struct {
	id uint64
}

func (channelDeleteRequest) new() *channelDeleteRequest {
	return &channelDeleteRequest{}
}

func (c *channelDeleteRequest) Fill(r *http.Request) error {
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

	c.id = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = channelDeleteRequest{}.new()
