package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	DmlConnection struct {
		ID        uint64              `json:"connectionID,string"`
		Handle    string              `json:"handle"`
		Label     string              `json:"label"`
		Params    DmlConnectionParams `json:"params"`
		CreatedAt time.Time           `json:"createdAt,omitempty"`
		UpdatedAt *time.Time          `json:"updatedAt,omitempty"`
		DeletedAt *time.Time          `json:"deletedAt,omitempty"`
	}

	DmlConnectionParams struct {
		Type            string         `json:"type"`
		Params          map[string]any `json:"params"`
		ModelIdent      string         `json:"modelIdent"`
		ModelIdentCheck []string       `json:"modelIdentCheck"`
	}
)

func (r DmlConnection) Clone() *DmlConnection {
	dup := r
	dup.Params = *r.Params.Clone()

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	return &dup
}

func (r DmlConnection) Diff(cmp *DmlConnection) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DmlConnection{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "connectionID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	for _, c := range r.Params.Diff(&cmp.Params) {
		c.Key = "params." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	return out
}

func (r DmlConnectionParams) Clone() *DmlConnectionParams {
	dup := r
	if r.Params != nil {
		dup.Params = make(map[string]any, len(r.Params))
		for k, v := range r.Params {
			dup.Params[k] = v
		}
	}

	if r.ModelIdentCheck != nil {
		dup.ModelIdentCheck = make([]string, len(r.ModelIdentCheck))
		copy(dup.ModelIdentCheck, r.ModelIdentCheck)
	}

	return &dup
}

func (r DmlConnectionParams) Diff(cmp *DmlConnectionParams) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DmlConnectionParams{}
	}
	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if !reflect.DeepEqual(r.Params, cmp.Params) {
		out = append(out, &revisions.Change{Key: "params", Old: []any{cmp.Params}, New: []any{r.Params}})
	}

	if r.ModelIdent != cmp.ModelIdent {
		out = append(out, &revisions.Change{Key: "modelIdent", Old: []any{cmp.ModelIdent}, New: []any{r.ModelIdent}})
	}

	if !reflect.DeepEqual(r.ModelIdentCheck, cmp.ModelIdentCheck) {
		out = append(out, &revisions.Change{Key: "modelIdentCheck", Old: []any{cmp.ModelIdentCheck}, New: []any{r.ModelIdentCheck}})
	}

	return out
}

func (r *DmlConnectionParams) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DmlConnectionParams) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseDmlConnectionParams(ss []string) (p DmlConnectionParams, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
