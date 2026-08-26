package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ConfiguredConnectionFilter struct {
		ConnectionID uint64   `json:"connectionID,string"`
		ProjectID    uint64   `json:"projectID,string"`
		Status       []string `json:"status"`
		Query        string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*ConfiguredConnection) (bool, error) `json:"-"`

		filter.Paging
		filter.Sorting
	}
)
