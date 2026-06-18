package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	KnowledgeBaseContext struct {
		Namespaces []KnowledgeBaseNamespaceContext `json:"namespaces"`
	}

	KnowledgeBaseNamespaceContext struct {
		NamespaceID uint64              `json:"namespaceID,string"`
		ModuleIDs   KnowledgeBaseIDList `json:"moduleIDs"`
	}

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
