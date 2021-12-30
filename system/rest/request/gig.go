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
	"github.com/cortezaproject/corteza-server/pkg/gig"
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
	GigCreate struct {
		// Worker POST parameter
		//
		// Gig worker
		Worker string

		// Preprocessors POST parameter
		//
		// Worker preprocessing to do
		Preprocessors gig.PreprocessorWrapSet

		// Postprocessors POST parameter
		//
		// Output postprocessing to do
		Postprocessors gig.PostprocessorWrapSet
	}

	GigRead struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigUpdate struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`

		// Decoders POST parameter
		//
		// Decoders to apply to the sources
		Decoders gig.DecoderWrapSet

		// Preprocessors POST parameter
		//
		// Worker preprocessing to do
		Preprocessors gig.PreprocessorWrapSet

		// Postprocessors POST parameter
		//
		// Output postprocessing to do
		Postprocessors gig.PostprocessorWrapSet
	}

	GigAddSource struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`

		// Upload POST parameter
		//
		// File source to add
		Upload *multipart.FileHeader

		// Uri POST parameter
		//
		// Source location to add
		Uri string

		// Decoders POST parameter
		//
		// Decoders to apply to the sources
		Decoders gig.DecoderWrapSet
	}

	GigOutput struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigPrepare struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigExec struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigStatus struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigComplete struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigTasks struct {
	}
)

// NewGigCreate request
func NewGigCreate() *GigCreate {
	return &GigCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"worker":         r.Worker,
		"preprocessors":  r.Preprocessors,
		"postprocessors": r.Postprocessors,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) GetWorker() string {
	return r.Worker
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) GetPreprocessors() gig.PreprocessorWrapSet {
	return r.Preprocessors
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) GetPostprocessors() gig.PostprocessorWrapSet {
	return r.Postprocessors
}

// Fill processes request and fills internal variables
func (r *GigCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["worker"]; ok && len(val) > 0 {
				r.Worker, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["preprocessors[]"]; ok {
				r.Preprocessors, err = gig.ParsePreprocessorWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["preprocessors"]; ok {
				r.Preprocessors, err = gig.ParsePreprocessorWrap(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["postprocessors[]"]; ok {
				r.Postprocessors, err = gig.ParsePostprocessorWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["postprocessors"]; ok {
				r.Postprocessors, err = gig.ParsePostprocessorWrap(val)
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

		if val, ok := req.Form["worker"]; ok && len(val) > 0 {
			r.Worker, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["preprocessors[]"]; ok {
			r.Preprocessors, err = gig.ParsePreprocessorWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["preprocessors"]; ok {
			r.Preprocessors, err = gig.ParsePreprocessorWrap(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["postprocessors[]"]; ok {
			r.Postprocessors, err = gig.ParsePostprocessorWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["postprocessors"]; ok {
			r.Postprocessors, err = gig.ParsePostprocessorWrap(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewGigRead request
func NewGigRead() *GigRead {
	return &GigRead{}
}

// Auditable returns all auditable/loggable parameters
func (r GigRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigRead) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigUpdate request
func NewGigUpdate() *GigUpdate {
	return &GigUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID":          r.GigID,
		"decoders":       r.Decoders,
		"preprocessors":  r.Preprocessors,
		"postprocessors": r.Postprocessors,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetGigID() uint64 {
	return r.GigID
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetDecoders() gig.DecoderWrapSet {
	return r.Decoders
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetPreprocessors() gig.PreprocessorWrapSet {
	return r.Preprocessors
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetPostprocessors() gig.PostprocessorWrapSet {
	return r.Postprocessors
}

// Fill processes request and fills internal variables
func (r *GigUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["preprocessors[]"]; ok {
				r.Preprocessors, err = gig.ParsePreprocessorWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["preprocessors"]; ok {
				r.Preprocessors, err = gig.ParsePreprocessorWrap(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["postprocessors[]"]; ok {
				r.Postprocessors, err = gig.ParsePostprocessorWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["postprocessors"]; ok {
				r.Postprocessors, err = gig.ParsePostprocessorWrap(val)
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

		//if val, ok := req.Form["decoders[]"]; ok && len(val) > 0  {
		//    r.Decoders, err = gig.DecoderWrapSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["preprocessors[]"]; ok {
			r.Preprocessors, err = gig.ParsePreprocessorWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["preprocessors"]; ok {
			r.Preprocessors, err = gig.ParsePreprocessorWrap(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["postprocessors[]"]; ok {
			r.Postprocessors, err = gig.ParsePostprocessorWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["postprocessors"]; ok {
			r.Postprocessors, err = gig.ParsePostprocessorWrap(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigAddSource request
func NewGigAddSource() *GigAddSource {
	return &GigAddSource{}
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID":    r.GigID,
		"upload":   r.Upload,
		"uri":      r.Uri,
		"decoders": r.Decoders,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) GetGigID() uint64 {
	return r.GigID
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) GetUri() string {
	return r.Uri
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) GetDecoders() gig.DecoderWrapSet {
	return r.Decoders
}

// Fill processes request and fills internal variables
func (r *GigAddSource) Fill(req *http.Request) (err error) {

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

			// Ignoring upload as its handled in the POST params section

			if val, ok := req.MultipartForm.Value["uri"]; ok && len(val) > 0 {
				r.Uri, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["decoders[]"]; ok {
				r.Decoders, err = gig.ParseDecoderWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["decoders"]; ok {
				r.Decoders, err = gig.ParseDecoderWrap(val)
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

		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

		if val, ok := req.Form["uri"]; ok && len(val) > 0 {
			r.Uri, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["decoders[]"]; ok {
			r.Decoders, err = gig.ParseDecoderWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["decoders"]; ok {
			r.Decoders, err = gig.ParseDecoderWrap(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigOutput request
func NewGigOutput() *GigOutput {
	return &GigOutput{}
}

// Auditable returns all auditable/loggable parameters
func (r GigOutput) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigOutput) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigOutput) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigPrepare request
func NewGigPrepare() *GigPrepare {
	return &GigPrepare{}
}

// Auditable returns all auditable/loggable parameters
func (r GigPrepare) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigPrepare) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigPrepare) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigExec request
func NewGigExec() *GigExec {
	return &GigExec{}
}

// Auditable returns all auditable/loggable parameters
func (r GigExec) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigExec) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigExec) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigStatus request
func NewGigStatus() *GigStatus {
	return &GigStatus{}
}

// Auditable returns all auditable/loggable parameters
func (r GigStatus) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigStatus) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigStatus) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigComplete request
func NewGigComplete() *GigComplete {
	return &GigComplete{}
}

// Auditable returns all auditable/loggable parameters
func (r GigComplete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigComplete) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigComplete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigTasks request
func NewGigTasks() *GigTasks {
	return &GigTasks{}
}

// Auditable returns all auditable/loggable parameters
func (r GigTasks) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *GigTasks) Fill(req *http.Request) (err error) {

	return err
}
