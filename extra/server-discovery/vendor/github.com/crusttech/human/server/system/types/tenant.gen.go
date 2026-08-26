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
	Tenant struct {
		ID          uint64                           `json:"tenantID,string"`
		Handle      string                           `json:"handle"`
		Status      TenantStatus                     `json:"status"`
		Config      TenantConfig                     `json:"config"`
		Meta        TenantMeta                       `json:"meta"`
		CreatedAt   time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time                       `json:"updatedAt,omitempty"`
		SuspendedAt *time.Time                       `json:"suspendedAt,omitempty"`
		DeletedAt   *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy   uint64                           `json:"createdBy,string"`
		UpdatedBy   uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy   uint64                           `json:"deletedBy,string,omitempty"`
		Labels      map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	TenantConfig struct {
		DALConnectionID uint64          `json:"dalConnectionID,string,omitempty"`
		AuthProviders   []string        `json:"authProviders,omitempty"`
		FeatureFlags    map[string]bool `json:"featureFlags,omitempty"`
		Quotas          TenantQuotas    `json:"quotas,omitempty"`
		Locale          string          `json:"locale,omitempty"`
		Timezone        string          `json:"timezone,omitempty"`
	}

	TenantQuotas struct {
		MaxUsers    int   `json:"maxUsers,omitempty"`
		MaxProjects int   `json:"maxProjects,omitempty"`
		MaxStorage  int64 `json:"maxStorage,omitempty"`
	}

	TenantMeta struct {
		Short       string   `json:"short,omitempty"`
		Description string   `json:"description,omitempty"`
		LogoID      uint64   `json:"logoID,string,omitempty"`
		Color       string   `json:"color,omitempty"`
		Tags        []string `json:"tags,omitempty"`
	}

	TenantStatus string
)

func (r Tenant) Clone() *Tenant {
	dup := r
	dup.Config = *r.Config.Clone()

	dup.Meta = *r.Meta.Clone()

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.SuspendedAt != nil {
		v := *r.SuspendedAt
		dup.SuspendedAt = &v
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

func (r Tenant) Diff(cmp *Tenant) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Tenant{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if !reflect.DeepEqual(r.Status, cmp.Status) {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	for _, c := range r.Config.Diff(&cmp.Config) {
		c.Key = "config." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.SuspendedAt, cmp.SuspendedAt) {
		out = append(out, &revisions.Change{Key: "suspendedAt", Old: []any{cmp.SuspendedAt}, New: []any{r.SuspendedAt}})
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

func (r TenantConfig) Clone() *TenantConfig {
	dup := r
	if r.AuthProviders != nil {
		dup.AuthProviders = make([]string, len(r.AuthProviders))
		copy(dup.AuthProviders, r.AuthProviders)
	}

	if r.FeatureFlags != nil {
		dup.FeatureFlags = make(map[string]bool, len(r.FeatureFlags))
		for k, v := range r.FeatureFlags {
			dup.FeatureFlags[k] = v
		}
	}

	dup.Quotas = *r.Quotas.Clone()

	return &dup
}

func (r TenantConfig) Diff(cmp *TenantConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &TenantConfig{}
	}
	if r.DALConnectionID != cmp.DALConnectionID {
		out = append(out, &revisions.Change{Key: "dalConnectionID", Old: []any{cmp.DALConnectionID}, New: []any{r.DALConnectionID}})
	}

	if !reflect.DeepEqual(r.AuthProviders, cmp.AuthProviders) {
		out = append(out, &revisions.Change{Key: "authProviders", Old: []any{cmp.AuthProviders}, New: []any{r.AuthProviders}})
	}

	if !reflect.DeepEqual(r.FeatureFlags, cmp.FeatureFlags) {
		out = append(out, &revisions.Change{Key: "featureFlags", Old: []any{cmp.FeatureFlags}, New: []any{r.FeatureFlags}})
	}

	for _, c := range r.Quotas.Diff(&cmp.Quotas) {
		c.Key = "quotas." + c.Key
		out = append(out, c)
	}

	if r.Locale != cmp.Locale {
		out = append(out, &revisions.Change{Key: "locale", Old: []any{cmp.Locale}, New: []any{r.Locale}})
	}

	if r.Timezone != cmp.Timezone {
		out = append(out, &revisions.Change{Key: "timezone", Old: []any{cmp.Timezone}, New: []any{r.Timezone}})
	}

	return out
}

func (r *TenantConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r TenantConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r TenantQuotas) Clone() *TenantQuotas {
	dup := r
	return &dup
}

func (r TenantQuotas) Diff(cmp *TenantQuotas) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &TenantQuotas{}
	}
	if r.MaxUsers != cmp.MaxUsers {
		out = append(out, &revisions.Change{Key: "maxUsers", Old: []any{cmp.MaxUsers}, New: []any{r.MaxUsers}})
	}

	if r.MaxProjects != cmp.MaxProjects {
		out = append(out, &revisions.Change{Key: "maxProjects", Old: []any{cmp.MaxProjects}, New: []any{r.MaxProjects}})
	}

	if r.MaxStorage != cmp.MaxStorage {
		out = append(out, &revisions.Change{Key: "maxStorage", Old: []any{cmp.MaxStorage}, New: []any{r.MaxStorage}})
	}

	return out
}

func (r *TenantQuotas) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r TenantQuotas) Value() (driver.Value, error) { return json.Marshal(r) }

func (r TenantMeta) Clone() *TenantMeta {
	dup := r
	if r.Tags != nil {
		dup.Tags = make([]string, len(r.Tags))
		copy(dup.Tags, r.Tags)
	}

	return &dup
}

func (r TenantMeta) Diff(cmp *TenantMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &TenantMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.LogoID != cmp.LogoID {
		out = append(out, &revisions.Change{Key: "logoID", Old: []any{cmp.LogoID}, New: []any{r.LogoID}})
	}

	if r.Color != cmp.Color {
		out = append(out, &revisions.Change{Key: "color", Old: []any{cmp.Color}, New: []any{r.Color}})
	}

	if !reflect.DeepEqual(r.Tags, cmp.Tags) {
		out = append(out, &revisions.Change{Key: "tags", Old: []any{cmp.Tags}, New: []any{r.Tags}})
	}

	return out
}

func (r *TenantMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r TenantMeta) Value() (driver.Value, error) { return json.Marshal(r) }

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusArchived  TenantStatus = "archived"
)

func ParseTenantConfig(ss []string) (p TenantConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseTenantMeta(ss []string) (p TenantMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
