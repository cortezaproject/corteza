package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	DalSchemaAlteration struct {
		ID           uint64                     `json:"alterationID,string"`
		BatchID      uint64                     `json:"batchID,string"`
		DependsOn    uint64                     `json:"dependsOn,string,omitempty"`
		Resource     string                     `json:"resource"`
		ResourceType string                     `json:"resourceType"`
		ConnectionID uint64                     `json:"connectionID,string"`
		Kind         string                     `json:"kind"`
		Params       *DalSchemaAlterationParams `json:"params"`
		Error        string                     `json:"error,omitempty"`
		CreatedAt    time.Time                  `json:"createdAt,omitempty"`
		UpdatedAt    *time.Time                 `json:"updatedAt,omitempty"`
		DeletedAt    *time.Time                 `json:"deletedAt,omitempty"`
		CompletedAt  *time.Time                 `json:"completedAt,omitempty"`
		DismissedAt  *time.Time                 `json:"dismissedAt,omitempty"`
		CreatedBy    uint64                     `json:"createdBy,string"`
		UpdatedBy    uint64                     `json:"updatedBy,string,omitempty"`
		DeletedBy    uint64                     `json:"deletedBy,string,omitempty"`
		CompletedBy  uint64                     `json:"completedBy,string,omitempty"`
		DismissedBy  uint64                     `json:"dismissedBy,string,omitempty"`
	}

	DalSchemaAlterationParams struct {
		AttributeAdd      *dal.AttributeAdd      `json:"attributeAdd,omitempty"`
		AttributeDelete   *dal.AttributeDelete   `json:"attributeDelete,omitempty"`
		AttributeReType   *dal.AttributeReType   `json:"attributeReType,omitempty"`
		AttributeReEncode *dal.AttributeReEncode `json:"attributeReEncode,omitempty"`
		ModelAdd          *dal.ModelAdd          `json:"modelAdd,omitempty"`
		ModelDelete       *dal.ModelDelete       `json:"modelDelete,omitempty"`
	}
)

func (r DalSchemaAlteration) Clone() *DalSchemaAlteration {
	dup := r
	if r.Params != nil {
		dup.Params = r.Params.Clone()
	}

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	if r.CompletedAt != nil {
		v := *r.CompletedAt
		dup.CompletedAt = &v
	}

	if r.DismissedAt != nil {
		v := *r.DismissedAt
		dup.DismissedAt = &v
	}

	return &dup
}

func (r DalSchemaAlteration) Diff(cmp *DalSchemaAlteration) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalSchemaAlteration{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "alterationID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.BatchID != cmp.BatchID {
		out = append(out, &revisions.Change{Key: "batchID", Old: []any{cmp.BatchID}, New: []any{r.BatchID}})
	}

	if r.DependsOn != cmp.DependsOn {
		out = append(out, &revisions.Change{Key: "dependsOn", Old: []any{cmp.DependsOn}, New: []any{r.DependsOn}})
	}

	if r.Resource != cmp.Resource {
		out = append(out, &revisions.Change{Key: "resource", Old: []any{cmp.Resource}, New: []any{r.Resource}})
	}

	if r.ResourceType != cmp.ResourceType {
		out = append(out, &revisions.Change{Key: "resourceType", Old: []any{cmp.ResourceType}, New: []any{r.ResourceType}})
	}

	if r.ConnectionID != cmp.ConnectionID {
		out = append(out, &revisions.Change{Key: "connectionID", Old: []any{cmp.ConnectionID}, New: []any{r.ConnectionID}})
	}

	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if (r.Params == nil) != (cmp.Params == nil) {
		out = append(out, &revisions.Change{Key: "params", Old: []any{cmp.Params}, New: []any{r.Params}})
	} else if r.Params != nil {
		for _, c := range r.Params.Diff(cmp.Params) {
			c.Key = "params." + c.Key
			out = append(out, c)
		}
	}

	if r.Error != cmp.Error {
		out = append(out, &revisions.Change{Key: "error", Old: []any{cmp.Error}, New: []any{r.Error}})
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	if !reflect.DeepEqual(r.CompletedAt, cmp.CompletedAt) {
		out = append(out, &revisions.Change{Key: "completedAt", Old: []any{cmp.CompletedAt}, New: []any{r.CompletedAt}})
	}

	if !reflect.DeepEqual(r.DismissedAt, cmp.DismissedAt) {
		out = append(out, &revisions.Change{Key: "dismissedAt", Old: []any{cmp.DismissedAt}, New: []any{r.DismissedAt}})
	}

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	if r.DeletedBy != cmp.DeletedBy {
		out = append(out, &revisions.Change{Key: "deletedBy", Old: []any{cmp.DeletedBy}, New: []any{r.DeletedBy}})
	}

	if r.CompletedBy != cmp.CompletedBy {
		out = append(out, &revisions.Change{Key: "completedBy", Old: []any{cmp.CompletedBy}, New: []any{r.CompletedBy}})
	}

	if r.DismissedBy != cmp.DismissedBy {
		out = append(out, &revisions.Change{Key: "dismissedBy", Old: []any{cmp.DismissedBy}, New: []any{r.DismissedBy}})
	}

	return out
}

