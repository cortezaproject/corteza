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

	"github.com/cortezaproject/corteza-server/messaging/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// ChannelList request parameters
type ChannelList struct {
	hasQuery bool
	rawQuery string
	Query    string
}

// NewChannelList request
func NewChannelList() *ChannelList {
	return &ChannelList{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query

	return out
}

// Fill processes request and fills internal variables
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
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}

	return err
}

var _ RequestFiller = NewChannelList()

// ChannelCreate request parameters
type ChannelCreate struct {
	hasName bool
	rawName string
	Name    string

	hasTopic bool
	rawTopic string
	Topic    string

	hasType bool
	rawType string
	Type    string

	hasMembershipPolicy bool
	rawMembershipPolicy string
	MembershipPolicy    types.ChannelMembershipPolicy

	hasMembers bool
	rawMembers []string
	Members    []string
}

// NewChannelCreate request
func NewChannelCreate() *ChannelCreate {
	return &ChannelCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["topic"] = r.Topic
	out["type"] = r.Type
	out["membershipPolicy"] = r.MembershipPolicy
	out["members"] = r.Members

	return out
}

// Fill processes request and fills internal variables
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
		r.hasName = true
		r.rawName = val
		r.Name = val
	}
	if val, ok := post["topic"]; ok {
		r.hasTopic = true
		r.rawTopic = val
		r.Topic = val
	}
	if val, ok := post["type"]; ok {
		r.hasType = true
		r.rawType = val
		r.Type = val
	}
	if val, ok := post["membershipPolicy"]; ok {
		r.hasMembershipPolicy = true
		r.rawMembershipPolicy = val
		r.MembershipPolicy = types.ChannelMembershipPolicy(val)
	}

	if val, ok := req.Form["members"]; ok {
		r.hasMembers = true
		r.rawMembers = val
		r.Members = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewChannelCreate()

// ChannelUpdate request parameters
type ChannelUpdate struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasName bool
	rawName string
	Name    string

	hasTopic bool
	rawTopic string
	Topic    string

	hasMembershipPolicy bool
	rawMembershipPolicy string
	MembershipPolicy    types.ChannelMembershipPolicy

	hasType bool
	rawType string
	Type    string

	hasOrganisationID bool
	rawOrganisationID string
	OrganisationID    uint64 `json:",string"`
}

// NewChannelUpdate request
func NewChannelUpdate() *ChannelUpdate {
	return &ChannelUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["name"] = r.Name
	out["topic"] = r.Topic
	out["membershipPolicy"] = r.MembershipPolicy
	out["type"] = r.Type
	out["organisationID"] = r.OrganisationID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}
	if val, ok := post["topic"]; ok {
		r.hasTopic = true
		r.rawTopic = val
		r.Topic = val
	}
	if val, ok := post["membershipPolicy"]; ok {
		r.hasMembershipPolicy = true
		r.rawMembershipPolicy = val
		r.MembershipPolicy = types.ChannelMembershipPolicy(val)
	}
	if val, ok := post["type"]; ok {
		r.hasType = true
		r.rawType = val
		r.Type = val
	}
	if val, ok := post["organisationID"]; ok {
		r.hasOrganisationID = true
		r.rawOrganisationID = val
		r.OrganisationID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewChannelUpdate()

// ChannelState request parameters
type ChannelState struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasState bool
	rawState string
	State    string
}

// NewChannelState request
func NewChannelState() *ChannelState {
	return &ChannelState{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelState) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["state"] = r.State

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["state"]; ok {
		r.hasState = true
		r.rawState = val
		r.State = val
	}

	return err
}

var _ RequestFiller = NewChannelState()

// ChannelSetFlag request parameters
type ChannelSetFlag struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasFlag bool
	rawFlag string
	Flag    string
}

// NewChannelSetFlag request
func NewChannelSetFlag() *ChannelSetFlag {
	return &ChannelSetFlag{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelSetFlag) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["flag"] = r.Flag

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["flag"]; ok {
		r.hasFlag = true
		r.rawFlag = val
		r.Flag = val
	}

	return err
}

var _ RequestFiller = NewChannelSetFlag()

