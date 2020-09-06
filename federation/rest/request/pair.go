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
	PairApprovePairing struct {
		// RequestID GET parameter
		//
		// Pair requestID
		RequestID uint64 `json:",string"`
	}

	PairCompletePairing struct {
		// Token POST parameter
		//
		// Auth token of the origin node
		Token string
	}
)

// NewPairApprovePairing request
func NewPairApprovePairing() *PairApprovePairing {
	return &PairApprovePairing{}
}

// Auditable returns all auditable/loggable parameters
func (r PairApprovePairing) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PairApprovePairing) GetRequestID() uint64 {
	return r.RequestID
}

// Fill processes request and fills internal variables
func (r *PairApprovePairing) Fill(req *http.Request) (err error) {
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

		if val, ok := tmp["requestID"]; ok && len(val) > 0 {
			r.RequestID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewPairCompletePairing request
func NewPairCompletePairing() *PairCompletePairing {
	return &PairCompletePairing{}
}

// Auditable returns all auditable/loggable parameters
func (r PairCompletePairing) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"token": r.Token,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PairCompletePairing) GetToken() string {
	return r.Token
}

// Fill processes request and fills internal variables
func (r *PairCompletePairing) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["token"]; ok && len(val) > 0 {
			r.Token, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