func (r DalSchemaAlterationParams) Clone() *DalSchemaAlterationParams {
	dup := r
	if r.AttributeAdd != nil {
		v := *r.AttributeAdd
		dup.AttributeAdd = &v
	}

	if r.AttributeDelete != nil {
		v := *r.AttributeDelete
		dup.AttributeDelete = &v
	}

	if r.AttributeReType != nil {
		v := *r.AttributeReType
		dup.AttributeReType = &v
	}

	if r.AttributeReEncode != nil {
		v := *r.AttributeReEncode
		dup.AttributeReEncode = &v
	}

	if r.ModelAdd != nil {
		v := *r.ModelAdd
		dup.ModelAdd = &v
	}

	if r.ModelDelete != nil {
		v := *r.ModelDelete
		dup.ModelDelete = &v
	}

	return &dup
}

func (r DalSchemaAlterationParams) Diff(cmp *DalSchemaAlterationParams) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalSchemaAlterationParams{}
	}
	if !reflect.DeepEqual(r.AttributeAdd, cmp.AttributeAdd) {
		out = append(out, &revisions.Change{Key: "attributeAdd", Old: []any{cmp.AttributeAdd}, New: []any{r.AttributeAdd}})
	}

	if !reflect.DeepEqual(r.AttributeDelete, cmp.AttributeDelete) {
		out = append(out, &revisions.Change{Key: "attributeDelete", Old: []any{cmp.AttributeDelete}, New: []any{r.AttributeDelete}})
	}

	if !reflect.DeepEqual(r.AttributeReType, cmp.AttributeReType) {
		out = append(out, &revisions.Change{Key: "attributeReType", Old: []any{cmp.AttributeReType}, New: []any{r.AttributeReType}})
	}

	if !reflect.DeepEqual(r.AttributeReEncode, cmp.AttributeReEncode) {
		out = append(out, &revisions.Change{Key: "attributeReEncode", Old: []any{cmp.AttributeReEncode}, New: []any{r.AttributeReEncode}})
	}

	if !reflect.DeepEqual(r.ModelAdd, cmp.ModelAdd) {
		out = append(out, &revisions.Change{Key: "modelAdd", Old: []any{cmp.ModelAdd}, New: []any{r.ModelAdd}})
	}

	if !reflect.DeepEqual(r.ModelDelete, cmp.ModelDelete) {
		out = append(out, &revisions.Change{Key: "modelDelete", Old: []any{cmp.ModelDelete}, New: []any{r.ModelDelete}})
	}

	return out
}

func (r *DalSchemaAlterationParams) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DalSchemaAlterationParams) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseDalSchemaAlterationParams(ss []string) (p DalSchemaAlterationParams, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
