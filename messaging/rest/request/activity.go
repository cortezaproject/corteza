package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `activity.go`, `activity.util.go` or `activity_test.go` to
	implement your API calls, helper functions and tests. The file `activity.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"io"
	"strings"

	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// ActivitySend request parameters
type ActivitySend struct {
	hasChannelID bool
	rawChannelID string
	ChannelID    uint64 `json:",string"`

	hasMessageID bool
	rawMessageID string
	MessageID    uint64 `json:",string"`

	hasKind bool
	rawKind string
	Kind    string
}

// NewActivitySend request
func NewActivitySend() *ActivitySend {
	return &ActivitySend{}
}

// Auditable returns all auditable/loggable parameters
func (r ActivitySend) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["messageID"] = r.MessageID
	out["kind"] = r.Kind

	return out
}

// Fill processes request and fills internal variables
func (r *ActivitySend) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := post["channelID"]; ok {
		r.hasChannelID = true
		r.rawChannelID = val
		r.ChannelID = parseUInt64(val)
	}
	if val, ok := post["messageID"]; ok {
		r.hasMessageID = true
		r.rawMessageID = val
		r.MessageID = parseUInt64(val)
	}
	if val, ok := post["kind"]; ok {
		r.hasKind = true
		r.rawKind = val
		r.Kind = val
	}

	return err
}

var _ RequestFiller = NewActivitySend()

// HasChannelID returns true if channelID was set
func (r *ActivitySend) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *ActivitySend) RawChannelID() string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *ActivitySend) GetChannelID() uint64 {
	return r.ChannelID
}

// HasMessageID returns true if messageID was set
func (r *ActivitySend) HasMessageID() bool {
	return r.hasMessageID
}

// RawMessageID returns raw value of messageID parameter
func (r *ActivitySend) RawMessageID() string {
	return r.rawMessageID
}

// GetMessageID returns casted value of  messageID parameter
func (r *ActivitySend) GetMessageID() uint64 {
	return r.MessageID
}

// HasKind returns true if kind was set
func (r *ActivitySend) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *ActivitySend) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *ActivitySend) GetKind() string {
	return r.Kind
}
