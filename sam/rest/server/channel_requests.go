package server

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
type ChannelListRequest struct {
	Query string
}

func (ChannelListRequest) new() *ChannelListRequest {
	return &ChannelListRequest{}
}

func (c *ChannelListRequest) Fill(r *http.Request) error {
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

	c.Query = get["query"]
	return nil
}

var _ RequestFiller = ChannelListRequest{}.new()

// Channel create request parameters
type ChannelCreateRequest struct {
	Name  string
	Topic string
}

func (ChannelCreateRequest) new() *ChannelCreateRequest {
	return &ChannelCreateRequest{}
}

func (c *ChannelCreateRequest) Fill(r *http.Request) error {
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

	c.Name = post["name"]

	c.Topic = post["topic"]
	return nil
}

var _ RequestFiller = ChannelCreateRequest{}.new()

// Channel edit request parameters
type ChannelEditRequest struct {
	ChannelId      uint64
	Name           string
	Topic          string
	Archive        bool
	OrganisationId uint64
}

func (ChannelEditRequest) new() *ChannelEditRequest {
	return &ChannelEditRequest{}
}

func (c *ChannelEditRequest) Fill(r *http.Request) error {
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

	c.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))

	c.Name = post["name"]

	c.Topic = post["topic"]

	c.Archive = parseBool(post["archive"])

	c.OrganisationId = parseUInt64(post["organisationId"])
	return nil
}

var _ RequestFiller = ChannelEditRequest{}.new()

// Channel read request parameters
type ChannelReadRequest struct {
	ChannelId uint64
}

func (ChannelReadRequest) new() *ChannelReadRequest {
	return &ChannelReadRequest{}
}

func (c *ChannelReadRequest) Fill(r *http.Request) error {
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

	c.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = ChannelReadRequest{}.new()

// Channel delete request parameters
type ChannelDeleteRequest struct {
	ChannelId uint64
}

func (ChannelDeleteRequest) new() *ChannelDeleteRequest {
	return &ChannelDeleteRequest{}
}

func (c *ChannelDeleteRequest) Fill(r *http.Request) error {
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

	c.ChannelId = parseUInt64(chi.URLParam(r, "channelId"))
	return nil
}

var _ RequestFiller = ChannelDeleteRequest{}.new()
