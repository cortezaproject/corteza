package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	Notification struct {
		ID uint64 `json:"notificationID,string"`

		Kind   NotificationKind   `json:"kind"`
		Config NotificationConfig `json:"config"`

		Recipient uint64    `json:"recipient,string"`
		CreatedBy uint64    `json:"createdBy,string"`
		CreatedAt time.Time `json:"createdAt,omitempty"`

		ReadAt    *time.Time `json:"readAt"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	NotificationFilter struct {
		NotificationID []uint64           `json:"notificationID"`
		Kind           []NotificationKind `json:"kind"`
		Recipient      uint64             `json:"recipient,string"`
		Read           filter.State       `json:"read"`
		Deleted        filter.State       `json:"deleted"`

		Check func(*Notification) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	NotificationKind string

	// NotificationConfig is a flexible configuration container for notification data
	NotificationConfig struct {
		Simple SimpleNotificationConfig `json:"simple"`
		Record RecordNotificationConfig `json:"record"`
	}

	// SimpleNotificationConfig defines the structure for a simple notification
	SimpleNotificationConfig struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	// RecordNotificationConfig defines the structure for a record notification
	RecordNotificationConfig struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ModuleID    uint64 `json:"moduleID,string"`
		NamespaceID uint64 `json:"namespaceID,string"`
		RecordID    uint64 `json:"recordID,string"`
		OpenMode    string `json:"openMode,omitempty"` // "modal", "newTab", or "sameTab" (default)
		Edit        bool   `json:"edit,omitempty"`     // Whether to open the record in edit mode
	}
)

const (
	// NotificationKindSimple is a basic notification with just title and description
	NotificationKindSimple NotificationKind = "simple"
	// NotificationKindRecord is a notification that links to a specific record
	NotificationKindRecord NotificationKind = "record"
)

// String returns the string representation of the notification kind
func (k NotificationKind) String() string {
	return string(k)
}

// CastToNotificationKind converts a string to a typed NotificationKind
func CastToNotificationKind(s string) NotificationKind {
	switch s {
	case string(NotificationKindSimple):
		return NotificationKindSimple
	case string(NotificationKindRecord):
		return NotificationKindRecord
	default:
		return NotificationKindSimple
	}
}

// Scan implements the sql.Scanner interface
func (nc *NotificationConfig) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), nc)
}

// Value implements the driver.Valuer interface
func (nc NotificationConfig) Value() (driver.Value, error) {
	return json.Marshal(nc)
}
