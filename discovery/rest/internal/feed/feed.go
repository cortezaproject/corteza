package feed

import (
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"time"
)

type (
	Response struct {
		Filter       Filter        `json:"filter"`
		ActivityLogs []ActivityLog `json:"activityLogs"`
	}

	ActivityLog struct {
		ID             uint64    `json:"activityID,string"`
		ResourceID     uint64    `json:"resourceID,string"`
		ResourceType   string    `json:"resourceType"`
		ResourceAction string    `json:"resourceAction"`
		Timestamp      time.Time `json:"timestamp"`
	}

	Filter struct {
		Limit    uint                 `json:"limit"`
		NextPage *filter.PagingCursor `json:"nextPage"`
	}
)
