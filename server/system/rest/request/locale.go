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
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/go-chi/chi/v5"
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
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	LocaleListResource struct {
		// Lang GET parameter
		//
		// Language
		Lang string

		// Resource GET parameter
		//
		// Resource
		Resource string

		// ResourceType GET parameter
		//
		// Resource type
		ResourceType string

		// OwnerID GET parameter
		//
		// OwnerID
		OwnerID uint64 `json:",string"`

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted resource translations
		Deleted uint64 `json:",string"`

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	LocaleCreateResource struct {
		// Lang POST parameter
		//
		// Lang
		Lang string

		// Resource POST parameter
		//
		// Resource
		Resource string

		// Key POST parameter
		//
		// Key
		Key string

		// Place POST parameter
		//
		// place
		Place int

		// Message POST parameter
		//
		// Message
		Message string

		// OwnerID POST parameter
		//
		// OwnerID
		OwnerID uint64 `json:",string"`
	}

	LocaleUpdateResource struct {
		// TranslationID PATH parameter
		//
		// ID
		TranslationID uint64 `json:",string"`

		// Lang POST parameter
		//
		// Lang
		Lang string

		// Resource POST parameter
		//
		// Resource
		Resource string

		// Key POST parameter
		//
		// Key
		Key string

		// Place POST parameter
		//
		// place
		Place int

		// Message POST parameter
		//
		// Message
		Message string

		// OwnerID POST parameter
		//
		// OwnerID
		OwnerID uint64 `json:",string"`
	}

	LocaleReadResource struct {
		// TranslationID PATH parameter
		//
		// ID
		TranslationID uint64 `json:",string"`
	}

	LocaleDeleteResource struct {
		// TranslationID PATH parameter
		//
		// ID
		TranslationID uint64 `json:",string"`
	}

	LocaleUndeleteResource struct {
		// TranslationID PATH parameter
		//
		// ID
		TranslationID uint64 `json:",string"`
	}

	LocaleList struct {
	}

	LocaleGet struct {
		// Lang PATH parameter
		//
		// Language
		Lang string

		// Application PATH parameter
		//
		// Application name
		Application string
	}
)

