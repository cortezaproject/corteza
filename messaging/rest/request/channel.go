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

// Channel list request parameters
type ChannelList struct {
	Query string
}

func NewChannelList() *ChannelList {
	return &ChannelList{}
}

func (r ChannelList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query

	return out
}

func (r *ChannelList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["query"]; ok {
		r.Query = val
	}

	return err
}

var _ RequestFiller = NewChannelList()

// Channel create request parameters
type ChannelCreate struct {
	Name    string
	Topic   string
	Type    string
	Members []string
}

func NewChannelCreate() *ChannelCreate {
	return &ChannelCreate{}
}

func (r ChannelCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["topic"] = r.Topic
	out["type"] = r.Type
	out["members"] = r.Members

	return out
}

func (r *ChannelCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["topic"]; ok {
		r.Topic = val
	}
	if val, ok := post["type"]; ok {
		r.Type = val
	}

	return err
}

var _ RequestFiller = NewChannelCreate()

// Channel update request parameters
type ChannelUpdate struct {
	ChannelID      uint64 `json:",string"`
	Name           string
	Topic          string
	Type           string
	OrganisationID uint64 `json:",string"`
}

func NewChannelUpdate() *ChannelUpdate {
	return &ChannelUpdate{}
}

func (r ChannelUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["name"] = r.Name
	out["topic"] = r.Topic
	out["type"] = r.Type
	out["organisationID"] = r.OrganisationID

	return out
}

func (r *ChannelUpdate) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["topic"]; ok {
		r.Topic = val
	}
	if val, ok := post["type"]; ok {
		r.Type = val
	}
	if val, ok := post["organisationID"]; ok {
		r.OrganisationID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewChannelUpdate()

// Channel state request parameters
type ChannelState struct {
	ChannelID uint64 `json:",string"`
	State     string
}

func NewChannelState() *ChannelState {
	return &ChannelState{}
}

func (r ChannelState) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["state"] = r.State

	return out
}

func (r *ChannelState) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["state"]; ok {
		r.State = val
	}

	return err
}

var _ RequestFiller = NewChannelState()

// Channel setFlag request parameters
type ChannelSetFlag struct {
	ChannelID uint64 `json:",string"`
	Flag      string
}

func NewChannelSetFlag() *ChannelSetFlag {
	return &ChannelSetFlag{}
}

func (r ChannelSetFlag) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["flag"] = r.Flag

	return out
}

func (r *ChannelSetFlag) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["flag"]; ok {
		r.Flag = val
	}

	return err
}

var _ RequestFiller = NewChannelSetFlag()

// Channel removeFlag request parameters
type ChannelRemoveFlag struct {
	ChannelID uint64 `json:",string"`
}

func NewChannelRemoveFlag() *ChannelRemoveFlag {
	return &ChannelRemoveFlag{}
}

func (r ChannelRemoveFlag) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID

	return out
}

func (r *ChannelRemoveFlag) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewChannelRemoveFlag()

// Channel read request parameters
type ChannelRead struct {
	ChannelID uint64 `json:",string"`
}

func NewChannelRead() *ChannelRead {
	return &ChannelRead{}
}

func (r ChannelRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID

	return out
}

func (r *ChannelRead) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewChannelRead()

// Channel members request parameters
type ChannelMembers struct {
	ChannelID uint64 `json:",string"`
}

func NewChannelMembers() *ChannelMembers {
	return &ChannelMembers{}
}

func (r ChannelMembers) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID

	return out
}

func (r *ChannelMembers) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewChannelMembers()

// Channel join request parameters
type ChannelJoin struct {
	ChannelID uint64 `json:",string"`
	UserID    uint64 `json:",string"`
}

func NewChannelJoin() *ChannelJoin {
	return &ChannelJoin{}
}

func (r ChannelJoin) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["userID"] = r.UserID

	return out
}

func (r *ChannelJoin) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewChannelJoin()

// Channel part request parameters
type ChannelPart struct {
	ChannelID uint64 `json:",string"`
	UserID    uint64 `json:",string"`
}

func NewChannelPart() *ChannelPart {
	return &ChannelPart{}
}

func (r ChannelPart) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["userID"] = r.UserID

	return out
}

func (r *ChannelPart) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewChannelPart()

// Channel invite request parameters
type ChannelInvite struct {
	ChannelID uint64 `json:",string"`
	UserID    []string
}

func NewChannelInvite() *ChannelInvite {
	return &ChannelInvite{}
}

func (r ChannelInvite) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["userID"] = r.UserID

	return out
}

func (r *ChannelInvite) Fill(req *http.Request) (err error) {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewChannelInvite()

// Channel attach request parameters
type ChannelAttach struct {
	ChannelID uint64 `json:",string"`
	ReplyTo   uint64 `json:",string"`
	Upload    *multipart.FileHeader
}

func NewChannelAttach() *ChannelAttach {
	return &ChannelAttach{}
}

func (r ChannelAttach) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["replyTo"] = r.ReplyTo
	out["upload.size"] = r.Upload.Size
	out["upload.filename"] = r.Upload.Filename

	return out
}

func (r *ChannelAttach) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseMultipartForm(32 << 20); err != nil {
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

	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["replyTo"]; ok {
		r.ReplyTo = parseUInt64(val)
	}
	if _, r.Upload, err = req.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	return err
}

var _ RequestFiller = NewChannelAttach()
