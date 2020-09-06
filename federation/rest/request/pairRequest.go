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
	PairRequestRequestPairing struct {
		// Identifier POST parameter
		//
		// Origin node identifier
		Identifier string

		// Token POST parameter
		//
		// Destination node token
		Token string
	}
)

// NewPairRequestRequestPairing request
func NewPairRequestRequestPairing() *PairRequestRequestPairing {
	return &PairRequestRequestPairing{}
}

// Auditable returns all auditable/loggable parameters
func (r PairRequestRequestPairing) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"identifier": r.Identifier,
		"token":      r.Token,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PairRequestRequestPairing) GetIdentifier() string {
	return r.Identifier
}

// Auditable returns all auditable/loggable parameters
func (r PairRequestRequestPairing) GetToken() string {
	return r.Token
}

// Fill processes request and fills internal variables
func (r *PairRequestRequestPairing) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["token"]; ok && len(val) > 0 {
			r.Token, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
