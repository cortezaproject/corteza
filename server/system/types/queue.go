package types

import (
	"encoding/json"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/spf13/cast"
)

type (
	QueueFilter struct {
		QueueID []string     `json:"queueID"`
		Query   string       `json:"query"`
		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
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

	QueueMeta struct {
		PollDelay      *time.Duration `json:"poll_delay"`
		DispatchEvents bool           `json:"dispatch_events"`
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
