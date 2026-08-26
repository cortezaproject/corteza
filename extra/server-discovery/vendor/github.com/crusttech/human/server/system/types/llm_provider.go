package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	LlmProviderFilter struct {
		LlmProviderID []uint64     `json:"llmProviderID"`
		Handle        string       `json:"handle"`
		Status        string       `json:"status"`
		Provider      string       `json:"provider"`
		Deleted       filter.State `json:"deleted"`

		Check func(*LlmProvider) (bool, error) `json:"-"`

		filter.Paging
		filter.Sorting
	}
)
