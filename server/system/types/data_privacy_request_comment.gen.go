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
	DataPrivacyRequestComment struct {
		ID        uint64     `json:"commentID,string"`
		RequestID uint64     `json:"requestID,string"`
		Comment   string     `json:"comment"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}
)

func (r DataPrivacyRequestComment) Clone() *DataPrivacyRequestComment {
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

func (r DataPrivacyRequestComment) Diff(cmp *DataPrivacyRequestComment) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DataPrivacyRequestComment{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "commentID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.RequestID != cmp.RequestID {
		out = append(out, &revisions.Change{Key: "requestID", Old: []any{cmp.RequestID}, New: []any{r.RequestID}})
	}

	if r.Comment != cmp.Comment {
		out = append(out, &revisions.Change{Key: "comment", Old: []any{cmp.Comment}, New: []any{r.Comment}})
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
