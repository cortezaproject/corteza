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
	Node struct {
		ID           uint64     `json:"nodeID,string"`
		TenantID     uint64     `json:"tenantID,string,omitempty"`
		SharedNodeID uint64     `json:"sharedNodeID,string"`
		Name         string     `json:"name"`
		BaseURL      string     `json:"baseURL"`
		Status       string     `json:"status"`
		Contact      string     `json:"contact"`
		PairToken    string     `json:"-"`
		AuthToken    string     `json:"-"`
		CreatedAt    time.Time  `json:"createdAt,omitempty"`
		UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
		DeletedAt    *time.Time `json:"deletedAt,omitempty"`
		CreatedBy    uint64     `json:"createdBy,string"`
		UpdatedBy    uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy    uint64     `json:"deletedBy,string,omitempty"`
	}
)

func (r Node) Clone() *Node {
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

func (r Node) Diff(cmp *Node) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Node{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "nodeID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.SharedNodeID != cmp.SharedNodeID {
		out = append(out, &revisions.Change{Key: "sharedNodeID", Old: []any{cmp.SharedNodeID}, New: []any{r.SharedNodeID}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.BaseURL != cmp.BaseURL {
		out = append(out, &revisions.Change{Key: "baseURL", Old: []any{cmp.BaseURL}, New: []any{r.BaseURL}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.Contact != cmp.Contact {
		out = append(out, &revisions.Change{Key: "contact", Old: []any{cmp.Contact}, New: []any{r.Contact}})
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
