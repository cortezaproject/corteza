package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	KnowledgeBaseFilter struct {
		KnowledgeBaseID []uint64     `json:"knowledgeBaseID"`
		Handle          string       `json:"handle"`
		Query           string       `json:"query"`
		Deleted         filter.State `json:"deleted"`

		Check func(*KnowledgeBase) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
