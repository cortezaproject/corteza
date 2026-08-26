package types

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/spf13/cast"
)

type (
	QueueFilter struct {
		QueueID []string     `json:"queueID"`
		Query   string       `json:"query"`
		Deleted filter.State `json:"deleted"`

		Check func(*Queue) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	QueueMessageFilter struct {
		Queue string

		Processed filter.State `json:"processed"`

		filter.Sorting
		filter.Paging
	}
)

func (h *QueueMeta) UnmarshalJSON(s []byte) error {
	type Alias QueueMeta

	aux := &struct {
		PollDelay string `json:"poll_delay"`
		*Alias
	}{
		Alias: (*Alias)(h),
	}

	h.DispatchEvents = false

	if err := json.Unmarshal(s, aux); err != nil {
		return err
	}

	if d, err := cast.ToDurationE(aux.PollDelay); err == nil {
		h.PollDelay = &d
	}

	return nil
}

func (m QueueMeta) MarshalJSON() ([]byte, error) {
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

func (s *Queue) CanDispatch() bool {
	return s.Meta.DispatchEvents
}
