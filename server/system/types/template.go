package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	TemplateFilter struct {
		TemplateID []string `json:"templateID"`
		Query      string   `json:"query"`
		Handle     string   `json:"handle"`
		Type       string   `json:"type"`
		OwnerID    uint64   `json:"ownerID,string"`
		Partial    bool     `json:"partial"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*Template) (bool, error) `json:"-"`

		Deleted filter.State `json:"deleted"`

		filter.Sorting
		filter.Paging
	}
)
