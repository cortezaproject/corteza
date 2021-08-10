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
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
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
