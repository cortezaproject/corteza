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
	NodeHandshakeInitialize struct {
		// NodeID PATH parameter
		//
		// NodeID
		NodeID uint64 `json:",string"`

		// NodeURI POST parameter
		//
		// Node A node URI
		NodeURI string

		// TokenB POST parameter
		//
		// Node B auth token
		TokenB string

		// NodeIDB POST parameter
		//
		// Node B nodeID
		NodeIDB uint64 `json:",string"`
	}
)

// NewNodeHandshakeInitialize request
func NewNodeHandshakeInitialize() *NodeHandshakeInitialize {
	return &NodeHandshakeInitialize{}
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":  r.NodeID,
		"nodeURI": r.NodeURI,
		"tokenB":  r.TokenB,
		"nodeIDB": r.NodeIDB,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) GetNodeURI() string {
	return r.NodeURI
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) GetTokenB() string {
	return r.TokenB
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeInitialize) GetNodeIDB() uint64 {
	return r.NodeIDB
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
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["nodeURI"]; ok && len(val) > 0 {
			r.NodeURI, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["tokenB"]; ok && len(val) > 0 {
			r.TokenB, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["nodeIDB"]; ok && len(val) > 0 {
			r.NodeIDB, err = payload.ParseUint64(val[0]), nil
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
