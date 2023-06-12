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
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/go-chi/chi/v5"
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
	AuthClientList struct {
		// Handle GET parameter
		//
		// Client handle
		Handle string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted clients
		Deleted uint

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// IncTotal GET parameter
		//
		// Include total counter
		IncTotal bool

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	AuthClientCreate struct {
		// Handle POST parameter
		//
		// Client handle
		Handle string

		// Meta POST parameter
		//
		// Additional info
		Meta *types.AuthClientMeta

		// ValidGrant POST parameter
		//
		// Valid grants (authorization_code
		ValidGrant string

		// RedirectURI POST parameter
		//
		// Space delimited list of redirect URIs
		RedirectURI string

		// Scope POST parameter
		//
		// Space delimited list of scopes
		Scope string

		// Trusted POST parameter
		//
		// Is client trusted (skip authorization)
		Trusted bool

		// Enabled POST parameter
		//
		// Is client enabled
		Enabled bool

		// ValidFrom POST parameter
		//
		// Date and time from when client becomes valid
		ValidFrom *time.Time

		// ExpiresAt POST parameter
		//
		// Date and time from client is no logner valid
		ExpiresAt *time.Time

		// Security POST parameter
		//
		// Security settings
		Security *types.AuthClientSecurity

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	AuthClientUpdate struct {
		// ClientID PATH parameter
		//
		// Client ID
		ClientID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Client handle
		Handle string

		// Meta POST parameter
		//
		// Additional info
		Meta *types.AuthClientMeta

		// ValidGrant POST parameter
		//
		// Valid grants (authorization_code
		ValidGrant string

		// RedirectURI POST parameter
		//
		// Space delimited list of redirect URIs
		RedirectURI string

		// Scope POST parameter
		//
		// Space delimited list of scopes
		Scope string

		// Trusted POST parameter
		//
		// Is client trusted (skip authorization)
		Trusted bool

		// Enabled POST parameter
		//
		// Is client enabled
		Enabled bool

		// ValidFrom POST parameter
		//
		// Date and time from when client becomes valid
		ValidFrom *time.Time

		// ExpiresAt POST parameter
		//
		// Date and time from client is no logner valid
		ExpiresAt *time.Time

		// Security POST parameter
		//
		// Security settings
		Security *types.AuthClientSecurity

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	AuthClientRead struct {
		// ClientID PATH parameter
		//
		// Client ID
		ClientID uint64 `json:",string"`
	}

	AuthClientDelete struct {
		// ClientID PATH parameter
		//
		// Client ID
		ClientID uint64 `json:",string"`
	}

	AuthClientUndelete struct {
		// ClientID PATH parameter
		//
		// Client ID
		ClientID uint64 `json:",string"`
	}

	AuthClientRegenerateSecret struct {
		// ClientID PATH parameter
		//
		// Client ID
		ClientID uint64 `json:",string"`
	}

	AuthClientExposeSecret struct {
		// ClientID PATH parameter
		//
		// Client ID
		ClientID uint64 `json:",string"`
	}
)

// NewAuthClientList request
func NewAuthClientList() *AuthClientList {
	return &AuthClientList{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":     r.Handle,
		"deleted":    r.Deleted,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *AuthClientList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
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
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incTotal"]; ok && len(val) > 0 {
			r.IncTotal, err = payload.ParseBool(val[0]), nil
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

// NewAuthClientCreate request
func NewAuthClientCreate() *AuthClientCreate {
	return &AuthClientCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":      r.Handle,
		"meta":        r.Meta,
		"validGrant":  r.ValidGrant,
		"redirectURI": r.RedirectURI,
		"scope":       r.Scope,
		"trusted":     r.Trusted,
		"enabled":     r.Enabled,
		"validFrom":   r.ValidFrom,
		"expiresAt":   r.ExpiresAt,
		"security":    r.Security,
		"labels":      r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetMeta() *types.AuthClientMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetValidGrant() string {
	return r.ValidGrant
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetRedirectURI() string {
	return r.RedirectURI
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetScope() string {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetTrusted() bool {
	return r.Trusted
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetValidFrom() *time.Time {
	return r.ValidFrom
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetExpiresAt() *time.Time {
	return r.ExpiresAt
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetSecurity() *types.AuthClientSecurity {
	return r.Security
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *AuthClientCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseAuthClientMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseAuthClientMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["validGrant"]; ok && len(val) > 0 {
				r.ValidGrant, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["redirectURI"]; ok && len(val) > 0 {
				r.RedirectURI, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["scope"]; ok && len(val) > 0 {
				r.Scope, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["trusted"]; ok && len(val) > 0 {
				r.Trusted, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["enabled"]; ok && len(val) > 0 {
				r.Enabled, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["validFrom"]; ok && len(val) > 0 {
				r.ValidFrom, err = payload.ParseISODatePtrWithErr(val[0])
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["expiresAt"]; ok && len(val) > 0 {
				r.ExpiresAt, err = payload.ParseISODatePtrWithErr(val[0])
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["security[]"]; ok {
				r.Security, err = types.ParseAuthClientSecurity(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["security"]; ok {
				r.Security, err = types.ParseAuthClientSecurity(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["labels[]"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["labels"]; ok {
				r.Labels, err = label.ParseStrings(val)
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseAuthClientMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseAuthClientMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["validGrant"]; ok && len(val) > 0 {
			r.ValidGrant, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["redirectURI"]; ok && len(val) > 0 {
			r.RedirectURI, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scope"]; ok && len(val) > 0 {
			r.Scope, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["trusted"]; ok && len(val) > 0 {
			r.Trusted, err = payload.ParseBool(val[0]), nil
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

		if val, ok := req.Form["validFrom"]; ok && len(val) > 0 {
			r.ValidFrom, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["expiresAt"]; ok && len(val) > 0 {
			r.ExpiresAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["security[]"]; ok {
			r.Security, err = types.ParseAuthClientSecurity(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["security"]; ok {
			r.Security, err = types.ParseAuthClientSecurity(val)
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

	return err
}

// NewAuthClientUpdate request
func NewAuthClientUpdate() *AuthClientUpdate {
	return &AuthClientUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"clientID":    r.ClientID,
		"handle":      r.Handle,
		"meta":        r.Meta,
		"validGrant":  r.ValidGrant,
		"redirectURI": r.RedirectURI,
		"scope":       r.Scope,
		"trusted":     r.Trusted,
		"enabled":     r.Enabled,
		"validFrom":   r.ValidFrom,
		"expiresAt":   r.ExpiresAt,
		"security":    r.Security,
		"labels":      r.Labels,
		"updatedAt":   r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetClientID() uint64 {
	return r.ClientID
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetMeta() *types.AuthClientMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetValidGrant() string {
	return r.ValidGrant
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetRedirectURI() string {
	return r.RedirectURI
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetScope() string {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetTrusted() bool {
	return r.Trusted
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetValidFrom() *time.Time {
	return r.ValidFrom
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetExpiresAt() *time.Time {
	return r.ExpiresAt
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetSecurity() *types.AuthClientSecurity {
	return r.Security
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *AuthClientUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseAuthClientMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseAuthClientMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["validGrant"]; ok && len(val) > 0 {
				r.ValidGrant, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["redirectURI"]; ok && len(val) > 0 {
				r.RedirectURI, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["scope"]; ok && len(val) > 0 {
				r.Scope, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["trusted"]; ok && len(val) > 0 {
				r.Trusted, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["enabled"]; ok && len(val) > 0 {
				r.Enabled, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["validFrom"]; ok && len(val) > 0 {
				r.ValidFrom, err = payload.ParseISODatePtrWithErr(val[0])
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["expiresAt"]; ok && len(val) > 0 {
				r.ExpiresAt, err = payload.ParseISODatePtrWithErr(val[0])
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["security[]"]; ok {
				r.Security, err = types.ParseAuthClientSecurity(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["security"]; ok {
				r.Security, err = types.ParseAuthClientSecurity(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["labels[]"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["labels"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["updatedAt"]; ok && len(val) > 0 {
				r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseAuthClientMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseAuthClientMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["validGrant"]; ok && len(val) > 0 {
			r.ValidGrant, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["redirectURI"]; ok && len(val) > 0 {
			r.RedirectURI, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scope"]; ok && len(val) > 0 {
			r.Scope, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["trusted"]; ok && len(val) > 0 {
			r.Trusted, err = payload.ParseBool(val[0]), nil
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

		if val, ok := req.Form["validFrom"]; ok && len(val) > 0 {
			r.ValidFrom, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["expiresAt"]; ok && len(val) > 0 {
			r.ExpiresAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["security[]"]; ok {
			r.Security, err = types.ParseAuthClientSecurity(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["security"]; ok {
			r.Security, err = types.ParseAuthClientSecurity(val)
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

		val = chi.URLParam(req, "clientID")
		r.ClientID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAuthClientRead request
func NewAuthClientRead() *AuthClientRead {
	return &AuthClientRead{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"clientID": r.ClientID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientRead) GetClientID() uint64 {
	return r.ClientID
}

// Fill processes request and fills internal variables
func (r *AuthClientRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "clientID")
		r.ClientID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAuthClientDelete request
func NewAuthClientDelete() *AuthClientDelete {
	return &AuthClientDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"clientID": r.ClientID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientDelete) GetClientID() uint64 {
	return r.ClientID
}

// Fill processes request and fills internal variables
func (r *AuthClientDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "clientID")
		r.ClientID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAuthClientUndelete request
func NewAuthClientUndelete() *AuthClientUndelete {
	return &AuthClientUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"clientID": r.ClientID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientUndelete) GetClientID() uint64 {
	return r.ClientID
}

// Fill processes request and fills internal variables
func (r *AuthClientUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "clientID")
		r.ClientID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAuthClientRegenerateSecret request
func NewAuthClientRegenerateSecret() *AuthClientRegenerateSecret {
	return &AuthClientRegenerateSecret{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientRegenerateSecret) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"clientID": r.ClientID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientRegenerateSecret) GetClientID() uint64 {
	return r.ClientID
}

// Fill processes request and fills internal variables
func (r *AuthClientRegenerateSecret) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "clientID")
		r.ClientID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAuthClientExposeSecret request
func NewAuthClientExposeSecret() *AuthClientExposeSecret {
	return &AuthClientExposeSecret{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientExposeSecret) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"clientID": r.ClientID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthClientExposeSecret) GetClientID() uint64 {
	return r.ClientID
}

// Fill processes request and fills internal variables
func (r *AuthClientExposeSecret) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "clientID")
		r.ClientID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
