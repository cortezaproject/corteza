package request

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
type ChannelList struct {
	Query string
}

func NewChannelList() *ChannelList {
	return &ChannelList{}
}

func (c *ChannelList) Fill(r *http.Request) error {
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

var _ RequestFiller = NewChannelList()

// Channel create request parameters
type ChannelCreate struct {
	Name  string
	Topic string
}

func NewChannelCreate() *ChannelCreate {
	return &ChannelCreate{}
}

func (c *ChannelCreate) Fill(r *http.Request) error {
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

var _ RequestFiller = NewChannelCreate()

// Channel edit request parameters
type ChannelEdit struct {
	ChannelID      uint64
	Name           string
	Topic          string
	Archive        bool
	OrganisationID uint64
}

func NewChannelEdit() *ChannelEdit {
	return &ChannelEdit{}
}

func (c *ChannelEdit) Fill(r *http.Request) error {
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

	c.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	c.Name = post["name"]
	c.Topic = post["topic"]
	c.Archive = parseBool(post["archive"])
	c.OrganisationID = parseUInt64(post["organisationID"])

	return nil
}

var _ RequestFiller = NewChannelEdit()

// Channel read request parameters
type ChannelRead struct {
	ChannelID uint64
}

func NewChannelRead() *ChannelRead {
	return &ChannelRead{}
}

func (c *ChannelRead) Fill(r *http.Request) error {
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

	c.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewChannelRead()

// Channel delete request parameters
type ChannelDelete struct {
	ChannelID uint64
}

func NewChannelDelete() *ChannelDelete {
	return &ChannelDelete{}
}

func (c *ChannelDelete) Fill(r *http.Request) error {
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

	c.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewChannelDelete()

// Channel members request parameters
type ChannelMembers struct {
	ChannelID uint64
}

func NewChannelMembers() *ChannelMembers {
	return &ChannelMembers{}
}

func (c *ChannelMembers) Fill(r *http.Request) error {
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

	c.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewChannelMembers()

// Channel join request parameters
type ChannelJoin struct {
	ChannelID uint64
}

func NewChannelJoin() *ChannelJoin {
	return &ChannelJoin{}
}

func (c *ChannelJoin) Fill(r *http.Request) error {
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

	c.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewChannelJoin()

// Channel part request parameters
type ChannelPart struct {
	ChannelID uint64
	UserID    uint64
}

func NewChannelPart() *ChannelPart {
	return &ChannelPart{}
}

func (c *ChannelPart) Fill(r *http.Request) error {
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

	c.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	c.UserID = parseUInt64(chi.URLParam(r, "userID"))

	return nil
}

var _ RequestFiller = NewChannelPart()

// Channel invite request parameters
type ChannelInvite struct {
	ChannelID uint64
	UserID    []uint64
}

func NewChannelInvite() *ChannelInvite {
	return &ChannelInvite{}
}

func (c *ChannelInvite) Fill(r *http.Request) error {
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

	c.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return nil
}

var _ RequestFiller = NewChannelInvite()
