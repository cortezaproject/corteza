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
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	ModuleList struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// Query GET parameter
		//
		// Search query
		Query string

		// Name GET parameter
		//
		// Search by name
		Name string

		// Handle GET parameter
		//
		// Search by handle
		Handle string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	ModuleCreate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// Name POST parameter
		//
		// Module Name
		Name string

		// Handle POST parameter
		//
		// Module handle
		Handle string

		// Fields POST parameter
		//
		// Fields JSON
		Fields types.ModuleFieldSet

		// Meta POST parameter
		//
		// Module meta data
		Meta sqlxTypes.JSONText

		// Labels POST parameter
		//
		// Module labels
		Labels map[string]string
	}

	ModuleRead struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`
	}

	ModuleUpdate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Name POST parameter
		//
		// Module Name
		Name string

		// Handle POST parameter
		//
		// Module Handle
		Handle string

		// Fields POST parameter
		//
		// Fields JSON
		Fields types.ModuleFieldSet

		// Meta POST parameter
		//
		// Module meta data
		Meta sqlxTypes.JSONText

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time

		// Labels POST parameter
		//
		// Module labels
		Labels map[string]string
	}

	ModuleDelete struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`
	}

	ModuleTriggerScript struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// ID
		ModuleID uint64 `json:",string"`

		// Script POST parameter
		//
		// Script to execute
		Script string
	}

	ModuleListTranslations struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// ID
		ModuleID uint64 `json:",string"`
	}

	ModuleUpdateTranslations struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// ID
		ModuleID uint64 `json:",string"`

		// Translation POST parameter
		//
		// Resource translation to upsert
		Translation locale.ResourceTranslationSet
	}
)

// NewModuleList request
func NewModuleList() *ModuleList {
	return &ModuleList{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"query":       r.Query,
		"name":        r.Name,
		"handle":      r.Handle,
		"limit":       r.Limit,
		"pageCursor":  r.PageCursor,
		"labels":      r.Labels,
		"sort":        r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ModuleList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := tmp["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["sort"]; ok && len(val) > 0 {
			r.Sort, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewModuleCreate request
func NewModuleCreate() *ModuleCreate {
	return &ModuleCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"name":        r.Name,
		"handle":      r.Handle,
		"fields":      r.Fields,
		"meta":        r.Meta,
		"labels":      r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ModuleCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ModuleCreate) GetFields() types.ModuleFieldSet {
	return r.Fields
}

// Auditable returns all auditable/loggable parameters
func (r ModuleCreate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r ModuleCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ModuleCreate) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["fields[]"]; ok && len(val) > 0  {
		//    r.Fields, err = types.ModuleFieldSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["meta"]; ok && len(val) > 0 {
			r.Meta, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewModuleRead request
func NewModuleRead() *ModuleRead {
	return &ModuleRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleRead) GetModuleID() uint64 {
	return r.ModuleID
}

// Fill processes request and fills internal variables
func (r *ModuleRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewModuleUpdate request
func NewModuleUpdate() *ModuleUpdate {
	return &ModuleUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"name":        r.Name,
		"handle":      r.Handle,
		"fields":      r.Fields,
		"meta":        r.Meta,
		"updatedAt":   r.UpdatedAt,
		"labels":      r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) GetFields() types.ModuleFieldSet {
	return r.Fields
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ModuleUpdate) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["fields[]"]; ok && len(val) > 0  {
		//    r.Fields, err = types.ModuleFieldSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["meta"]; ok && len(val) > 0 {
			r.Meta, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["updatedAt"]; ok && len(val) > 0 {
			r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewModuleDelete request
func NewModuleDelete() *ModuleDelete {
	return &ModuleDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleDelete) GetModuleID() uint64 {
	return r.ModuleID
}

// Fill processes request and fills internal variables
func (r *ModuleDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewModuleTriggerScript request
func NewModuleTriggerScript() *ModuleTriggerScript {
	return &ModuleTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleTriggerScript) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"script":      r.Script,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleTriggerScript) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleTriggerScript) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleTriggerScript) GetScript() string {
	return r.Script
}

// Fill processes request and fills internal variables
func (r *ModuleTriggerScript) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["script"]; ok && len(val) > 0 {
			r.Script, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewModuleListTranslations request
func NewModuleListTranslations() *ModuleListTranslations {
	return &ModuleListTranslations{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleListTranslations) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleListTranslations) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleListTranslations) GetModuleID() uint64 {
	return r.ModuleID
}

// Fill processes request and fills internal variables
func (r *ModuleListTranslations) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewModuleUpdateTranslations request
func NewModuleUpdateTranslations() *ModuleUpdateTranslations {
	return &ModuleUpdateTranslations{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdateTranslations) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"translation": r.Translation,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdateTranslations) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdateTranslations) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdateTranslations) GetTranslation() locale.ResourceTranslationSet {
	return r.Translation
}

// Fill processes request and fills internal variables
func (r *ModuleUpdateTranslations) Fill(req *http.Request) (err error) {

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

		//if val, ok := req.Form["translation[]"]; ok && len(val) > 0  {
		//    r.Translation, err = locale.ResourceTranslationSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
