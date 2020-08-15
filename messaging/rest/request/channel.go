package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
)

type (
	// Internal API interface
	ChannelList struct {
		// Query GET parameter
		//
		// Search query
		Query string
	}

	ChannelCreate struct {
		// Name POST parameter
		//
		// Name of Channel
		Name string

		// Topic POST parameter
		//
		// Subject of Channel
		Topic string

		// Type POST parameter
		//
		// Channel type
		Type string

		// MembershipPolicy POST parameter
		//
		// Membership policy (eg: featured, forced)?
		MembershipPolicy types.ChannelMembershipPolicy

		// Members POST parameter
		//
		// Initial members of the channel
		Members []string
	}

	ChannelUpdate struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// Name POST parameter
		//
		// Name of Channel
		Name string

		// Topic POST parameter
		//
		// Subject of Channel
		Topic string

		// MembershipPolicy POST parameter
		//
		// Membership policy (eg: featured, forced)?
		MembershipPolicy types.ChannelMembershipPolicy

		// Type POST parameter
		//
		// Channel type
		Type string

		// OrganisationID POST parameter
		//
		// Move channel to different organisation
		OrganisationID uint64 `json:",string"`
	}

	ChannelState struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// State POST parameter
		//
		// Valid values: delete, undelete, archive, unarchive
		State string
	}

	ChannelSetFlag struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// Flag POST parameter
		//
		// Valid values: pinned, hidden, ignored
		Flag string
	}

	ChannelRemoveFlag struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`
	}

	ChannelRead struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`
	}

	ChannelMembers struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`
	}

	ChannelJoin struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// Member ID
		UserID uint64 `json:",string"`
	}

	ChannelPart struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// Member ID
		UserID uint64 `json:",string"`
	}

	ChannelInvite struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// UserID POST parameter
		//
		// User ID
		UserID []string
	}

	ChannelAttach struct {
		// ChannelID PATH parameter
		//
		// Channel ID
		ChannelID uint64 `json:",string"`

		// ReplyTo POST parameter
		//
		// Upload as a reply
		ReplyTo uint64 `json:",string"`

		// Upload POST parameter
		//
		// File to upload
		Upload *multipart.FileHeader
	}
)

// NewChannelList request
func NewChannelList() *ChannelList {
	return &ChannelList{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query": r.Query,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelList) GetQuery() string {
	return r.Query
}

// Fill processes request and fills internal variables
func (r *ChannelList) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewChannelCreate request
func NewChannelCreate() *ChannelCreate {
	return &ChannelCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":             r.Name,
		"topic":            r.Topic,
		"type":             r.Type,
		"membershipPolicy": r.MembershipPolicy,
		"members":          r.Members,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ChannelCreate) GetTopic() string {
	return r.Topic
}

// Auditable returns all auditable/loggable parameters
func (r ChannelCreate) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r ChannelCreate) GetMembershipPolicy() types.ChannelMembershipPolicy {
	return r.MembershipPolicy
}

// Auditable returns all auditable/loggable parameters
func (r ChannelCreate) GetMembers() []string {
	return r.Members
}

// Fill processes request and fills internal variables
func (r *ChannelCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["topic"]; ok && len(val) > 0 {
			r.Topic, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["membershipPolicy"]; ok && len(val) > 0 {
			r.MembershipPolicy, err = types.ChannelMembershipPolicy(val[0]), nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["members[]"]; ok && len(val) > 0  {
		//    r.Members, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	return err
}

// NewChannelUpdate request
func NewChannelUpdate() *ChannelUpdate {
	return &ChannelUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID":        r.ChannelID,
		"name":             r.Name,
		"topic":            r.Topic,
		"membershipPolicy": r.MembershipPolicy,
		"type":             r.Type,
		"organisationID":   r.OrganisationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelUpdate) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r ChannelUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ChannelUpdate) GetTopic() string {
	return r.Topic
}

// Auditable returns all auditable/loggable parameters
func (r ChannelUpdate) GetMembershipPolicy() types.ChannelMembershipPolicy {
	return r.MembershipPolicy
}

// Auditable returns all auditable/loggable parameters
func (r ChannelUpdate) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r ChannelUpdate) GetOrganisationID() uint64 {
	return r.OrganisationID
}

// Fill processes request and fills internal variables
func (r *ChannelUpdate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["topic"]; ok && len(val) > 0 {
			r.Topic, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["membershipPolicy"]; ok && len(val) > 0 {
			r.MembershipPolicy, err = types.ChannelMembershipPolicy(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["organisationID"]; ok && len(val) > 0 {
			r.OrganisationID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelState request
func NewChannelState() *ChannelState {
	return &ChannelState{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelState) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"state":     r.State,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelState) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r ChannelState) GetState() string {
	return r.State
}

// Fill processes request and fills internal variables
func (r *ChannelState) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["state"]; ok && len(val) > 0 {
			r.State, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelSetFlag request
func NewChannelSetFlag() *ChannelSetFlag {
	return &ChannelSetFlag{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelSetFlag) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"flag":      r.Flag,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelSetFlag) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r ChannelSetFlag) GetFlag() string {
	return r.Flag
}

// Fill processes request and fills internal variables
func (r *ChannelSetFlag) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["flag"]; ok && len(val) > 0 {
			r.Flag, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelRemoveFlag request
func NewChannelRemoveFlag() *ChannelRemoveFlag {
	return &ChannelRemoveFlag{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelRemoveFlag) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelRemoveFlag) GetChannelID() uint64 {
	return r.ChannelID
}

// Fill processes request and fills internal variables
func (r *ChannelRemoveFlag) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelRead request
func NewChannelRead() *ChannelRead {
	return &ChannelRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelRead) GetChannelID() uint64 {
	return r.ChannelID
}

// Fill processes request and fills internal variables
func (r *ChannelRead) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelMembers request
func NewChannelMembers() *ChannelMembers {
	return &ChannelMembers{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelMembers) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelMembers) GetChannelID() uint64 {
	return r.ChannelID
}

// Fill processes request and fills internal variables
func (r *ChannelMembers) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelJoin request
func NewChannelJoin() *ChannelJoin {
	return &ChannelJoin{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelJoin) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"userID":    r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelJoin) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r ChannelJoin) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *ChannelJoin) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelPart request
func NewChannelPart() *ChannelPart {
	return &ChannelPart{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelPart) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"userID":    r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelPart) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r ChannelPart) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *ChannelPart) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelInvite request
func NewChannelInvite() *ChannelInvite {
	return &ChannelInvite{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelInvite) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"userID":    r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelInvite) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r ChannelInvite) GetUserID() []string {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *ChannelInvite) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		//if val, ok := req.Form["userID[]"]; ok && len(val) > 0  {
		//    r.UserID, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChannelAttach request
func NewChannelAttach() *ChannelAttach {
	return &ChannelAttach{}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelAttach) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"replyTo":   r.ReplyTo,
		"upload":    r.Upload,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChannelAttach) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r ChannelAttach) GetReplyTo() uint64 {
	return r.ReplyTo
}

// Auditable returns all auditable/loggable parameters
func (r ChannelAttach) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Fill processes request and fills internal variables
func (r *ChannelAttach) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["replyTo"]; ok && len(val) > 0 {
			r.ReplyTo, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "channelID")
		r.ChannelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
