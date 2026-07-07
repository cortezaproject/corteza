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
	Notification struct {
		ID        uint64             `json:"notificationID,string"`
		TenantID  uint64             `json:"tenantID,string,omitempty"`
		ProjectID uint64             `json:"projectID,string,omitempty"`
		Kind      NotificationKind   `json:"kind"`
		Config    NotificationConfig `json:"config"`
		Recipient uint64             `json:"recipient,string"`
		CreatedBy uint64             `json:"createdBy,string"`
		ReadAt    *time.Time         `json:"readAt"`
		CreatedAt time.Time          `json:"createdAt,omitempty"`
		UpdatedAt *time.Time         `json:"updatedAt,omitempty"`
		DeletedAt *time.Time         `json:"deletedAt,omitempty"`
	}

	NotificationConfig struct {
		Simple SimpleNotificationConfig `json:"simple"`
		Record RecordNotificationConfig `json:"record"`
	}

	SimpleNotificationConfig struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	RecordNotificationConfig struct {
		Title       string       `json:"title"`
		Description string       `json:"description"`
		ModuleID    uint64       `json:"moduleID,string"`
		NamespaceID uint64       `json:"namespaceID,string"`
		RecordID    uint64       `json:"recordID,string"`
		OpenMode    OpenModeType `json:"openMode,omitempty"`
		Edit        bool         `json:"edit,omitempty"`
	}

	OpenModeType string

	NotificationKind string
)

func (r Notification) Clone() *Notification {
	dup := r
	dup.Config = *r.Config.Clone()

	if r.ReadAt != nil {
		v := *r.ReadAt
		dup.ReadAt = &v
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

func (r Notification) Diff(cmp *Notification) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Notification{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "notificationID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	for _, c := range r.Config.Diff(&cmp.Config) {
		c.Key = "config." + c.Key
		out = append(out, c)
	}

	if r.Recipient != cmp.Recipient {
		out = append(out, &revisions.Change{Key: "recipient", Old: []any{cmp.Recipient}, New: []any{r.Recipient}})
	}

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if !reflect.DeepEqual(r.ReadAt, cmp.ReadAt) {
		out = append(out, &revisions.Change{Key: "readAt", Old: []any{cmp.ReadAt}, New: []any{r.ReadAt}})
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

func (r NotificationConfig) Clone() *NotificationConfig {
	dup := r
	dup.Simple = *r.Simple.Clone()

	dup.Record = *r.Record.Clone()

	return &dup
}

func (r NotificationConfig) Diff(cmp *NotificationConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NotificationConfig{}
	}
	for _, c := range r.Simple.Diff(&cmp.Simple) {
		c.Key = "simple." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Record.Diff(&cmp.Record) {
		c.Key = "record." + c.Key
		out = append(out, c)
	}

	return out
}

func (r SimpleNotificationConfig) Clone() *SimpleNotificationConfig {
	dup := r
	return &dup
}

func (r SimpleNotificationConfig) Diff(cmp *SimpleNotificationConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &SimpleNotificationConfig{}
	}
	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *SimpleNotificationConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r SimpleNotificationConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r RecordNotificationConfig) Clone() *RecordNotificationConfig {
	dup := r
	return &dup
}

func (r RecordNotificationConfig) Diff(cmp *RecordNotificationConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &RecordNotificationConfig{}
	}
	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.ModuleID != cmp.ModuleID {
		out = append(out, &revisions.Change{Key: "moduleID", Old: []any{cmp.ModuleID}, New: []any{r.ModuleID}})
	}

	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	if r.RecordID != cmp.RecordID {
		out = append(out, &revisions.Change{Key: "recordID", Old: []any{cmp.RecordID}, New: []any{r.RecordID}})
	}

	if !reflect.DeepEqual(r.OpenMode, cmp.OpenMode) {
		out = append(out, &revisions.Change{Key: "openMode", Old: []any{cmp.OpenMode}, New: []any{r.OpenMode}})
	}

	if r.Edit != cmp.Edit {
		out = append(out, &revisions.Change{Key: "edit", Old: []any{cmp.Edit}, New: []any{r.Edit}})
	}

	return out
}

func (r *RecordNotificationConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r RecordNotificationConfig) Value() (driver.Value, error) { return json.Marshal(r) }

const (
	OpenModeTypeOpenModeModal   OpenModeType = "modal"
	OpenModeTypeOpenModeNewTab  OpenModeType = "newTab"
	OpenModeTypeOpenModeSameTab OpenModeType = "sameTab"
)

const (
	NotificationKindSimple NotificationKind = "simple"
	NotificationKindRecord NotificationKind = "record"
)
