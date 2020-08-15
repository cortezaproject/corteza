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
	StatusList struct {
	}

	StatusSet struct {
		// Icon POST parameter
		//
		// Status icon
		Icon string

		// Message POST parameter
		//
		// Status message
		Message string

		// Expires POST parameter
		//
		// Clear status when it expires (eg: when-active, afternoon, tomorrow 1h, 30m, 1 PM, 2019-05-20)
		Expires string
	}

	StatusDelete struct {
	}
)

// NewStatusList request
func NewStatusList() *StatusList {
	return &StatusList{}
}

// Auditable returns all auditable/loggable parameters
func (r StatusList) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *StatusList) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	return err
}

// NewStatusSet request
func NewStatusSet() *StatusSet {
	return &StatusSet{}
}

// Auditable returns all auditable/loggable parameters
func (r StatusSet) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"icon":    r.Icon,
		"message": r.Message,
		"expires": r.Expires,
	}
}

// Auditable returns all auditable/loggable parameters
func (r StatusSet) GetIcon() string {
	return r.Icon
}

// Auditable returns all auditable/loggable parameters
func (r StatusSet) GetMessage() string {
	return r.Message
}

// Auditable returns all auditable/loggable parameters
func (r StatusSet) GetExpires() string {
	return r.Expires
}

// Fill processes request and fills internal variables
func (r *StatusSet) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["icon"]; ok && len(val) > 0 {
			r.Icon, err = val[0], nil
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

		if val, ok := req.Form["expires"]; ok && len(val) > 0 {
			r.Expires, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewStatusDelete request
func NewStatusDelete() *StatusDelete {
	return &StatusDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r StatusDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *StatusDelete) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	return err
}
