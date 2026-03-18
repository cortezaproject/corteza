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
	"github.com/cortezaproject/corteza/server/system/types"
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
	LlmProviderList struct {
		// Provider GET parameter
		//
		// Filter by provider type
		Provider string

		// Status GET parameter
		//
		// Filter by status
		Status string
	}

	LlmProviderCreate struct {
		// Handle POST parameter
		//
		// Handle
		Handle string

		// Provider POST parameter
		//
		// Provider type
		Provider string

		// Status POST parameter
		//
		// Status
		Status string

		// ApiKey POST parameter
		//
		// API key
		ApiKey string

		// Meta POST parameter
		//
		// Meta
		Meta types.LLMProviderMeta

		// Config POST parameter
		//
		// Config
		Config types.LLMProviderConfig
	}

	LlmProviderRead struct {
		// LlmProviderID PATH parameter
		//
		// LLM Provider ID
		LlmProviderID uint64 `json:",string"`
	}

	LlmProviderUpdate struct {
		// LlmProviderID PATH parameter
		//
		// LLM Provider ID
		LlmProviderID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Provider POST parameter
		//
		// Provider type
		Provider string

		// Status POST parameter
		//
		// Status
		Status string

		// ApiKey POST parameter
		//
		// API key
		ApiKey string

		// Meta POST parameter
		//
		// Meta
		Meta types.LLMProviderMeta

		// Config POST parameter
		//
		// Config
		Config types.LLMProviderConfig
	}

	LlmProviderDelete struct {
		// LlmProviderID PATH parameter
		//
		// LLM Provider ID
		LlmProviderID uint64 `json:",string"`
	}

	LlmProviderModels struct {
		// LlmProviderID PATH parameter
		//
		// LLM Provider ID
		LlmProviderID uint64 `json:",string"`
	}

	LlmProviderValidate struct {
		// LlmProviderID PATH parameter
		//
		// LLM Provider ID
		LlmProviderID uint64 `json:",string"`
	}
)

// NewLlmProviderList request
func NewLlmProviderList() *LlmProviderList {
	return &LlmProviderList{}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"provider": r.Provider,
		"status":   r.Status,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderList) GetProvider() string {
	return r.Provider
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderList) GetStatus() string {
	return r.Status
}

// Fill processes request and fills internal variables
func (r *LlmProviderList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["provider"]; ok && len(val) > 0 {
			r.Provider, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewLlmProviderCreate request
func NewLlmProviderCreate() *LlmProviderCreate {
	return &LlmProviderCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":   r.Handle,
		"provider": r.Provider,
		"status":   r.Status,
		"apiKey":   r.ApiKey,
		"meta":     r.Meta,
		"config":   r.Config,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderCreate) GetProvider() string {
	return r.Provider
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderCreate) GetApiKey() string {
	return r.ApiKey
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderCreate) GetMeta() types.LLMProviderMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderCreate) GetConfig() types.LLMProviderConfig {
	return r.Config
}

// Fill processes request and fills internal variables
func (r *LlmProviderCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["provider"]; ok && len(val) > 0 {
				r.Provider, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["apiKey"]; ok && len(val) > 0 {
				r.ApiKey, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseLLMProviderMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseLLMProviderMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseLLMProviderConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseLLMProviderConfig(val)
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

		if val, ok := req.Form["provider"]; ok && len(val) > 0 {
			r.Provider, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["apiKey"]; ok && len(val) > 0 {
			r.ApiKey, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseLLMProviderMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseLLMProviderMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseLLMProviderConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseLLMProviderConfig(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewLlmProviderRead request
func NewLlmProviderRead() *LlmProviderRead {
	return &LlmProviderRead{}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"llmProviderID": r.LlmProviderID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderRead) GetLlmProviderID() uint64 {
	return r.LlmProviderID
}

// Fill processes request and fills internal variables
func (r *LlmProviderRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "llmProviderID")
		r.LlmProviderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewLlmProviderUpdate request
func NewLlmProviderUpdate() *LlmProviderUpdate {
	return &LlmProviderUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"llmProviderID": r.LlmProviderID,
		"handle":        r.Handle,
		"provider":      r.Provider,
		"status":        r.Status,
		"apiKey":        r.ApiKey,
		"meta":          r.Meta,
		"config":        r.Config,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderUpdate) GetLlmProviderID() uint64 {
	return r.LlmProviderID
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderUpdate) GetProvider() string {
	return r.Provider
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderUpdate) GetApiKey() string {
	return r.ApiKey
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderUpdate) GetMeta() types.LLMProviderMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderUpdate) GetConfig() types.LLMProviderConfig {
	return r.Config
}

// Fill processes request and fills internal variables
func (r *LlmProviderUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["provider"]; ok && len(val) > 0 {
				r.Provider, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["apiKey"]; ok && len(val) > 0 {
				r.ApiKey, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseLLMProviderMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseLLMProviderMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseLLMProviderConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseLLMProviderConfig(val)
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

		if val, ok := req.Form["provider"]; ok && len(val) > 0 {
			r.Provider, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["apiKey"]; ok && len(val) > 0 {
			r.ApiKey, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseLLMProviderMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseLLMProviderMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseLLMProviderConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseLLMProviderConfig(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "llmProviderID")
		r.LlmProviderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewLlmProviderDelete request
func NewLlmProviderDelete() *LlmProviderDelete {
	return &LlmProviderDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"llmProviderID": r.LlmProviderID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderDelete) GetLlmProviderID() uint64 {
	return r.LlmProviderID
}

// Fill processes request and fills internal variables
func (r *LlmProviderDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "llmProviderID")
		r.LlmProviderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewLlmProviderModels request
func NewLlmProviderModels() *LlmProviderModels {
	return &LlmProviderModels{}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderModels) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"llmProviderID": r.LlmProviderID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderModels) GetLlmProviderID() uint64 {
	return r.LlmProviderID
}

// Fill processes request and fills internal variables
func (r *LlmProviderModels) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "llmProviderID")
		r.LlmProviderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewLlmProviderValidate request
func NewLlmProviderValidate() *LlmProviderValidate {
	return &LlmProviderValidate{}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderValidate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"llmProviderID": r.LlmProviderID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LlmProviderValidate) GetLlmProviderID() uint64 {
	return r.LlmProviderID
}

// Fill processes request and fills internal variables
func (r *LlmProviderValidate) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "llmProviderID")
		r.LlmProviderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
