package messagebus

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	QueueMessage struct {
		ID        uint64     `json:"messageID"`
		Queue     string     `json:"queue"`
		Payload   []byte     `json:"payload"`
		Created   *time.Time `json:"created"`
		Processed *time.Time `json:"processed"`
	}

	QueueMessageFilter struct {
		Queue string

		Processed filter.State `json:"processed"`

		filter.Sorting
		filter.Paging
	}
)
