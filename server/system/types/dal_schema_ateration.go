package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/sql"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	DalSchemaAlteration struct {
		ID           uint64 `json:"alterationID,string"`
		BatchID      uint64 `json:"batchID,string"`
		DependsOn    uint64 `json:"dependsOn,string,omitempty"`
		ConnectionID uint64 `json:"connectionID,string"`
		Resource     string `json:"resource"`
		ResourceType string `json:"resourceType"`

		Kind   string                     `json:"kind"`
		Params *DalSchemaAlterationParams `json:"params"`

		Error string `json:"error,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`

		CompletedAt *time.Time `json:"completedAt,omitempty"`
		CompletedBy uint64     `json:"completedBy,string,omitempty"`
		DismissedAt *time.Time `json:"dismissedAt,omitempty"`
		DismissedBy uint64     `json:"dismissedBy,string,omitempty"`
	}

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
		BatchID      uint64   `json:"batchID,string"`
		Kind         string   `json:"kind"`
		Resource     []string `json:"resource"`

		Deleted   filter.State `json:"deleted"`
		Completed filter.State `json:"completed"`
		Dismissed filter.State `json:"dismissed"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

func (meta *DalSchemaAlterationParams) Scan(src any) error           { return sql.ParseJSON(src, meta) }
func (meta *DalSchemaAlterationParams) Value() (driver.Value, error) { return json.Marshal(meta) }
