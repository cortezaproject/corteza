package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	ConfiguredConnection struct {
		ID           uint64                           `json:"configurationID,string"`
		TenantID     uint64                           `json:"tenantID,string,omitempty"`
		ProjectID    uint64                           `json:"projectID,string,omitempty"`
		ConnectionID uint64                           `json:"connectionID,string"`
		Name         string                           `json:"name"`
		Status       string                           `json:"status"`
		Connection   Connection                       `json:"connection"`
		Config       ConfiguredConnectionConfig       `json:"config"`
		CreatedAt    time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt    *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt    *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy    uint64                           `json:"createdBy,string"`
		UpdatedBy    uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy    uint64                           `json:"deletedBy,string,omitempty"`
		Labels       map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	ConfiguredConnectionConfig struct {
		NamespaceID     uint64                      `json:"namespaceID,string"`
		DalConnectionID uint64                      `json:"dalConnectionID,string"`
		CredentialID    uint64                      `json:"credentialID,string"`
		Params          []ConfiguredConnectionParam `json:"params,omitempty"`
		Discovery       map[string]json.RawMessage  `json:"discovery,omitempty"`
	}

	ConfiguredConnectionParam struct {
		Scope []string `json:"scope"`
		Name  string   `json:"name"`
		Value string   `json:"value"`
	}

	ConfiguredConnectionCheckStatus struct {
		OK      bool   `json:"ok"`
		Message string `json:"message,omitempty"`
	}

	ConfiguredConnectionCheckResult struct {
		Connectivity ConfiguredConnectionCheckStatus  `json:"connectivity"`
		Auth         ConfiguredConnectionCheckStatus  `json:"auth"`
		Probe        *ConfiguredConnectionCheckStatus `json:"probe,omitempty"`
	}
)

func (r ConfiguredConnection) Clone() *ConfiguredConnection {
	dup := r
	dup.Config = *r.Config.Clone()

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}
	return &dup
}

func (r ConfiguredConnection) Diff(cmp *ConfiguredConnection) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConfiguredConnection{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "configurationID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.ConnectionID != cmp.ConnectionID {
		out = append(out, &revisions.Change{Key: "connectionID", Old: []any{cmp.ConnectionID}, New: []any{r.ConnectionID}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if !reflect.DeepEqual(r.Connection, cmp.Connection) {
		out = append(out, &revisions.Change{Key: "connection", Old: []any{cmp.Connection}, New: []any{r.Connection}})
	}

	for _, c := range r.Config.Diff(&cmp.Config) {
		c.Key = "config." + c.Key
		out = append(out, c)
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

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	if r.DeletedBy != cmp.DeletedBy {
		out = append(out, &revisions.Change{Key: "deletedBy", Old: []any{cmp.DeletedBy}, New: []any{r.DeletedBy}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r ConfiguredConnectionConfig) Clone() *ConfiguredConnectionConfig {
	dup := r
	if r.Params != nil {
		dup.Params = make([]ConfiguredConnectionParam, len(r.Params))
		for i := range r.Params {
			dup.Params[i] = *r.Params[i].Clone()
		}
	}

	if r.Discovery != nil {
		dup.Discovery = make(map[string]json.RawMessage, len(r.Discovery))
		for k, v := range r.Discovery {
			dup.Discovery[k] = v
		}
	}

	return &dup
}

func (r ConfiguredConnectionConfig) Diff(cmp *ConfiguredConnectionConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConfiguredConnectionConfig{}
	}
	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	if r.DalConnectionID != cmp.DalConnectionID {
		out = append(out, &revisions.Change{Key: "dalConnectionID", Old: []any{cmp.DalConnectionID}, New: []any{r.DalConnectionID}})
	}

	if r.CredentialID != cmp.CredentialID {
		out = append(out, &revisions.Change{Key: "credentialID", Old: []any{cmp.CredentialID}, New: []any{r.CredentialID}})
	}

	if !reflect.DeepEqual(r.Params, cmp.Params) {
		out = append(out, &revisions.Change{Key: "params", Old: []any{cmp.Params}, New: []any{r.Params}})
	}

	if !reflect.DeepEqual(r.Discovery, cmp.Discovery) {
		out = append(out, &revisions.Change{Key: "discovery", Old: []any{cmp.Discovery}, New: []any{r.Discovery}})
	}

	return out
}

func (r *ConfiguredConnectionConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConfiguredConnectionConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConfiguredConnectionParam) Clone() *ConfiguredConnectionParam {
	dup := r
	if r.Scope != nil {
		dup.Scope = make([]string, len(r.Scope))
		copy(dup.Scope, r.Scope)
	}

	return &dup
}

func (r ConfiguredConnectionParam) Diff(cmp *ConfiguredConnectionParam) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConfiguredConnectionParam{}
	}
	if !reflect.DeepEqual(r.Scope, cmp.Scope) {
		out = append(out, &revisions.Change{Key: "scope", Old: []any{cmp.Scope}, New: []any{r.Scope}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Value != cmp.Value {
		out = append(out, &revisions.Change{Key: "value", Old: []any{cmp.Value}, New: []any{r.Value}})
	}

	return out
}

func (r ConfiguredConnectionCheckStatus) Clone() *ConfiguredConnectionCheckStatus {
	dup := r
	return &dup
}

func (r ConfiguredConnectionCheckStatus) Diff(cmp *ConfiguredConnectionCheckStatus) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConfiguredConnectionCheckStatus{}
	}
	if r.OK != cmp.OK {
		out = append(out, &revisions.Change{Key: "ok", Old: []any{cmp.OK}, New: []any{r.OK}})
	}

	if r.Message != cmp.Message {
		out = append(out, &revisions.Change{Key: "message", Old: []any{cmp.Message}, New: []any{r.Message}})
	}

	return out
}

func (r *ConfiguredConnectionCheckStatus) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConfiguredConnectionCheckStatus) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConfiguredConnectionCheckResult) Clone() *ConfiguredConnectionCheckResult {
	dup := r
	dup.Connectivity = *r.Connectivity.Clone()

	dup.Auth = *r.Auth.Clone()

	if r.Probe != nil {
		dup.Probe = r.Probe.Clone()
	}

	return &dup
}

func (r ConfiguredConnectionCheckResult) Diff(cmp *ConfiguredConnectionCheckResult) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConfiguredConnectionCheckResult{}
	}
	for _, c := range r.Connectivity.Diff(&cmp.Connectivity) {
		c.Key = "connectivity." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Auth.Diff(&cmp.Auth) {
		c.Key = "auth." + c.Key
		out = append(out, c)
	}

	if (r.Probe == nil) != (cmp.Probe == nil) {
		out = append(out, &revisions.Change{Key: "probe", Old: []any{cmp.Probe}, New: []any{r.Probe}})
	} else if r.Probe != nil {
		for _, c := range r.Probe.Diff(cmp.Probe) {
			c.Key = "probe." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *ConfiguredConnectionCheckResult) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConfiguredConnectionCheckResult) Value() (driver.Value, error) { return json.Marshal(r) }

func (m *Connection) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m Connection) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConnection(ss []string) (p Connection, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseConfiguredConnectionConfig(ss []string) (p ConfiguredConnectionConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
