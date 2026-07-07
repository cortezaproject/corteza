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
	ProjectGroupEntry struct {
		ID             uint64    `json:"-"`
		ProjectGroupID uint64    `json:"projectGroupID,string"`
		ResourceRef    string    `json:"resourceRef"`
		CreatedAt      time.Time `json:"createdAt,omitempty"`
	}
)

func (r ProjectGroupEntry) Clone() *ProjectGroupEntry {
	dup := r
	return &dup
}

func (r ProjectGroupEntry) Diff(cmp *ProjectGroupEntry) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectGroupEntry{}
	}
	if r.ProjectGroupID != cmp.ProjectGroupID {
		out = append(out, &revisions.Change{Key: "projectGroupID", Old: []any{cmp.ProjectGroupID}, New: []any{r.ProjectGroupID}})
	}

	if r.ResourceRef != cmp.ResourceRef {
		out = append(out, &revisions.Change{Key: "resourceRef", Old: []any{cmp.ResourceRef}, New: []any{r.ResourceRef}})
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	return out
}
