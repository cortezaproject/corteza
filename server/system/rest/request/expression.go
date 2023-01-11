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
	ExpressionEvaluate struct {
		// Variables POST parameter
		//
		// variables
		Variables map[string]interface{}

		// Expressions POST parameter
		//
		// expressions
		Expressions map[string]string
	}
)

// NewExpressionEvaluate request
func NewExpressionEvaluate() *ExpressionEvaluate {
	return &ExpressionEvaluate{}
}

// Auditable returns all auditable/loggable parameters
func (r ExpressionEvaluate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"variables":   r.Variables,
		"expressions": r.Expressions,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ExpressionEvaluate) GetVariables() map[string]interface{} {
	return r.Variables
}

// Auditable returns all auditable/loggable parameters
func (r ExpressionEvaluate) GetExpressions() map[string]string {
	return r.Expressions
}

// Fill processes request and fills internal variables
func (r *ExpressionEvaluate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["variables[]"]; ok {
				r.Variables, err = parseMapStringInterface(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["variables"]; ok {
				r.Variables, err = parseMapStringInterface(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["expressions[]"]; ok {
				r.Expressions, err = parseMapStringString(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["expressions"]; ok {
				r.Expressions, err = parseMapStringString(val)
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

		if val, ok := req.Form["variables[]"]; ok {
			r.Variables, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["variables"]; ok {
			r.Variables, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["expressions[]"]; ok {
			r.Expressions, err = parseMapStringString(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["expressions"]; ok {
			r.Expressions, err = parseMapStringString(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}
