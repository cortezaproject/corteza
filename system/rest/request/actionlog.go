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
	"time"
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
	ActionlogList struct {
		// From GET parameter
		//
		// From
		From *time.Time

		// To GET parameter
		//
		// To
		To *time.Time

		// Resource GET parameter
		//
		// Resource
		Resource string

		// Action GET parameter
		//
		// Action
		Action string

		// ActorID GET parameter
		//
		// Filter by one or more actors
		ActorID []string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// Offset GET parameter
		//
		// Offset
		Offset uint

		// Page GET parameter
		//
		// Page number (1-based)
		Page uint

		// PerPage GET parameter
		//
		// Returned items per page (default 50)
		PerPage uint
	}
)

// NewActionlogList request
func NewActionlogList() *ActionlogList {
	return &ActionlogList{}
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"from":     r.From,
		"to":       r.To,
		"resource": r.Resource,
		"action":   r.Action,
		"actorID":  r.ActorID,
		"limit":    r.Limit,
		"offset":   r.Offset,
		"page":     r.Page,
		"perPage":  r.PerPage,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetFrom() *time.Time {
	return r.From
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetTo() *time.Time {
	return r.To
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetAction() string {
	return r.Action
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetActorID() []string {
	return r.ActorID
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetOffset() uint {
	return r.Offset
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetPage() uint {
	return r.Page
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetPerPage() uint {
	return r.PerPage
}

// Fill processes request and fills internal variables
func (r *ActionlogList) Fill(req *http.Request) (err error) {
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

		if val, ok := tmp["from"]; ok && len(val) > 0 {
			r.From, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["to"]; ok && len(val) > 0 {
			r.To, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["action"]; ok && len(val) > 0 {
			r.Action, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["actorID[]"]; ok {
			r.ActorID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["actorID"]; ok {
			r.ActorID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["offset"]; ok && len(val) > 0 {
			r.Offset, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["page"]; ok && len(val) > 0 {
			r.Page, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["perPage"]; ok && len(val) > 0 {
			r.PerPage, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
