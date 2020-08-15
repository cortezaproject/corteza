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
	OrganisationList struct {
		// Query GET parameter
		//
		// Search query
		Query string
	}

	OrganisationCreate struct {
		// Name POST parameter
		//
		// Organisation Name
		Name string
	}

	OrganisationUpdate struct {
		// Id PATH parameter
		//
		// Organisation ID
		Id uint64 `json:",string"`

		// Name POST parameter
		//
		// Organisation Name
		Name string
	}

	OrganisationDelete struct {
		// Id PATH parameter
		//
		// Organisation ID
		Id uint64 `json:",string"`
	}

	OrganisationRead struct {
		// Id PATH parameter
		//
		// Organisation ID
		Id uint64 `json:",string"`
	}

	OrganisationArchive struct {
		// Id PATH parameter
		//
		// Organisation ID
		Id uint64 `json:",string"`
	}
)

// NewOrganisationList request
func NewOrganisationList() *OrganisationList {
	return &OrganisationList{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query": r.Query,
	}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationList) GetQuery() string {
	return r.Query
}

// Fill processes request and fills internal variables
func (r *OrganisationList) Fill(req *http.Request) (err error) {
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
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewOrganisationCreate request
func NewOrganisationCreate() *OrganisationCreate {
	return &OrganisationCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name": r.Name,
	}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationCreate) GetName() string {
	return r.Name
}

// Fill processes request and fills internal variables
func (r *OrganisationCreate) Fill(req *http.Request) (err error) {
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
	}

	return err
}

// NewOrganisationUpdate request
func NewOrganisationUpdate() *OrganisationUpdate {
	return &OrganisationUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"id":   r.Id,
		"name": r.Name,
	}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationUpdate) GetId() uint64 {
	return r.Id
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationUpdate) GetName() string {
	return r.Name
}

// Fill processes request and fills internal variables
func (r *OrganisationUpdate) Fill(req *http.Request) (err error) {
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
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "id")
		r.Id, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewOrganisationDelete request
func NewOrganisationDelete() *OrganisationDelete {
	return &OrganisationDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"id": r.Id,
	}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationDelete) GetId() uint64 {
	return r.Id
}

// Fill processes request and fills internal variables
func (r *OrganisationDelete) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "id")
		r.Id, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewOrganisationRead request
func NewOrganisationRead() *OrganisationRead {
	return &OrganisationRead{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"id": r.Id,
	}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationRead) GetId() uint64 {
	return r.Id
}

// Fill processes request and fills internal variables
func (r *OrganisationRead) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "id")
		r.Id, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewOrganisationArchive request
func NewOrganisationArchive() *OrganisationArchive {
	return &OrganisationArchive{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationArchive) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"id": r.Id,
	}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationArchive) GetId() uint64 {
	return r.Id
}

// Fill processes request and fills internal variables
func (r *OrganisationArchive) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "id")
		r.Id, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
