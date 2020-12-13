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
	NodeSearch struct {
		// Query GET parameter
		//
		// Filter nodes by name and host
		Query string

		// Status GET parameter
		//
		// Filter by status
		Status string
	}

	NodeCreate struct {
		// Host POST parameter
		//
		// Node B host
		Host string

		// BaseURL POST parameter
		//
		// Federation API base URL
		BaseURL string

		// Name POST parameter
		//
		// Name for this node
		Name string

		// PairingURI POST parameter
		//
		// Pairing URI
		PairingURI string
	}

	NodeGenerateURI struct {
		// NodeID PATH parameter
		//
		// NodeID
		NodeID uint64 `json:",string"`
	}

	NodePair struct {
		// NodeID PATH parameter
		//
		// NodeID
		NodeID uint64 `json:",string"`
	}

	NodeHandshakeConfirm struct {
		// NodeID PATH parameter
		//
		// NodeID
		NodeID uint64 `json:",string"`
	}

	NodeHandshakeComplete struct {
		// NodeID PATH parameter
		//
		// NodeID
		NodeID uint64 `json:",string"`

		// TokenA POST parameter
		//
		// Node A token
		TokenA string
	}
)

// NewNodeSearch request
func NewNodeSearch() *NodeSearch {
	return &NodeSearch{}
}

// Auditable returns all auditable/loggable parameters
func (r NodeSearch) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":  r.Query,
		"status": r.Status,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodeSearch) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r NodeSearch) GetStatus() string {
	return r.Status
}

// Fill processes request and fills internal variables
func (r *NodeSearch) Fill(req *http.Request) (err error) {
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
		if val, ok := tmp["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewNodeCreate request
func NewNodeCreate() *NodeCreate {
	return &NodeCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"host":       r.Host,
		"baseURL":    r.BaseURL,
		"name":       r.Name,
		"pairingURI": r.PairingURI,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetHost() string {
	return r.Host
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetBaseURL() string {
	return r.BaseURL
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetPairingURI() string {
	return r.PairingURI
}

// Fill processes request and fills internal variables
func (r *NodeCreate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["host"]; ok && len(val) > 0 {
			r.Host, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["baseURL"]; ok && len(val) > 0 {
			r.BaseURL, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["pairingURI"]; ok && len(val) > 0 {
			r.PairingURI, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewNodeGenerateURI request
func NewNodeGenerateURI() *NodeGenerateURI {
	return &NodeGenerateURI{}
}

// Auditable returns all auditable/loggable parameters
func (r NodeGenerateURI) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID": r.NodeID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodeGenerateURI) GetNodeID() uint64 {
	return r.NodeID
}

// Fill processes request and fills internal variables
func (r *NodeGenerateURI) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNodePair request
func NewNodePair() *NodePair {
	return &NodePair{}
}

// Auditable returns all auditable/loggable parameters
func (r NodePair) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID": r.NodeID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodePair) GetNodeID() uint64 {
	return r.NodeID
}

// Fill processes request and fills internal variables
func (r *NodePair) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNodeHandshakeConfirm request
func NewNodeHandshakeConfirm() *NodeHandshakeConfirm {
	return &NodeHandshakeConfirm{}
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeConfirm) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID": r.NodeID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeConfirm) GetNodeID() uint64 {
	return r.NodeID
}

// Fill processes request and fills internal variables
func (r *NodeHandshakeConfirm) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNodeHandshakeComplete request
func NewNodeHandshakeComplete() *NodeHandshakeComplete {
	return &NodeHandshakeComplete{}
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeComplete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID": r.NodeID,
		"tokenA": r.TokenA,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeComplete) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r NodeHandshakeComplete) GetTokenA() string {
	return r.TokenA
}

// Fill processes request and fills internal variables
func (r *NodeHandshakeComplete) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["tokenA"]; ok && len(val) > 0 {
			r.TokenA, err = val[0], nil
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
