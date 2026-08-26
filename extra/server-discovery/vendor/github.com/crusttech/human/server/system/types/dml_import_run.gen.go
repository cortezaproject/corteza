package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/pkg/revisions"
	"reflect"
)

type (
	DmlImportRun struct {
		ID           uint64            `json:"runID,string"`
		ConnectionID uint64            `json:"connectionID,string"`
		MappingID    uint64            `json:"mappingID,string"`
		Method       DmlImportMethod   `json:"method"`
		Status       string            `json:"status"`
		Processed    uint64            `json:"processed"`
		Failed       uint64            `json:"failed"`
		Error        string            `json:"error,omitempty"`
		Cursor       map[string]string `json:"cursor,omitempty"`
	}

	DmlImportMethod string
)

func (r DmlImportRun) Clone() *DmlImportRun {
	dup := r
	if r.Cursor != nil {
		dup.Cursor = make(map[string]string, len(r.Cursor))
		for k, v := range r.Cursor {
			dup.Cursor[k] = v
		}
	}

	return &dup
}

func (r DmlImportRun) Diff(cmp *DmlImportRun) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DmlImportRun{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "runID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.ConnectionID != cmp.ConnectionID {
		out = append(out, &revisions.Change{Key: "connectionID", Old: []any{cmp.ConnectionID}, New: []any{r.ConnectionID}})
	}

	if r.MappingID != cmp.MappingID {
		out = append(out, &revisions.Change{Key: "mappingID", Old: []any{cmp.MappingID}, New: []any{r.MappingID}})
	}

	if !reflect.DeepEqual(r.Method, cmp.Method) {
		out = append(out, &revisions.Change{Key: "method", Old: []any{cmp.Method}, New: []any{r.Method}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.Processed != cmp.Processed {
		out = append(out, &revisions.Change{Key: "processed", Old: []any{cmp.Processed}, New: []any{r.Processed}})
	}

	if r.Failed != cmp.Failed {
		out = append(out, &revisions.Change{Key: "failed", Old: []any{cmp.Failed}, New: []any{r.Failed}})
	}

	if r.Error != cmp.Error {
		out = append(out, &revisions.Change{Key: "error", Old: []any{cmp.Error}, New: []any{r.Error}})
	}

	if !reflect.DeepEqual(r.Cursor, cmp.Cursor) {
		out = append(out, &revisions.Change{Key: "cursor", Old: []any{cmp.Cursor}, New: []any{r.Cursor}})
	}

	return out
}

const (
	DmlImportMethodBackground DmlImportMethod = "background"
)
