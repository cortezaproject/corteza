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
	NodeHandshakeInitialize struct {
		// NodeID PATH parameter
		//
		// NodeID
		NodeID uint64 `json:",string"`

		// PairToken POST parameter
		//
		// Pairing token to authenticate handshake initialization request
		PairToken string

		// SharedNodeID POST parameter
		//
		// Remote (invoker's) node ID
		SharedNodeID uint64 `json:",string"`

		// AuthToken POST parameter
		//
		// Authentication token so that remote
		AuthToken string
	}
)

// NewNodeHandshakeInitialize request
func NewNodeHandshakeInitialize() *NodeHandshakeInitialize {
	return &NodeHandshakeInitialize{}
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":       r.NodeID,
		"pairToken":    r.PairToken,
		"sharedNodeID": r.SharedNodeID,
		"authToken":    r.AuthToken,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) GetPairToken() string {
	return r.PairToken
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) GetSharedNodeID() uint64 {
	return r.SharedNodeID
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) GetAuthToken() string {
	return r.AuthToken
}

// Fill processes request and fills internal variables
func (r *NodeHandshakeInitialize) Fill(req *http.Request) (err error) {

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
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["pairToken"]; ok && len(val) > 0 {
				r.PairToken, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["sharedNodeID"]; ok && len(val) > 0 {
				r.SharedNodeID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["authToken"]; ok && len(val) > 0 {
				r.AuthToken, err = val[0], nil
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

		if val, ok := req.Form["pairToken"]; ok && len(val) > 0 {
			r.PairToken, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["sharedNodeID"]; ok && len(val) > 0 {
			r.SharedNodeID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["authToken"]; ok && len(val) > 0 {
			r.AuthToken, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
