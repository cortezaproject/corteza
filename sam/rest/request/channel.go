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
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}
var _ = multipart.FileHeader{}

// Channel list request parameters
type ChannelList struct {
	Query string
}

func NewChannelList() *ChannelList {
	return &ChannelList{}
}

func (c *ChannelList) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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

	if val, ok := get["query"]; ok {

		c.Query = val
	}

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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

	if val, ok := post["name"]; ok {

		c.Name = val
	}
	if val, ok := post["topic"]; ok {

		c.Topic = val
	}

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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
	if val, ok := post["name"]; ok {

		c.Name = val
	}
	if val, ok := post["topic"]; ok {

		c.Topic = val
	}
	if val, ok := post["archive"]; ok {

		c.Archive = parseBool(val)
	}
	if val, ok := post["organisationID"]; ok {

		c.OrganisationID = parseUInt64(val)
	}

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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

	return err
}

var _ RequestFiller = NewChannelMembers()

// Channel join request parameters
type ChannelJoin struct {
	ChannelID uint64
	UserID    uint64
}

func NewChannelJoin() *ChannelJoin {
	return &ChannelJoin{}
}

func (c *ChannelJoin) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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

	return err
}

var _ RequestFiller = NewChannelInvite()

// Channel attach request parameters
type ChannelAttach struct {
	ChannelID uint64
	Upload    *multipart.FileHeader
}

func NewChannelAttach() *ChannelAttach {
	return &ChannelAttach{}
}

func (c *ChannelAttach) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(c)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

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
	if _, c.Upload, err = r.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	return err
}

var _ RequestFiller = NewChannelAttach()
