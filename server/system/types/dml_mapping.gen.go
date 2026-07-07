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
	DmlMapping struct {
		ID              uint64          `json:"mappingID,string"`
		ConnectionID    uint64          `json:"connectionID,string"`
		NamespaceHandle string          `json:"namespaceHandle"`
		SourceIdent     string          `json:"sourceIdent"`
		ModuleHandle    string          `json:"moduleHandle"`
		ModuleName      string          `json:"moduleName"`
		Skip            bool            `json:"skip"`
		Identifier      string          `json:"identifier"`
		Columns         DmlColumnMapSet `json:"columns"`
		CreatedAt       time.Time       `json:"createdAt,omitempty"`
		UpdatedAt       *time.Time      `json:"updatedAt,omitempty"`
		DeletedAt       *time.Time      `json:"deletedAt,omitempty"`
	}

	DmlColumnMap struct {
		SourceIdent string `json:"sourceIdent"`
		FieldName   string `json:"fieldName"`
		Label       string `json:"label,omitempty"`
		FieldKind   string `json:"fieldKind"`
		Skip        bool   `json:"skip"`
	}
)

func (r DmlMapping) Clone() *DmlMapping {
	dup := r
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

func (r DmlMapping) Diff(cmp *DmlMapping) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DmlMapping{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "mappingID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.ConnectionID != cmp.ConnectionID {
		out = append(out, &revisions.Change{Key: "connectionID", Old: []any{cmp.ConnectionID}, New: []any{r.ConnectionID}})
	}

	if r.NamespaceHandle != cmp.NamespaceHandle {
		out = append(out, &revisions.Change{Key: "namespaceHandle", Old: []any{cmp.NamespaceHandle}, New: []any{r.NamespaceHandle}})
	}

	if r.SourceIdent != cmp.SourceIdent {
		out = append(out, &revisions.Change{Key: "sourceIdent", Old: []any{cmp.SourceIdent}, New: []any{r.SourceIdent}})
	}

	if r.ModuleHandle != cmp.ModuleHandle {
		out = append(out, &revisions.Change{Key: "moduleHandle", Old: []any{cmp.ModuleHandle}, New: []any{r.ModuleHandle}})
	}

	if r.ModuleName != cmp.ModuleName {
		out = append(out, &revisions.Change{Key: "moduleName", Old: []any{cmp.ModuleName}, New: []any{r.ModuleName}})
	}

	if r.Skip != cmp.Skip {
		out = append(out, &revisions.Change{Key: "skip", Old: []any{cmp.Skip}, New: []any{r.Skip}})
	}

	if r.Identifier != cmp.Identifier {
		out = append(out, &revisions.Change{Key: "identifier", Old: []any{cmp.Identifier}, New: []any{r.Identifier}})
	}

	if !reflect.DeepEqual(r.Columns, cmp.Columns) {
		out = append(out, &revisions.Change{Key: "columns", Old: []any{cmp.Columns}, New: []any{r.Columns}})
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

func (r DmlColumnMap) Clone() *DmlColumnMap {
	dup := r
	return &dup
}

func (r DmlColumnMap) Diff(cmp *DmlColumnMap) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DmlColumnMap{}
	}
	if r.SourceIdent != cmp.SourceIdent {
		out = append(out, &revisions.Change{Key: "sourceIdent", Old: []any{cmp.SourceIdent}, New: []any{r.SourceIdent}})
	}

	if r.FieldName != cmp.FieldName {
		out = append(out, &revisions.Change{Key: "fieldName", Old: []any{cmp.FieldName}, New: []any{r.FieldName}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	if r.FieldKind != cmp.FieldKind {
		out = append(out, &revisions.Change{Key: "fieldKind", Old: []any{cmp.FieldKind}, New: []any{r.FieldKind}})
	}

	if r.Skip != cmp.Skip {
		out = append(out, &revisions.Change{Key: "skip", Old: []any{cmp.Skip}, New: []any{r.Skip}})
	}

	return out
}

func (r *DmlColumnMap) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DmlColumnMap) Value() (driver.Value, error) { return json.Marshal(r) }

func (m *DmlColumnMapSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DmlColumnMapSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDmlColumnMapSet(ss []string) (p DmlColumnMapSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
