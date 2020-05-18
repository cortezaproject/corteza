package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `actionlog.go`, `actionlog.util.go` or `actionlog_test.go` to
	implement your API calls, helper functions and tests. The file `actionlog.go`
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

	"time"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// ActionlogList request parameters
type ActionlogList struct {
	hasFrom bool
	rawFrom string
	From    *time.Time

	hasTo bool
	rawTo string
	To    *time.Time

	hasResource bool
	rawResource string
	Resource    string

	hasAction bool
	rawAction string
	Action    string

	hasActorID bool
	rawActorID []string
	ActorID    []string

	hasLimit bool
	rawLimit string
	Limit    uint

	hasOffset bool
	rawOffset string
	Offset    uint

	hasPage bool
	rawPage string
	Page    uint

	hasPerPage bool
	rawPerPage string
	PerPage    uint
}

// NewActionlogList request
func NewActionlogList() *ActionlogList {
	return &ActionlogList{}
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["from"] = r.From
	out["to"] = r.To
	out["resource"] = r.Resource
	out["action"] = r.Action
	out["actorID"] = r.ActorID
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage

	return out
}

// Fill processes request and fills internal variables
func (r *ActionlogList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["from"]; ok {
		r.hasFrom = true
		r.rawFrom = val

		if r.From, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := get["to"]; ok {
		r.hasTo = true
		r.rawTo = val

		if r.To, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := get["resource"]; ok {
		r.hasResource = true
		r.rawResource = val
		r.Resource = val
	}
	if val, ok := get["action"]; ok {
		r.hasAction = true
		r.rawAction = val
		r.Action = val
	}

	if val, ok := urlQuery["actorID[]"]; ok {
		r.hasActorID = true
		r.rawActorID = val
		r.ActorID = parseStrings(val)
	} else if val, ok = urlQuery["actorID"]; ok {
		r.hasActorID = true
		r.rawActorID = val
		r.ActorID = parseStrings(val)
	}

	if val, ok := get["limit"]; ok {
		r.hasLimit = true
		r.rawLimit = val
		r.Limit = parseUint(val)
	}
	if val, ok := get["offset"]; ok {
		r.hasOffset = true
		r.rawOffset = val
		r.Offset = parseUint(val)
	}
	if val, ok := get["page"]; ok {
		r.hasPage = true
		r.rawPage = val
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.hasPerPage = true
		r.rawPerPage = val
		r.PerPage = parseUint(val)
	}

	return err
}

var _ RequestFiller = NewActionlogList()

// HasFrom returns true if from was set
func (r *ActionlogList) HasFrom() bool {
	return r.hasFrom
}

// RawFrom returns raw value of from parameter
func (r *ActionlogList) RawFrom() string {
	return r.rawFrom
}

// GetFrom returns casted value of  from parameter
func (r *ActionlogList) GetFrom() *time.Time {
	return r.From
}

// HasTo returns true if to was set
func (r *ActionlogList) HasTo() bool {
	return r.hasTo
}

// RawTo returns raw value of to parameter
func (r *ActionlogList) RawTo() string {
	return r.rawTo
}

// GetTo returns casted value of  to parameter
func (r *ActionlogList) GetTo() *time.Time {
	return r.To
}

// HasResource returns true if resource was set
func (r *ActionlogList) HasResource() bool {
	return r.hasResource
}

// RawResource returns raw value of resource parameter
func (r *ActionlogList) RawResource() string {
	return r.rawResource
}

// GetResource returns casted value of  resource parameter
func (r *ActionlogList) GetResource() string {
	return r.Resource
}

// HasAction returns true if action was set
func (r *ActionlogList) HasAction() bool {
	return r.hasAction
}

// RawAction returns raw value of action parameter
func (r *ActionlogList) RawAction() string {
	return r.rawAction
}

// GetAction returns casted value of  action parameter
func (r *ActionlogList) GetAction() string {
	return r.Action
}

// HasActorID returns true if actorID was set
func (r *ActionlogList) HasActorID() bool {
	return r.hasActorID
}

// RawActorID returns raw value of actorID parameter
func (r *ActionlogList) RawActorID() []string {
	return r.rawActorID
}

// GetActorID returns casted value of  actorID parameter
func (r *ActionlogList) GetActorID() []string {
	return r.ActorID
}

// HasLimit returns true if limit was set
func (r *ActionlogList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *ActionlogList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *ActionlogList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *ActionlogList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *ActionlogList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *ActionlogList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *ActionlogList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *ActionlogList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *ActionlogList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *ActionlogList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *ActionlogList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *ActionlogList) GetPerPage() uint {
	return r.PerPage
}
