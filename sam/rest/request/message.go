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
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}

// Message create request parameters
type MessageCreate struct {
	Message   string
	ChannelID uint64
}

func NewMessageCreate() *MessageCreate {
	return &MessageCreate{}
}

func (m *MessageCreate) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	if val, ok := post["message"]; ok {

		m.Message = val
	}
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	if val, ok := get["lastMessageID"]; ok {

		m.LastMessageID = parseUInt64(val)
	}
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))
	if val, ok := post["message"]; ok {

		m.Message = val
	}

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

		m.Query = val
	}
	if val, ok := get["message_type"]; ok {

		m.Message_type = val
	}
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.Reaction = chi.URLParam(r, "reaction")
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
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
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
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

	m.MessageID = parseUInt64(chi.URLParam(r, "messageID"))
	m.Reaction = chi.URLParam(r, "reaction")
	m.ChannelID = parseUInt64(chi.URLParam(r, "channelID"))

	return err
}

var _ RequestFiller = NewMessageUnreact()
