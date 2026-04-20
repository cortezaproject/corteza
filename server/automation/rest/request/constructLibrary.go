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
	"github.com/crusttech/human/server/pkg/payload"
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
	ConstructLibraryFunctions struct {
	}

	ConstructLibraryTriggers struct {
	}
)

// NewConstructLibraryFunctions request
func NewConstructLibraryFunctions() *ConstructLibraryFunctions {
	return &ConstructLibraryFunctions{}
}

// Auditable returns all auditable/loggable parameters
func (r ConstructLibraryFunctions) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *ConstructLibraryFunctions) Fill(req *http.Request) (err error) {

	return err
}

// NewConstructLibraryTriggers request
func NewConstructLibraryTriggers() *ConstructLibraryTriggers {
	return &ConstructLibraryTriggers{}
}

// Auditable returns all auditable/loggable parameters
func (r ConstructLibraryTriggers) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *ConstructLibraryTriggers) Fill(req *http.Request) (err error) {

	return err
}
