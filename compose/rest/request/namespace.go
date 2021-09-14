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
	NamespaceList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// Slug GET parameter
		//
		// Search by namespace slug
		Slug string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	NamespaceCreate struct {
		// Name POST parameter
		//
		// Name
		Name string

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Slug POST parameter
		//
		// Slug (url path part)
		Slug string

		// Enabled POST parameter
		//
		// Enabled
		Enabled bool

		// Meta POST parameter
		//
		// Meta data
		Meta sqlxTypes.JSONText
	}

	NamespaceRead struct {
		// NamespaceID PATH parameter
		//
		// ID
		NamespaceID uint64 `json:",string"`
	}

	NamespaceUpdate struct {
		// NamespaceID PATH parameter
		//
		// ID
		NamespaceID uint64 `json:",string"`

		// Name POST parameter
		//
		// Name
		Name string

		// Slug POST parameter
		//
		// Slug (url path part)
		Slug string

		// Enabled POST parameter
		//
		// Enabled
		Enabled bool

		// Meta POST parameter
		//
		// Meta data
		Meta sqlxTypes.JSONText

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	NamespaceDelete struct {
		// NamespaceID PATH parameter
		//
		// ID
		NamespaceID uint64 `json:",string"`
	}

	NamespaceUpload struct {
		// Upload POST parameter
		//
		// File to upload
		Upload *multipart.FileHeader
	}

	NamespaceTriggerScript struct {
		// NamespaceID PATH parameter
		//
		// ID
		NamespaceID uint64 `json:",string"`

		// Script POST parameter
		//
		// Script to execute
		Script string
	}

	NamespaceListLocale struct {
		// NamespaceID PATH parameter
		//
		// ID
		NamespaceID uint64 `json:",string"`
	}

	NamespaceUpdateLocale struct {
		// NamespaceID PATH parameter
		//
		// ID
		NamespaceID uint64 `json:",string"`

		// Locale POST parameter
		//
		// ...
		Locale locale.ResourceTranslationSet
	}
)

// NewNamespaceList request
func NewNamespaceList() *NamespaceList {
	return &NamespaceList{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"slug":       r.Slug,
		"limit":      r.Limit,
		"labels":     r.Labels,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceList) GetSlug() string {
	return r.Slug
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *NamespaceList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["slug"]; ok && len(val) > 0 {
			r.Slug, err = val[0], nil
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
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
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

	return err
}

// NewNamespaceCreate request
func NewNamespaceCreate() *NamespaceCreate {
	return &NamespaceCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":    r.Name,
		"labels":  r.Labels,
		"slug":    r.Slug,
		"enabled": r.Enabled,
		"meta":    r.Meta,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceCreate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceCreate) GetSlug() string {
	return r.Slug
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceCreate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// Fill processes request and fills internal variables
func (r *NamespaceCreate) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["slug"]; ok && len(val) > 0 {
			r.Slug, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["enabled"]; ok && len(val) > 0 {
			r.Enabled, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta"]; ok && len(val) > 0 {
			r.Meta, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewNamespaceRead request
func NewNamespaceRead() *NamespaceRead {
	return &NamespaceRead{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Fill processes request and fills internal variables
func (r *NamespaceRead) Fill(req *http.Request) (err error) {

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

// NewNamespaceUpdate request
func NewNamespaceUpdate() *NamespaceUpdate {
	return &NamespaceUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"name":        r.Name,
		"slug":        r.Slug,
		"enabled":     r.Enabled,
		"meta":        r.Meta,
		"labels":      r.Labels,
		"updatedAt":   r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) GetSlug() string {
	return r.Slug
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *NamespaceUpdate) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["slug"]; ok && len(val) > 0 {
			r.Slug, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["enabled"]; ok && len(val) > 0 {
			r.Enabled, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

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

		if val, ok := req.Form["updatedAt"]; ok && len(val) > 0 {
			r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
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

// NewNamespaceDelete request
func NewNamespaceDelete() *NamespaceDelete {
	return &NamespaceDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Fill processes request and fills internal variables
func (r *NamespaceDelete) Fill(req *http.Request) (err error) {

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

// NewNamespaceUpload request
func NewNamespaceUpload() *NamespaceUpload {
	return &NamespaceUpload{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpload) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"upload": r.Upload,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpload) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Fill processes request and fills internal variables
func (r *NamespaceUpload) Fill(req *http.Request) (err error) {

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

		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

	}

	return err
}

// NewNamespaceTriggerScript request
func NewNamespaceTriggerScript() *NamespaceTriggerScript {
	return &NamespaceTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceTriggerScript) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"script":      r.Script,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceTriggerScript) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceTriggerScript) GetScript() string {
	return r.Script
}

// Fill processes request and fills internal variables
func (r *NamespaceTriggerScript) Fill(req *http.Request) (err error) {

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

	}

	return err
}

// NewNamespaceListLocale request
func NewNamespaceListLocale() *NamespaceListLocale {
	return &NamespaceListLocale{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceListLocale) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceListLocale) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Fill processes request and fills internal variables
func (r *NamespaceListLocale) Fill(req *http.Request) (err error) {

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

// NewNamespaceUpdateLocale request
func NewNamespaceUpdateLocale() *NamespaceUpdateLocale {
	return &NamespaceUpdateLocale{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdateLocale) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"locale":      r.Locale,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdateLocale) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdateLocale) GetLocale() locale.ResourceTranslationSet {
	return r.Locale
}

// Fill processes request and fills internal variables
func (r *NamespaceUpdateLocale) Fill(req *http.Request) (err error) {

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

		//if val, ok := req.Form["locale[]"]; ok && len(val) > 0  {
		//    r.Locale, err = locale.ResourceTranslationSet(val), nil
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

	}

	return err
}
