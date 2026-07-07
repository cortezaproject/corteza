package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
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
)

func (k NotificationKind) String() string {
	return string(k)
}

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

func (nc *NotificationConfig) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), nc)
}

func (nc NotificationConfig) Value() (driver.Value, error) {
	return json.Marshal(nc)
}
