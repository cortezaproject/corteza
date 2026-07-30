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
	ProjectAiSystemEntry struct {
		ID                uint64    `json:"-"`
		ProjectAiSystemID uint64    `json:"projectAiSystemID,string"`
		ResourceRef       string    `json:"resourceRef"`
		CreatedAt         time.Time `json:"createdAt,omitempty"`
	}
)

func (r ProjectAiSystemEntry) Clone() *ProjectAiSystemEntry {
	dup := r
	return &dup
}

func (r ProjectAiSystemEntry) Diff(cmp *ProjectAiSystemEntry) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectAiSystemEntry{}
	}
	if r.ProjectAiSystemID != cmp.ProjectAiSystemID {
		out = append(out, &revisions.Change{Key: "projectAiSystemID", Old: []any{cmp.ProjectAiSystemID}, New: []any{r.ProjectAiSystemID}})
	}

	if r.ResourceRef != cmp.ResourceRef {
		out = append(out, &revisions.Change{Key: "resourceRef", Old: []any{cmp.ResourceRef}, New: []any{r.ResourceRef}})
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	return out
}
