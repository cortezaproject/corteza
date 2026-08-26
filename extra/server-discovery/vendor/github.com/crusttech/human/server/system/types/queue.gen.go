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
	Queue struct {
		ID        uint64     `json:"queueID,string"`
		TenantID  uint64     `json:"tenantID,string,omitempty"`
		ProjectID uint64     `json:"projectID,string,omitempty"`
		Consumer  string     `json:"consumer"`
		Queue     string     `json:"queue"`
		Meta      QueueMeta  `json:"meta"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	QueueMeta struct {
		PollDelay      *time.Duration `json:"poll_delay"`
		DispatchEvents bool           `json:"dispatch_events"`
	}
)

func (r Queue) Clone() *Queue {
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

	return &dup
}

func (r Queue) Diff(cmp *Queue) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Queue{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "queueID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Consumer != cmp.Consumer {
		out = append(out, &revisions.Change{Key: "consumer", Old: []any{cmp.Consumer}, New: []any{r.Consumer}})
	}

	if r.Queue != cmp.Queue {
		out = append(out, &revisions.Change{Key: "queue", Old: []any{cmp.Queue}, New: []any{r.Queue}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
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

func (r QueueMeta) Clone() *QueueMeta {
	dup := r
	if r.PollDelay != nil {
		v := *r.PollDelay
		dup.PollDelay = &v
	}

	return &dup
}

func (r QueueMeta) Diff(cmp *QueueMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &QueueMeta{}
	}
	if !reflect.DeepEqual(r.PollDelay, cmp.PollDelay) {
		out = append(out, &revisions.Change{Key: "poll_delay", Old: []any{cmp.PollDelay}, New: []any{r.PollDelay}})
	}

	if r.DispatchEvents != cmp.DispatchEvents {
		out = append(out, &revisions.Change{Key: "dispatch_events", Old: []any{cmp.DispatchEvents}, New: []any{r.DispatchEvents}})
	}

	return out
}

func (r *QueueMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r QueueMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseQueueMeta(ss []string) (p QueueMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
