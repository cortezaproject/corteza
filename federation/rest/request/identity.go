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
	IdentityGenerateNodeIdentity struct {
		// Domain POST parameter
		//
		// Domain of the destination node
		Domain string
	}

	IdentityRegisterOriginNode struct {
		// Identifier POST parameter
		//
		// Origin node identifier
		Identifier string
	}
)

// NewIdentityGenerateNodeIdentity request
func NewIdentityGenerateNodeIdentity() *IdentityGenerateNodeIdentity {
	return &IdentityGenerateNodeIdentity{}
}

// Auditable returns all auditable/loggable parameters
func (r IdentityGenerateNodeIdentity) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"domain": r.Domain,
	}
}

// Auditable returns all auditable/loggable parameters
func (r IdentityGenerateNodeIdentity) GetDomain() string {
	return r.Domain
}

// Fill processes request and fills internal variables
func (r *IdentityGenerateNodeIdentity) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["domain"]; ok && len(val) > 0 {
			r.Domain, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewIdentityRegisterOriginNode request
func NewIdentityRegisterOriginNode() *IdentityRegisterOriginNode {
	return &IdentityRegisterOriginNode{}
}

// Auditable returns all auditable/loggable parameters
func (r IdentityRegisterOriginNode) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"identifier": r.Identifier,
	}
}

// Auditable returns all auditable/loggable parameters
func (r IdentityRegisterOriginNode) GetIdentifier() string {
	return r.Identifier
}

// Fill processes request and fills internal variables
func (r *IdentityRegisterOriginNode) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["identifier"]; ok && len(val) > 0 {
			r.Identifier, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