// ChannelRemoveFlag request parameters
type ChannelRemoveFlag struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewChannelRemoveFlag request
func NewChannelRemoveFlag() *ChannelRemoveFlag {
	return &ChannelRemoveFlag{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelRemoveFlag) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewChannelRemoveFlag()

// ChannelRead request parameters
type ChannelRead struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewChannelRead request
func NewChannelRead() *ChannelRead {
	return &ChannelRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewChannelRead()

// ChannelMembers request parameters
type ChannelMembers struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`
}

// NewChannelMembers request
func NewChannelMembers() *ChannelMembers {
	return &ChannelMembers{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelMembers) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	return err
}

var _ RequestFiller = NewChannelMembers()

// ChannelJoin request parameters
type ChannelJoin struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewChannelJoin request
func NewChannelJoin() *ChannelJoin {
	return &ChannelJoin{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelJoin) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewChannelJoin()

// ChannelPart request parameters
type ChannelPart struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewChannelPart request
func NewChannelPart() *ChannelPart {
	return &ChannelPart{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelPart) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewChannelPart()

// ChannelInvite request parameters
type ChannelInvite struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasUserID bool
	rawUserID []string
	UserID    []string
}

// NewChannelInvite request
func NewChannelInvite() *ChannelInvite {
	return &ChannelInvite{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelInvite) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))

	if val, ok := req.Form["userID"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewChannelInvite()

// ChannelAttach request parameters
type ChannelAttach struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasReplyTo bool
	rawReplyTo string
	ReplyTo    uint64 `json:",string"`

	hasUpload bool
	rawUpload string
	Upload    *multipart.FileHeader
}

// NewChannelAttach request
func NewChannelAttach() *ChannelAttach {
	return &ChannelAttach{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelAttach) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["replyTo"] = r.ReplyTo
	out["upload.size"] = r.Upload.Size
	out["upload.filename"] = r.Upload.Filename

	return out
}

// Fill processes request and fills internal variables
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

	r.hasChannelID = true
	r.rawChannelID = chi.URLParam(req, "channelID")
	r.ChannelID = parseUInt64(chi.URLParam(req, "channelID"))
	if val, ok := post["replyTo"]; ok {
		r.hasReplyTo = true
		r.rawReplyTo = val
		r.ReplyTo = parseUInt64(val)
	}
	if _, r.Upload, err = req.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error processing uploaded file")
	}

	return err
}

var _ RequestFiller = NewChannelAttach()

// HasQuery returns true if query was set
func (r *ChannelList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *ChannelList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *ChannelList) GetQuery() string {
	return r.Query
}

// HasName returns true if name was set
func (r *ChannelCreate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ChannelCreate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ChannelCreate) GetName() string {
	return r.Name
}

// HasTopic returns true if topic was set
func (r *ChannelCreate) HasTopic() bool {
	return r.hasTopic
}

// RawTopic returns raw value of topic parameter
func (r *ChannelCreate) RawTopic() string {
	return r.rawTopic
}

// GetTopic returns casted value of  topic parameter
func (r *ChannelCreate) GetTopic() string {
	return r.Topic
}

// HasType returns true if type was set
func (r *ChannelCreate) HasType() bool {
	return r.hasType
}

// RawType returns raw value of type parameter
func (r *ChannelCreate) RawType() string {
	return r.rawType
}

// GetType returns casted value of  type parameter
func (r *ChannelCreate) GetType() string {
	return r.Type
}

// HasMembershipPolicy returns true if membershipPolicy was set
func (r *ChannelCreate) HasMembershipPolicy() bool {
	return r.hasMembershipPolicy
}

// RawMembershipPolicy returns raw value of membershipPolicy parameter
func (r *ChannelCreate) RawMembershipPolicy() string {
	return r.rawMembershipPolicy
}

// GetMembershipPolicy returns casted value of  membershipPolicy parameter
func (r *ChannelCreate) GetMembershipPolicy() types.ChannelMembershipPolicy {
	return r.MembershipPolicy
}

// HasMembers returns true if members was set
func (r *ChannelCreate) HasMembers() bool {
	return r.hasMembers
}

// RawMembers returns raw value of members parameter
func (r *ChannelCreate) RawMembers() []string {
	return r.rawMembers
}

// GetMembers returns casted value of  members parameter
func (r *ChannelCreate) GetMembers() []string {
	return r.Members
}

// HasChannelID returns true if channelID was set
func (r *ChannelUpdate) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelUpdate) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelUpdate) GetChannelID() uint64 {
	return r.ChannelID
}

// HasName returns true if name was set
func (r *ChannelUpdate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ChannelUpdate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ChannelUpdate) GetName() string {
	return r.Name
}

// HasTopic returns true if topic was set
func (r *ChannelUpdate) HasTopic() bool {
	return r.hasTopic
}

// RawTopic returns raw value of topic parameter
func (r *ChannelUpdate) RawTopic() string {
	return r.rawTopic
}

// GetTopic returns casted value of  topic parameter
func (r *ChannelUpdate) GetTopic() string {
	return r.Topic
}

// HasMembershipPolicy returns true if membershipPolicy was set
func (r *ChannelUpdate) HasMembershipPolicy() bool {
	return r.hasMembershipPolicy
}

// RawMembershipPolicy returns raw value of membershipPolicy parameter
func (r *ChannelUpdate) RawMembershipPolicy() string {
	return r.rawMembershipPolicy
}

// GetMembershipPolicy returns casted value of  membershipPolicy parameter
func (r *ChannelUpdate) GetMembershipPolicy() types.ChannelMembershipPolicy {
	return r.MembershipPolicy
}

// HasType returns true if type was set
func (r *ChannelUpdate) HasType() bool {
	return r.hasType
}

// RawType returns raw value of type parameter
func (r *ChannelUpdate) RawType() string {
	return r.rawType
}

// GetType returns casted value of  type parameter
func (r *ChannelUpdate) GetType() string {
	return r.Type
}

// HasOrganisationID returns true if organisationID was set
func (r *ChannelUpdate) HasOrganisationID() bool {
	return r.hasOrganisationID
}

// RawOrganisationID returns raw value of organisationID parameter
func (r *ChannelUpdate) RawOrganisationID() string {
	return r.rawOrganisationID
}

// GetOrganisationID returns casted value of  organisationID parameter
func (r *ChannelUpdate) GetOrganisationID() uint64 {
	return r.OrganisationID
}

// HasChannelID returns true if channelID was set
func (r *ChannelState) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelState) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelState) GetChannelID() uint64 {
	return r.ChannelID
}

// HasState returns true if state was set
func (r *ChannelState) HasState() bool {
	return r.hasState
}

// RawState returns raw value of state parameter
func (r *ChannelState) RawState() string {
	return r.rawState
}

// GetState returns casted value of  state parameter
func (r *ChannelState) GetState() string {
	return r.State
}

// HasChannelID returns true if channelID was set
func (r *ChannelSetFlag) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelSetFlag) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelSetFlag) GetChannelID() uint64 {
	return r.ChannelID
}

// HasFlag returns true if flag was set
func (r *ChannelSetFlag) HasFlag() bool {
	return r.hasFlag
}

// RawFlag returns raw value of flag parameter
func (r *ChannelSetFlag) RawFlag() string {
	return r.rawFlag
}

// GetFlag returns casted value of  flag parameter
func (r *ChannelSetFlag) GetFlag() string {
	return r.Flag
}

// HasChannelID returns true if channelID was set
func (r *ChannelRemoveFlag) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelRemoveFlag) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelRemoveFlag) GetChannelID() uint64 {
	return r.ChannelID
}

// HasChannelID returns true if channelID was set
func (r *ChannelRead) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelRead) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelRead) GetChannelID() uint64 {
	return r.ChannelID
}

// HasChannelID returns true if channelID was set
func (r *ChannelMembers) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelMembers) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelMembers) GetChannelID() uint64 {
	return r.ChannelID
}

// HasChannelID returns true if channelID was set
func (r *ChannelJoin) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelJoin) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelJoin) GetChannelID() uint64 {
	return r.ChannelID
}

// HasUserID returns true if userID was set
func (r *ChannelJoin) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *ChannelJoin) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *ChannelJoin) GetUserID() uint64 {
	return r.UserID
}

// HasChannelID returns true if channelID was set
func (r *ChannelPart) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelPart) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelPart) GetChannelID() uint64 {
	return r.ChannelID
}

// HasUserID returns true if userID was set
func (r *ChannelPart) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *ChannelPart) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *ChannelPart) GetUserID() uint64 {
	return r.UserID
}

// HasChannelID returns true if channelID was set
func (r *ChannelInvite) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelInvite) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelInvite) GetChannelID() uint64 {
	return r.ChannelID
}

// HasUserID returns true if userID was set
func (r *ChannelInvite) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *ChannelInvite) RawUserID() []string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *ChannelInvite) GetUserID() []string {
	return r.UserID
}

// HasChannelID returns true if channelID was set
func (r *ChannelAttach) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ChannelAttach) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ChannelAttach) GetChannelID() uint64 {
	return r.ChannelID
}

// HasReplyTo returns true if replyTo was set
func (r *ChannelAttach) HasReplyTo() bool {
	return r.hasReplyTo
}

// RawReplyTo returns raw value of replyTo parameter
func (r *ChannelAttach) RawReplyTo() string {
	return r.rawReplyTo
}

// GetReplyTo returns casted value of  replyTo parameter
func (r *ChannelAttach) GetReplyTo() uint64 {
	return r.ReplyTo
}

// HasUpload returns true if upload was set
func (r *ChannelAttach) HasUpload() bool {
	return r.hasUpload
}

// RawUpload returns raw value of upload parameter
func (r *ChannelAttach) RawUpload() string {
	return r.rawUpload
}

// GetUpload returns casted value of  upload parameter
func (r *ChannelAttach) GetUpload() *multipart.FileHeader {
	return r.Upload
}
