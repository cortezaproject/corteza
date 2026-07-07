package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
)

type (
	ModuleMapping struct {
		TenantID           uint64                `json:"tenantID,string,omitempty"`
		ProjectID          uint64                `json:"projectID,string,omitempty"`
		NodeID             uint64                `json:"nodeID,string"`
		FederationModuleID uint64                `json:"federationModuleID,string"`
		ComposeModuleID    uint64                `json:"composeModuleID,string"`
		ComposeNamespaceID uint64                `json:"composeNamespaceID,string"`
		FieldMapping       ModuleFieldMappingSet `json:"fields"`
	}

	ModuleFieldMapping struct {
		Origin      ModuleField `json:"origin"`
		Destination ModuleField `json:"destination"`
	}
)

func (r ModuleMapping) Clone() *ModuleMapping {
	dup := r
	return &dup
}

func (r ModuleMapping) Diff(cmp *ModuleMapping) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ModuleMapping{}
	}
	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.NodeID != cmp.NodeID {
		out = append(out, &revisions.Change{Key: "nodeID", Old: []any{cmp.NodeID}, New: []any{r.NodeID}})
	}

	if r.FederationModuleID != cmp.FederationModuleID {
		out = append(out, &revisions.Change{Key: "federationModuleID", Old: []any{cmp.FederationModuleID}, New: []any{r.FederationModuleID}})
	}

	if r.ComposeModuleID != cmp.ComposeModuleID {
		out = append(out, &revisions.Change{Key: "composeModuleID", Old: []any{cmp.ComposeModuleID}, New: []any{r.ComposeModuleID}})
	}

	if r.ComposeNamespaceID != cmp.ComposeNamespaceID {
		out = append(out, &revisions.Change{Key: "composeNamespaceID", Old: []any{cmp.ComposeNamespaceID}, New: []any{r.ComposeNamespaceID}})
	}

	if !reflect.DeepEqual(r.FieldMapping, cmp.FieldMapping) {
		out = append(out, &revisions.Change{Key: "fields", Old: []any{cmp.FieldMapping}, New: []any{r.FieldMapping}})
	}

	return out
}

func (r ModuleFieldMapping) Clone() *ModuleFieldMapping {
	dup := r
	return &dup
}

func (r ModuleFieldMapping) Diff(cmp *ModuleFieldMapping) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ModuleFieldMapping{}
	}
	if !reflect.DeepEqual(r.Origin, cmp.Origin) {
		out = append(out, &revisions.Change{Key: "origin", Old: []any{cmp.Origin}, New: []any{r.Origin}})
	}

	if !reflect.DeepEqual(r.Destination, cmp.Destination) {
		out = append(out, &revisions.Change{Key: "destination", Old: []any{cmp.Destination}, New: []any{r.Destination}})
	}

	return out
}

func (r *ModuleFieldMapping) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ModuleFieldMapping) Value() (driver.Value, error) { return json.Marshal(r) }

func (m *ModuleFieldMappingSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ModuleFieldMappingSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseModuleFieldMappingSet(ss []string) (p ModuleFieldMappingSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
