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
	DataPrivacyRequest struct {
		ID          uint64                       `json:"requestID,string"`
		TenantID    uint64                       `json:"tenantID,string,omitempty"`
		ProjectID   uint64                       `json:"projectID,string,omitempty"`
		Kind        RequestKind                  `json:"kind"`
		Status      RequestStatus                `json:"status"`
		Payload     DataPrivacyRequestPayloadSet `json:"payload,omitempty"`
		RequestedAt time.Time                    `json:"requestedAt,omitempty"`
		RequestedBy uint64                       `json:"requestedBy,string"`
		CompletedAt *time.Time                   `json:"completedAt,omitempty"`
		CompletedBy uint64                       `json:"completedBy,string,omitempty"`
		CreatedAt   time.Time                    `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time                   `json:"updatedAt,omitempty"`
		DeletedAt   *time.Time                   `json:"deletedAt,omitempty"`
		CreatedBy   uint64                       `json:"createdBy,string"`
		UpdatedBy   uint64                       `json:"updatedBy,string,omitempty"`
		DeletedBy   uint64                       `json:"deletedBy,string,omitempty"`
	}

	RequestKind string

	RequestStatus string
)

func (r DataPrivacyRequest) Clone() *DataPrivacyRequest {
	dup := r
	if r.CompletedAt != nil {
		v := *r.CompletedAt
		dup.CompletedAt = &v
	}

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

func (r DataPrivacyRequest) Diff(cmp *DataPrivacyRequest) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DataPrivacyRequest{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "requestID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if !reflect.DeepEqual(r.Kind, cmp.Kind) {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if !reflect.DeepEqual(r.Status, cmp.Status) {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if !reflect.DeepEqual(r.Payload, cmp.Payload) {
		out = append(out, &revisions.Change{Key: "payload", Old: []any{cmp.Payload}, New: []any{r.Payload}})
	}

	if !reflect.DeepEqual(r.RequestedAt, cmp.RequestedAt) {
		out = append(out, &revisions.Change{Key: "requestedAt", Old: []any{cmp.RequestedAt}, New: []any{r.RequestedAt}})
	}

	if r.RequestedBy != cmp.RequestedBy {
		out = append(out, &revisions.Change{Key: "requestedBy", Old: []any{cmp.RequestedBy}, New: []any{r.RequestedBy}})
	}

	if !reflect.DeepEqual(r.CompletedAt, cmp.CompletedAt) {
		out = append(out, &revisions.Change{Key: "completedAt", Old: []any{cmp.CompletedAt}, New: []any{r.CompletedAt}})
	}

	if r.CompletedBy != cmp.CompletedBy {
		out = append(out, &revisions.Change{Key: "completedBy", Old: []any{cmp.CompletedBy}, New: []any{r.CompletedBy}})
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

const (
	// RequestKindCorrect to correct module fields
	RequestKindCorrect RequestKind = "correct"
	// RequestKindDelete to delete module fields
	RequestKindDelete RequestKind = "delete"
	// RequestKindExport to export module fields
	RequestKindExport RequestKind = "export"
)

const (
	// RequestStatusPending initially request will be in pending status
	RequestStatusPending RequestStatus = "pending"
	// RequestStatusCanceled owner of request has cancelled the request
	RequestStatusCanceled RequestStatus = "canceled"
	// RequestStatusApproved data officer has of request has cancelled the request
	RequestStatusApproved RequestStatus = "approved"
	// RequestStatusRejected data officer has denied the request
	RequestStatusRejected RequestStatus = "rejected"
)

func (m *DataPrivacyRequestPayloadSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DataPrivacyRequestPayloadSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDataPrivacyRequestPayloadSet(ss []string) (p DataPrivacyRequestPayloadSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
