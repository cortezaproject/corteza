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
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/go-chi/chi/v5"
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
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	FeedChanges struct {
		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// From GET parameter
		//
		// From timestamp
		From *time.Time

		// To GET parameter
		//
		// To timestamp
		To *time.Time
	}
)

// NewFeedChanges request
func NewFeedChanges() *FeedChanges {
	return &FeedChanges{}
}

// Auditable returns all auditable/loggable parameters
func (r FeedChanges) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"from":       r.From,
		"to":         r.To,
	}
}

// Auditable returns all auditable/loggable parameters
func (r FeedChanges) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r FeedChanges) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r FeedChanges) GetFrom() *time.Time {
	return r.From
}

// Auditable returns all auditable/loggable parameters
func (r FeedChanges) GetTo() *time.Time {
	return r.To
}

// Fill processes request and fills internal variables
func (r *FeedChanges) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
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
	}

	return err
}
