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
		ID      uint64
		Handler string
		Queue   string
		Meta    QueueSettingsMeta

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	QueueSettingsFilter struct {
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
	h.DispatchEvents = true

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
	return json.Marshal(struct {
		PollDelay      string `json:"poll_delay"`
		DispatchEvents bool   `json:"dispatch_events,omitempty"`
	}{
		PollDelay:      m.PollDelay.String(),
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
