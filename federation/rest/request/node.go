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
	NodeCreate struct {
		// MyDomain POST parameter
		//
		// [TMP] field that determines my domain
		MyDomain string

		// Domain POST parameter
		//
		// Node B domain
		Domain string

		// Name POST parameter
		//
		// Name for this node
		Name string

		// AdminContact POST parameter
		//
		// Node B admin contact email
		AdminContact string

		// NodeURI POST parameter
		//
		// Node A URI
		NodeURI string
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

// NewNodeCreate request
func NewNodeCreate() *NodeCreate {
	return &NodeCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"myDomain":     r.MyDomain,
		"domain":       r.Domain,
		"name":         r.Name,
		"adminContact": r.AdminContact,
		"nodeURI":      r.NodeURI,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetMyDomain() string {
	return r.MyDomain
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetDomain() string {
	return r.Domain
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetAdminContact() string {
	return r.AdminContact
}

// Auditable returns all auditable/loggable parameters
func (r NodeCreate) GetNodeURI() string {
	return r.NodeURI
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

		if val, ok := req.Form["myDomain"]; ok && len(val) > 0 {
			r.MyDomain, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["domain"]; ok && len(val) > 0 {
			r.Domain, err = val[0], nil
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

		if val, ok := req.Form["adminContact"]; ok && len(val) > 0 {
			r.AdminContact, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["nodeURI"]; ok && len(val) > 0 {
			r.NodeURI, err = val[0], nil
			if err != nil {
				return err
			}
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
