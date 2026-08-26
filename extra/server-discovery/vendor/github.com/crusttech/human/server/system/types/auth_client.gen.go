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
	AuthClient struct {
		ID          uint64                           `json:"authClientID,string"`
		TenantID    uint64                           `json:"tenantID,string,omitempty"`
		ProjectID   uint64                           `json:"projectID,string,omitempty"`
		Handle      string                           `json:"handle"`
		Meta        *AuthClientMeta                  `json:"meta,omitempty"`
		Secret      string                           `json:"secret,omitempty"`
		Scope       string                           `json:"scope"`
		ValidGrant  string                           `json:"validGrant"`
		RedirectURI string                           `json:"redirectURI"`
		Enabled     bool                             `json:"enabled"`
		Trusted     bool                             `json:"trusted"`
		ValidFrom   *time.Time                       `json:"validFrom,omitempty"`
		ExpiresAt   *time.Time                       `json:"expiresAt,omitempty"`
		Security    *AuthClientSecurity              `json:"security"`
		OwnedBy     uint64                           `json:"ownedBy"`
		CreatedAt   time.Time                        `json:"createdAt"`
		UpdatedAt   *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt   *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy   uint64                           `json:"createdBy"`
		UpdatedBy   uint64                           `json:"updatedBy,omitempty"`
		DeletedBy   uint64                           `json:"deletedBy,omitempty"`
		Labels      map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	AuthClientMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	AuthClientSecurity struct {
		ImpersonateUser uint64   `json:"impersonateUser,string,omitempty"`
		UserGroup       uint64   `json:"userGroup,string,omitempty"`
		PermittedRoles  []string `json:"permittedRoles,omitempty"`
		ProhibitedRoles []string `json:"prohibitedRoles,omitempty"`
		ForcedRoles     []string `json:"forcedRoles,omitempty"`
	}
)

func (r AuthClient) Clone() *AuthClient {
	dup := r
	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	if r.ValidFrom != nil {
		v := *r.ValidFrom
		dup.ValidFrom = &v
	}

	if r.ExpiresAt != nil {
		v := *r.ExpiresAt
		dup.ExpiresAt = &v
	}

	if r.Security != nil {
		dup.Security = r.Security.Clone()
	}

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

func (r AuthClient) Diff(cmp *AuthClient) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AuthClient{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "authClientID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if (r.Meta == nil) != (cmp.Meta == nil) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	} else if r.Meta != nil {
		for _, c := range r.Meta.Diff(cmp.Meta) {
			c.Key = "meta." + c.Key
			out = append(out, c)
		}
	}

	if r.Secret != cmp.Secret {
		out = append(out, &revisions.Change{Key: "secret", Old: []any{cmp.Secret}, New: []any{r.Secret}})
	}

	if r.Scope != cmp.Scope {
		out = append(out, &revisions.Change{Key: "scope", Old: []any{cmp.Scope}, New: []any{r.Scope}})
	}

	if r.ValidGrant != cmp.ValidGrant {
		out = append(out, &revisions.Change{Key: "validGrant", Old: []any{cmp.ValidGrant}, New: []any{r.ValidGrant}})
	}

	if r.RedirectURI != cmp.RedirectURI {
		out = append(out, &revisions.Change{Key: "redirectURI", Old: []any{cmp.RedirectURI}, New: []any{r.RedirectURI}})
	}

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.Trusted != cmp.Trusted {
		out = append(out, &revisions.Change{Key: "trusted", Old: []any{cmp.Trusted}, New: []any{r.Trusted}})
	}

	if !reflect.DeepEqual(r.ValidFrom, cmp.ValidFrom) {
		out = append(out, &revisions.Change{Key: "validFrom", Old: []any{cmp.ValidFrom}, New: []any{r.ValidFrom}})
	}

	if !reflect.DeepEqual(r.ExpiresAt, cmp.ExpiresAt) {
		out = append(out, &revisions.Change{Key: "expiresAt", Old: []any{cmp.ExpiresAt}, New: []any{r.ExpiresAt}})
	}

	if (r.Security == nil) != (cmp.Security == nil) {
		out = append(out, &revisions.Change{Key: "security", Old: []any{cmp.Security}, New: []any{r.Security}})
	} else if r.Security != nil {
		for _, c := range r.Security.Diff(cmp.Security) {
			c.Key = "security." + c.Key
			out = append(out, c)
		}
	}

	if r.OwnedBy != cmp.OwnedBy {
		out = append(out, &revisions.Change{Key: "ownedBy", Old: []any{cmp.OwnedBy}, New: []any{r.OwnedBy}})
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

func (r AuthClientMeta) Clone() *AuthClientMeta {
	dup := r
	return &dup
}

func (r AuthClientMeta) Diff(cmp *AuthClientMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AuthClientMeta{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *AuthClientMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AuthClientMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AuthClientSecurity) Clone() *AuthClientSecurity {
	dup := r
	if r.PermittedRoles != nil {
		dup.PermittedRoles = make([]string, len(r.PermittedRoles))
		copy(dup.PermittedRoles, r.PermittedRoles)
	}

	if r.ProhibitedRoles != nil {
		dup.ProhibitedRoles = make([]string, len(r.ProhibitedRoles))
		copy(dup.ProhibitedRoles, r.ProhibitedRoles)
	}

	if r.ForcedRoles != nil {
		dup.ForcedRoles = make([]string, len(r.ForcedRoles))
		copy(dup.ForcedRoles, r.ForcedRoles)
	}

	return &dup
}

func (r AuthClientSecurity) Diff(cmp *AuthClientSecurity) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AuthClientSecurity{}
	}
	if r.ImpersonateUser != cmp.ImpersonateUser {
		out = append(out, &revisions.Change{Key: "impersonateUser", Old: []any{cmp.ImpersonateUser}, New: []any{r.ImpersonateUser}})
	}

	if r.UserGroup != cmp.UserGroup {
		out = append(out, &revisions.Change{Key: "userGroup", Old: []any{cmp.UserGroup}, New: []any{r.UserGroup}})
	}

	if !reflect.DeepEqual(r.PermittedRoles, cmp.PermittedRoles) {
		out = append(out, &revisions.Change{Key: "permittedRoles", Old: []any{cmp.PermittedRoles}, New: []any{r.PermittedRoles}})
	}

	if !reflect.DeepEqual(r.ProhibitedRoles, cmp.ProhibitedRoles) {
		out = append(out, &revisions.Change{Key: "prohibitedRoles", Old: []any{cmp.ProhibitedRoles}, New: []any{r.ProhibitedRoles}})
	}

	if !reflect.DeepEqual(r.ForcedRoles, cmp.ForcedRoles) {
		out = append(out, &revisions.Change{Key: "forcedRoles", Old: []any{cmp.ForcedRoles}, New: []any{r.ForcedRoles}})
	}

	return out
}

func (r *AuthClientSecurity) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AuthClientSecurity) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseAuthClientMeta(ss []string) (p *AuthClientMeta, err error) {
	p = &AuthClientMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func ParseAuthClientSecurity(ss []string) (p *AuthClientSecurity, err error) {
	p = &AuthClientSecurity{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}
