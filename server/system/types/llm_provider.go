package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
)

type (
	LlmProviderFilter struct {
		LlmProviderID id.Uint64s   `json:"llmProviderID"`
		Handle        string       `json:"handle"`
		Status        string       `json:"status"`
		Provider      string       `json:"provider"`
		Query         string       `json:"query"`
		Deleted       filter.State `json:"deleted"`

		Check func(*LlmProvider) (bool, error) `json:"-"`

		filter.Paging
		filter.Sorting
	}
)
