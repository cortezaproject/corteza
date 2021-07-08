package messagebus

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/spf13/cast"
)

type (
	QueueSettingsMeta struct {
		PollDelay      *time.Duration `json:"poll_delay"`
		DispatchEvents bool           `json:"dispatch_events"`
	}

	QueueSettings struct {
		ID       uint64            `json:"queueID,string"`
		Consumer string            `json:"consumer"`
		Queue    string            `json:"queue"`
		Meta     QueueSettingsMeta `json:"meta"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	QueueSettingsFilter struct {
		Queue    string       `json:"queue"`
		Consumer ConsumerType `json:"handler"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*QueueSettings) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (h *QueueSettingsMeta) UnmarshalJSON(s []byte) error {
	type Alias QueueSettingsMeta

	aux := &struct {
		PollDelay string `json:"poll_delay"`
		*Alias
	}{
		Alias: (*Alias)(h),
	}

	// set default
	h.DispatchEvents = false

	if err := json.Unmarshal(s, aux); err != nil {
		return err
	}

	if d, err := cast.ToDurationE(aux.PollDelay); err == nil {
		h.PollDelay = &d
	}

	return nil
}

func (m QueueSettingsMeta) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m QueueSettingsMeta) MarshalJSON() ([]byte, error) {

	pollDelay := ""
	if m.PollDelay != nil {
		pollDelay = m.PollDelay.String()
	}

	return json.Marshal(struct {
		PollDelay      string `json:"poll_delay"`
		DispatchEvents bool   `json:"dispatch_events"`
	}{
		PollDelay:      pollDelay,
		DispatchEvents: m.DispatchEvents,
	})
}

func (m *QueueSettingsMeta) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		*m = QueueSettingsMeta{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), m); err != nil {
			return errors.New(fmt.Sprintf("cannot scan '%v' into QueueSettingsMeta", value))
		}
	}

	return nil
}

func (s *QueueSettings) CanDispatch() bool {
	return s.Meta.DispatchEvents
}

func ParseQueueSettingsMeta(ss []string) (p QueueSettingsMeta, err error) {
	p = QueueSettingsMeta{}

	if len(ss) == 0 {
		return
	}

	return p, json.Unmarshal([]byte(ss[0]), &p)
}
