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
	DataPrivacyListSensitiveData struct {
		// SensitivityLevelID GET parameter
		//
		// Sensitivity Level ID
		SensitivityLevelID uint64 `json:",string"`
	}
)

// NewDataPrivacyListSensitiveData request
func NewDataPrivacyListSensitiveData() *DataPrivacyListSensitiveData {
	return &DataPrivacyListSensitiveData{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListSensitiveData) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sensitivityLevelID": r.SensitivityLevelID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListSensitiveData) GetSensitivityLevelID() uint64 {
	return r.SensitivityLevelID
}

// Fill processes request and fills internal variables
func (r *DataPrivacyListSensitiveData) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["sensitivityLevelID"]; ok && len(val) > 0 {
			r.SensitivityLevelID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
