package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
)

type (
	KnowledgeBaseFilter struct {
		KnowledgeBaseID id.Uint64s   `json:"knowledgeBaseID"`
		Handle          string       `json:"handle"`
		Query           string       `json:"query"`
		Deleted         filter.State `json:"deleted"`

		Check func(*KnowledgeBase) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
