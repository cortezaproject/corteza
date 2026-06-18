package types

import (
	"github.com/crusttech/human/server/pkg/dal"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
	DalSchemaAlterationParams struct {
		AttributeAdd      *dal.AttributeAdd      `json:"attributeAdd,omitempty"`
		AttributeDelete   *dal.AttributeDelete   `json:"attributeDelete,omitempty"`
		AttributeReType   *dal.AttributeReType   `json:"attributeReType,omitempty"`
		AttributeReEncode *dal.AttributeReEncode `json:"attributeReEncode,omitempty"`
		ModelAdd          *dal.ModelAdd          `json:"modelAdd,omitempty"`
		ModelDelete       *dal.ModelDelete       `json:"modelDelete,omitempty"`
	}

	DalSchemaAlterationFilter struct {
		AlterationID []string `json:"alterationID"`
		BatchID      []uint64 `json:"batchID,string"`
		Kind         string   `json:"kind"`
		Resource     []string `json:"resource"`
		ResourceType string   `json:"resourceType"`

		Deleted   filter.State `json:"deleted"`
		Completed filter.State `json:"completed"`
		Dismissed filter.State `json:"dismissed"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)
