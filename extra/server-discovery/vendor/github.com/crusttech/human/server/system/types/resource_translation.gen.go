package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/pkg/revisions"
	"reflect"
	"time"
)

type (
	ResourceTranslation struct {
		ID        uint64     `json:"translationID,string"`
		TenantID  uint64     `json:"tenantID,string,omitempty"`
		ProjectID uint64     `json:"projectID,string,omitempty"`
		Lang      Lang       `json:"lang"`
		Resource  string     `json:"resource"`
		K         string     `json:"key"`
		Message   string     `json:"message"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}
)

func (r ResourceTranslation) Clone() *ResourceTranslation {
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

func (r ResourceTranslation) Diff(cmp *ResourceTranslation) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ResourceTranslation{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "translationID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if !reflect.DeepEqual(r.Lang, cmp.Lang) {
		out = append(out, &revisions.Change{Key: "lang", Old: []any{cmp.Lang}, New: []any{r.Lang}})
	}

	if r.Resource != cmp.Resource {
		out = append(out, &revisions.Change{Key: "resource", Old: []any{cmp.Resource}, New: []any{r.Resource}})
	}

	if r.K != cmp.K {
		out = append(out, &revisions.Change{Key: "key", Old: []any{cmp.K}, New: []any{r.K}})
	}

	if r.Message != cmp.Message {
		out = append(out, &revisions.Change{Key: "message", Old: []any{cmp.Message}, New: []any{r.Message}})
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

	if r.OwnedBy != cmp.OwnedBy {
		out = append(out, &revisions.Change{Key: "ownedBy", Old: []any{cmp.OwnedBy}, New: []any{r.OwnedBy}})
	}

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	if r.DeletedBy != cmp.DeletedBy {
		out = append(out, &revisions.Change{Key: "deletedBy", Old: []any{cmp.DeletedBy}, New: []any{r.DeletedBy}})
	}

	return out
}
