package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	Template struct {
		ID         uint64                           `json:"templateID,string"`
		TenantID   uint64                           `json:"tenantID,string,omitempty"`
		ProjectID  uint64                           `json:"projectID,string,omitempty"`
		OwnerID    uint64                           `json:"ownerID,string"`
		Handle     string                           `json:"handle"`
		Language   string                           `json:"language"`
		Type       DocumentType                     `json:"type"`
		Partial    bool                             `json:"partial"`
		Meta       TemplateMeta                     `json:"meta"`
		Template   string                           `json:"template"`
		CreatedAt  time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
		LastUsedAt *time.Time                       `json:"lastUsedAt,omitempty"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	TemplateMeta struct {
		Short       string `json:"short"`
		Description string `json:"description,omitempty"`
	}

	DocumentType string
)

func (r Template) Clone() *Template {
	dup := r
	dup.Meta = *r.Meta.Clone()

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	if r.LastUsedAt != nil {
		v := *r.LastUsedAt
		dup.LastUsedAt = &v
	}

	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}
	return &dup
}

func (r Template) Diff(cmp *Template) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Template{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "templateID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.OwnerID != cmp.OwnerID {
		out = append(out, &revisions.Change{Key: "ownerID", Old: []any{cmp.OwnerID}, New: []any{r.OwnerID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.Language != cmp.Language {
		out = append(out, &revisions.Change{Key: "language", Old: []any{cmp.Language}, New: []any{r.Language}})
	}

	if !reflect.DeepEqual(r.Type, cmp.Type) {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Partial != cmp.Partial {
		out = append(out, &revisions.Change{Key: "partial", Old: []any{cmp.Partial}, New: []any{r.Partial}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if r.Template != cmp.Template {
		out = append(out, &revisions.Change{Key: "template", Old: []any{cmp.Template}, New: []any{r.Template}})
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

	if !reflect.DeepEqual(r.LastUsedAt, cmp.LastUsedAt) {
		out = append(out, &revisions.Change{Key: "lastUsedAt", Old: []any{cmp.LastUsedAt}, New: []any{r.LastUsedAt}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r TemplateMeta) Clone() *TemplateMeta {
	dup := r
	return &dup
}

func (r TemplateMeta) Diff(cmp *TemplateMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &TemplateMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *TemplateMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r TemplateMeta) Value() (driver.Value, error) { return json.Marshal(r) }

const (
	DocumentTypePlain DocumentType = "text/plain"
	DocumentTypeHTML  DocumentType = "text/html"
	DocumentTypePDF   DocumentType = "application/pdf"
)

func ParseTemplateMeta(ss []string) (p TemplateMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
