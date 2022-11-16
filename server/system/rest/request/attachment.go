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
	AttachmentRead struct {
		// Kind PATH parameter
		//
		// Kind
		Kind string

		// AttachmentID PATH parameter
		//
		// Attachment ID
		AttachmentID uint64 `json:",string"`

		// Sign GET parameter
		//
		// Signature
		Sign string

		// UserID GET parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	AttachmentDelete struct {
		// Kind PATH parameter
		//
		// Kind
		Kind string

		// AttachmentID PATH parameter
		//
		// Attachment ID
		AttachmentID uint64 `json:",string"`

		// Sign GET parameter
		//
		// Signature
		Sign string

		// UserID GET parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	AttachmentOriginal struct {
		// Kind PATH parameter
		//
		// Kind
		Kind string

		// AttachmentID PATH parameter
		//
		// Attachment ID
		AttachmentID uint64 `json:",string"`

		// Name PATH parameter
		//
		// File name
		Name string

		// Sign GET parameter
		//
		// Signature
		Sign string

		// UserID GET parameter
		//
		// User ID
		UserID uint64 `json:",string"`

		// Download GET parameter
		//
		// Force file download
		Download bool
	}

	AttachmentPreview struct {
		// Kind PATH parameter
		//
		// Kind
		Kind string

		// AttachmentID PATH parameter
		//
		// Attachment ID
		AttachmentID uint64 `json:",string"`

		// Ext PATH parameter
		//
		// Preview extension/format
		Ext string

		// Sign GET parameter
		//
		// Signature
		Sign string

		// UserID GET parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}
)

// NewAttachmentRead request
func NewAttachmentRead() *AttachmentRead {
	return &AttachmentRead{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"kind":         r.Kind,
		"attachmentID": r.AttachmentID,
		"sign":         r.Sign,
		"userID":       r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentRead) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentRead) GetAttachmentID() uint64 {
	return r.AttachmentID
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentRead) GetSign() string {
	return r.Sign
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentRead) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *AttachmentRead) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["sign"]; ok && len(val) > 0 {
			r.Sign, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "kind")
		r.Kind, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "attachmentID")
		r.AttachmentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAttachmentDelete request
func NewAttachmentDelete() *AttachmentDelete {
	return &AttachmentDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"kind":         r.Kind,
		"attachmentID": r.AttachmentID,
		"sign":         r.Sign,
		"userID":       r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentDelete) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentDelete) GetAttachmentID() uint64 {
	return r.AttachmentID
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentDelete) GetSign() string {
	return r.Sign
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentDelete) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *AttachmentDelete) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["sign"]; ok && len(val) > 0 {
			r.Sign, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "kind")
		r.Kind, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "attachmentID")
		r.AttachmentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAttachmentOriginal request
func NewAttachmentOriginal() *AttachmentOriginal {
	return &AttachmentOriginal{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentOriginal) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"kind":         r.Kind,
		"attachmentID": r.AttachmentID,
		"name":         r.Name,
		"sign":         r.Sign,
		"userID":       r.UserID,
		"download":     r.Download,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentOriginal) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentOriginal) GetAttachmentID() uint64 {
	return r.AttachmentID
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentOriginal) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentOriginal) GetSign() string {
	return r.Sign
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentOriginal) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentOriginal) GetDownload() bool {
	return r.Download
}

// Fill processes request and fills internal variables
func (r *AttachmentOriginal) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["sign"]; ok && len(val) > 0 {
			r.Sign, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["download"]; ok && len(val) > 0 {
			r.Download, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "kind")
		r.Kind, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "attachmentID")
		r.AttachmentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "name")
		r.Name, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAttachmentPreview request
func NewAttachmentPreview() *AttachmentPreview {
	return &AttachmentPreview{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentPreview) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"kind":         r.Kind,
		"attachmentID": r.AttachmentID,
		"ext":          r.Ext,
		"sign":         r.Sign,
		"userID":       r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentPreview) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentPreview) GetAttachmentID() uint64 {
	return r.AttachmentID
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentPreview) GetExt() string {
	return r.Ext
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentPreview) GetSign() string {
	return r.Sign
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentPreview) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *AttachmentPreview) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["sign"]; ok && len(val) > 0 {
			r.Sign, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "kind")
		r.Kind, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "attachmentID")
		r.AttachmentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "ext")
		r.Ext, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
