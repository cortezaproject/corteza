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
	DataPrivacyRecordList struct {
		// SensitivityLevelID GET parameter
		//
		// Sensitivity Level ID
		SensitivityLevelID uint64 `json:",string"`

		// ConnectionID GET parameter
		//
		// Filter by connection ID
		ConnectionID []string
	}

	DataPrivacyModuleList struct {
		// ConnectionID GET parameter
		//
		// Filter by connection ID
		ConnectionID []string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}
)

// NewDataPrivacyRecordList request
func NewDataPrivacyRecordList() *DataPrivacyRecordList {
	return &DataPrivacyRecordList{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRecordList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sensitivityLevelID": r.SensitivityLevelID,
		"connectionID":       r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRecordList) GetSensitivityLevelID() uint64 {
	return r.SensitivityLevelID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRecordList) GetConnectionID() []string {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRecordList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["sensitivityLevelID"]; ok && len(val) > 0 {
			r.SensitivityLevelID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["connectionID[]"]; ok {
			r.ConnectionID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["connectionID"]; ok {
			r.ConnectionID, err = val, nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDataPrivacyModuleList request
func NewDataPrivacyModuleList() *DataPrivacyModuleList {
	return &DataPrivacyModuleList{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyModuleList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"limit":        r.Limit,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyModuleList) GetConnectionID() []string {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyModuleList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyModuleList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyModuleList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *DataPrivacyModuleList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["connectionID[]"]; ok {
			r.ConnectionID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["connectionID"]; ok {
			r.ConnectionID, err = val, nil
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
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["sort"]; ok && len(val) > 0 {
			r.Sort, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
