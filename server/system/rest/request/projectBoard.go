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
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/go-chi/chi/v5"
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
	ProjectBoardBoard struct {
		// ProjectID GET parameter
		//
		// Chain-root project scope
		ProjectID uint64 `json:",string"`

		// RevisionID GET parameter
		//
		// Scope to a single revision (default is chain-wide, including unassigned items)
		RevisionID uint64 `json:",string"`

		// Status GET parameter
		//
		// Return only this column (required when paging with a cursor)
		Status string

		// PageCursor GET parameter
		//
		// Continue a single column's paging from its previous page
		PageCursor string

		// Limit GET parameter
		//
		// Cards per column page (default 20)
		Limit uint
	}
)

// NewProjectBoardBoard request
func NewProjectBoardBoard() *ProjectBoardBoard {
	return &ProjectBoardBoard{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBoardBoard) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":  r.ProjectID,
		"revisionID": r.RevisionID,
		"status":     r.Status,
		"pageCursor": r.PageCursor,
		"limit":      r.Limit,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBoardBoard) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBoardBoard) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBoardBoard) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBoardBoard) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBoardBoard) GetLimit() uint {
	return r.Limit
}

// Fill processes request and fills internal variables
func (r *ProjectBoardBoard) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["projectID"]; ok && len(val) > 0 {
			r.ProjectID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["revisionID"]; ok && len(val) > 0 {
			r.RevisionID, err = payload.ParseUint64(val[0]), nil
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
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
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
	}

	return err
}
