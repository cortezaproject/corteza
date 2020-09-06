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
	ModuleRead struct {
		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`
	}
)

// NewModuleRead request
func NewModuleRead() *ModuleRead {
	return &ModuleRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"moduleID": r.ModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleRead) GetModuleID() uint64 {
	return r.ModuleID
}

// Fill processes request and fills internal variables
func (r *ModuleRead) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
