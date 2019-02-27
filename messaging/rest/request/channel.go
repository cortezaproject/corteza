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
	"io"
	"mime/multipart"
	"net/http"
	"strings"

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

func (cReq *ChannelList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

		cReq.Query = val
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

func (cReq *ChannelCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

		cReq.Name = val
	}
	if val, ok := post["topic"]; ok {

		cReq.Topic = val
	}
	if val, ok := post["type"]; ok {

		cReq.Type = val
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

func (cReq *ChannelUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["name"]; ok {

		cReq.Name = val
	}
	if val, ok := post["topic"]; ok {

		cReq.Topic = val
	}
	if val, ok := post["type"]; ok {

		cReq.Type = val
	}
	if val, ok := post["organisationID"]; ok {

		cReq.OrganisationID = parseUInt64(val)
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

func (cReq *ChannelState) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["state"]; ok {

		cReq.State = val
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

func (cReq *ChannelSetFlag) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["flag"]; ok {

		cReq.Flag = val
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

func (cReq *ChannelRemoveFlag) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

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

func (cReq *ChannelRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

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

func (cReq *ChannelMembers) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

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

func (cReq *ChannelJoin) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	cReq.UserID = parseUInt64(chi.URLParam(r, "userID"))

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

func (cReq *ChannelPart) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	cReq.UserID = parseUInt64(chi.URLParam(r, "userID"))

	return err
}

var _ RequestFiller = NewChannelPart()

// Channel invite request parameters
type ChannelInvite struct {
	ChannelID uint64   `json:",string"`
	UserID    []uint64 `json:",string"`
}

func NewChannelInvite() *ChannelInvite {
	return &ChannelInvite{}
}

func (cReq *ChannelInvite) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	cReq.UserID = parseUInt64A(r.Form["userID"])

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

func (cReq *ChannelAttach) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseMultipartForm(32 << 20); err != nil {
		return err
	}

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

	cReq.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["replyTo"]; ok {

		cReq.ReplyTo = parseUInt64(val)
	}
	if _, cReq.Upload, err = r.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	return err
}

var _ RequestFiller = NewChannelAttach()