// NewLocaleListResource request
func NewLocaleListResource() *LocaleListResource {
	return &LocaleListResource{}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"lang":         r.Lang,
		"resource":     r.Resource,
		"resourceType": r.ResourceType,
		"ownerID":      r.OwnerID,
		"deleted":      r.Deleted,
		"limit":        r.Limit,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) GetLang() string {
	return r.Lang
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) GetOwnerID() uint64 {
	return r.OwnerID
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) GetDeleted() uint64 {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r LocaleListResource) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *LocaleListResource) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["lang"]; ok && len(val) > 0 {
			r.Lang, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["resourceType"]; ok && len(val) > 0 {
			r.ResourceType, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["ownerID"]; ok && len(val) > 0 {
			r.OwnerID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint64(val[0]), nil
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
		if val, ok := tmp["sort"]; ok && len(val) > 0 {
			r.Sort, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewLocaleCreateResource request
func NewLocaleCreateResource() *LocaleCreateResource {
	return &LocaleCreateResource{}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleCreateResource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"lang":     r.Lang,
		"resource": r.Resource,
		"key":      r.Key,
		"place":    r.Place,
		"message":  r.Message,
		"ownerID":  r.OwnerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleCreateResource) GetLang() string {
	return r.Lang
}

// Auditable returns all auditable/loggable parameters
func (r LocaleCreateResource) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r LocaleCreateResource) GetKey() string {
	return r.Key
}

// Auditable returns all auditable/loggable parameters
func (r LocaleCreateResource) GetPlace() int {
	return r.Place
}

// Auditable returns all auditable/loggable parameters
func (r LocaleCreateResource) GetMessage() string {
	return r.Message
}

// Auditable returns all auditable/loggable parameters
func (r LocaleCreateResource) GetOwnerID() uint64 {
	return r.OwnerID
}

// Fill processes request and fills internal variables
func (r *LocaleCreateResource) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["lang"]; ok && len(val) > 0 {
				r.Lang, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["resource"]; ok && len(val) > 0 {
				r.Resource, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["key"]; ok && len(val) > 0 {
				r.Key, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["place"]; ok && len(val) > 0 {
				r.Place, err = payload.ParseInt(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["message"]; ok && len(val) > 0 {
				r.Message, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["ownerID"]; ok && len(val) > 0 {
				r.OwnerID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["lang"]; ok && len(val) > 0 {
			r.Lang, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["key"]; ok && len(val) > 0 {
			r.Key, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["place"]; ok && len(val) > 0 {
			r.Place, err = payload.ParseInt(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["message"]; ok && len(val) > 0 {
			r.Message, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownerID"]; ok && len(val) > 0 {
			r.OwnerID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewLocaleUpdateResource request
func NewLocaleUpdateResource() *LocaleUpdateResource {
	return &LocaleUpdateResource{}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUpdateResource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"translationID": r.TranslationID,
		"lang":          r.Lang,
		"resource":      r.Resource,
		"key":           r.Key,
		"place":         r.Place,
		"message":       r.Message,
		"ownerID":       r.OwnerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUpdateResource) GetTranslationID() uint64 {
	return r.TranslationID
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUpdateResource) GetLang() string {
	return r.Lang
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUpdateResource) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUpdateResource) GetKey() string {
	return r.Key
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUpdateResource) GetPlace() int {
	return r.Place
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUpdateResource) GetMessage() string {
	return r.Message
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUpdateResource) GetOwnerID() uint64 {
	return r.OwnerID
}

// Fill processes request and fills internal variables
func (r *LocaleUpdateResource) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["lang"]; ok && len(val) > 0 {
				r.Lang, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["resource"]; ok && len(val) > 0 {
				r.Resource, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["key"]; ok && len(val) > 0 {
				r.Key, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["place"]; ok && len(val) > 0 {
				r.Place, err = payload.ParseInt(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["message"]; ok && len(val) > 0 {
				r.Message, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["ownerID"]; ok && len(val) > 0 {
				r.OwnerID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["lang"]; ok && len(val) > 0 {
			r.Lang, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["key"]; ok && len(val) > 0 {
			r.Key, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["place"]; ok && len(val) > 0 {
			r.Place, err = payload.ParseInt(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["message"]; ok && len(val) > 0 {
			r.Message, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownerID"]; ok && len(val) > 0 {
			r.OwnerID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "translationID")
		r.TranslationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewLocaleReadResource request
func NewLocaleReadResource() *LocaleReadResource {
	return &LocaleReadResource{}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleReadResource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"translationID": r.TranslationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleReadResource) GetTranslationID() uint64 {
	return r.TranslationID
}

// Fill processes request and fills internal variables
func (r *LocaleReadResource) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "translationID")
		r.TranslationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewLocaleDeleteResource request
func NewLocaleDeleteResource() *LocaleDeleteResource {
	return &LocaleDeleteResource{}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleDeleteResource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"translationID": r.TranslationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleDeleteResource) GetTranslationID() uint64 {
	return r.TranslationID
}

// Fill processes request and fills internal variables
func (r *LocaleDeleteResource) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "translationID")
		r.TranslationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewLocaleUndeleteResource request
func NewLocaleUndeleteResource() *LocaleUndeleteResource {
	return &LocaleUndeleteResource{}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUndeleteResource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"translationID": r.TranslationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleUndeleteResource) GetTranslationID() uint64 {
	return r.TranslationID
}

// Fill processes request and fills internal variables
func (r *LocaleUndeleteResource) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "translationID")
		r.TranslationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewLocaleList request
func NewLocaleList() *LocaleList {
	return &LocaleList{}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleList) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *LocaleList) Fill(req *http.Request) (err error) {

	return err
}

// NewLocaleGet request
func NewLocaleGet() *LocaleGet {
	return &LocaleGet{}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleGet) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"lang":        r.Lang,
		"application": r.Application,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LocaleGet) GetLang() string {
	return r.Lang
}

// Auditable returns all auditable/loggable parameters
func (r LocaleGet) GetApplication() string {
	return r.Application
}

// Fill processes request and fills internal variables
func (r *LocaleGet) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "lang")
		r.Lang, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "application")
		r.Application, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
