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
	LabelList struct {
		// Kind GET parameter
		//
		// Filter by resource kind (namespace, module, etc)
		Kind string

		// Name GET parameter
		//
		// Filter by label key/name
		Name string

		// Value GET parameter
		//
		// Filter by value
		Value []string

		// Limit GET parameter
		//
		// Limit
		Limit uint
	}

	LabelDelete struct {
		// Name PATH parameter
		//
		// Label name to delete
		Name string

		// Kind GET parameter
		//
		// Limit deletion to a specific resource kind
		Kind string
	}
)

// NewLabelList request
func NewLabelList() *LabelList {
	return &LabelList{}
}

// Auditable returns all auditable/loggable parameters
func (r LabelList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"kind":  r.Kind,
		"name":  r.Name,
		"value": r.Value,
		"limit": r.Limit,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LabelList) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r LabelList) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r LabelList) GetValue() []string {
	return r.Value
}

// Auditable returns all auditable/loggable parameters
func (r LabelList) GetLimit() uint {
	return r.Limit
}

// Fill processes request and fills internal variables
func (r *LabelList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
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
		if val, ok := tmp["value[]"]; ok {
			r.Value, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["value"]; ok {
			r.Value, err = val, nil
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
	}

	return err
}

// NewLabelDelete request
func NewLabelDelete() *LabelDelete {
	return &LabelDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r LabelDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name": r.Name,
		"kind": r.Kind,
	}
}

// Auditable returns all auditable/loggable parameters
func (r LabelDelete) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r LabelDelete) GetKind() string {
	return r.Kind
}

// Fill processes request and fills internal variables
func (r *LabelDelete) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "name")
		r.Name, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
