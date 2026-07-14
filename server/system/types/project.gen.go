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
	Project struct {
		ID               uint64                           `json:"projectID,string"`
		TenantID         uint64                           `json:"tenantID,string,omitempty"`
		Handle           string                           `json:"handle"`
		Status           ProjectStatus                    `json:"status"`
		Config           ProjectConfig                    `json:"config"`
		Mode             ProjectMode                      `json:"mode"`
		Meta             ProjectMeta                      `json:"meta"`
		Governance       ProjectGovernance                `json:"governance"`
		ProjectID        uint64                           `json:"rootProjectID,string,omitempty"`
		ParentRevisionID uint64                           `json:"parentRevisionID,string,omitempty"`
		Revision         int                              `json:"revision,omitempty"`
		CreatedAt        time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt        *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt        *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy        uint64                           `json:"createdBy,string"`
		UpdatedBy        uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy        uint64                           `json:"deletedBy,string,omitempty"`
		Labels           map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	ProjectConfig struct {
		Visibility         ProjectVisibility         `json:"visibility,omitempty"`
		DefaultMemberRole  ProjectMemberRole         `json:"defaultMemberRole,omitempty"`
		FeatureFlags       map[string]bool           `json:"featureFlags,omitempty"`
		NamespaceID        uint64                    `json:"namespaceID,string,omitempty"`
		DeployerCategories ProjectDeployerCategories `json:"deployerCategories,omitempty"`
		FriaRequired       bool                      `json:"friaRequired,omitempty"`
		ResourceManagement ProjectResourceManagement `json:"resourceManagement,omitempty"`
	}

	ProjectDeployerCategories struct {
		PublicAuthorityAnnex3    bool `json:"publicAuthorityAnnex3,omitempty"`
		PrivateEssentialServices bool `json:"privateEssentialServices,omitempty"`
		InsuranceBanking         bool `json:"insuranceBanking,omitempty"`
	}

	ProjectResourceManagement struct {
		AI          map[string]any                `json:"ai,omitempty"`
		Infra       map[string]any                `json:"infra,omitempty"`
		Connections []*ProjectPermittedConnection `json:"connections,omitempty"`
	}

	ProjectMeta struct {
		Short       string   `json:"short,omitempty"`
		Description string   `json:"description,omitempty"`
		Icon        string   `json:"icon,omitempty"`
		Color       string   `json:"color,omitempty"`
		Tags        []string `json:"tags,omitempty"`
	}

	ProjectGovernanceStep struct {
		Values     map[string]any          `json:"values,omitempty"`
		Status     ProjectGovernanceStatus `json:"status"`
		ReviewNote string                  `json:"reviewNote,omitempty"`
	}

	ProjectPermittedConnection struct {
		ID                  string `json:"id"`
		Name                string `json:"name"`
		Connector           string `json:"connector,omitempty"`
		Type                string `json:"type,omitempty"`
		Description         string `json:"description,omitempty"`
		ActionIfUnavailable string `json:"actionIfUnavailable,omitempty"`
		Replacement         string `json:"replacement,omitempty"`
		IsAiSystem          string `json:"isAiSystem,omitempty"`
	}

	ProjectStatus string

	ProjectVisibility string

	ProjectMode string

	ProjectMemberRole string

	ProjectGovernanceStatus string
)

func (r Project) Clone() *Project {
	dup := r
	dup.Config = *r.Config.Clone()

	dup.Meta = *r.Meta.Clone()

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

func (r Project) Diff(cmp *Project) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Project{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
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

	if !reflect.DeepEqual(r.Mode, cmp.Mode) {
		out = append(out, &revisions.Change{Key: "mode", Old: []any{cmp.Mode}, New: []any{r.Mode}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Governance, cmp.Governance) {
		out = append(out, &revisions.Change{Key: "governance", Old: []any{cmp.Governance}, New: []any{r.Governance}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "rootProjectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.ParentRevisionID != cmp.ParentRevisionID {
		out = append(out, &revisions.Change{Key: "parentRevisionID", Old: []any{cmp.ParentRevisionID}, New: []any{r.ParentRevisionID}})
	}

	if r.Revision != cmp.Revision {
		out = append(out, &revisions.Change{Key: "revision", Old: []any{cmp.Revision}, New: []any{r.Revision}})
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

func (r ProjectConfig) Clone() *ProjectConfig {
	dup := r
	if r.FeatureFlags != nil {
		dup.FeatureFlags = make(map[string]bool, len(r.FeatureFlags))
		for k, v := range r.FeatureFlags {
			dup.FeatureFlags[k] = v
		}
	}

	dup.DeployerCategories = *r.DeployerCategories.Clone()

	dup.ResourceManagement = *r.ResourceManagement.Clone()

	return &dup
}

func (r ProjectConfig) Diff(cmp *ProjectConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectConfig{}
	}
	if !reflect.DeepEqual(r.Visibility, cmp.Visibility) {
		out = append(out, &revisions.Change{Key: "visibility", Old: []any{cmp.Visibility}, New: []any{r.Visibility}})
	}

	if !reflect.DeepEqual(r.DefaultMemberRole, cmp.DefaultMemberRole) {
		out = append(out, &revisions.Change{Key: "defaultMemberRole", Old: []any{cmp.DefaultMemberRole}, New: []any{r.DefaultMemberRole}})
	}

	if !reflect.DeepEqual(r.FeatureFlags, cmp.FeatureFlags) {
		out = append(out, &revisions.Change{Key: "featureFlags", Old: []any{cmp.FeatureFlags}, New: []any{r.FeatureFlags}})
	}

	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	for _, c := range r.DeployerCategories.Diff(&cmp.DeployerCategories) {
		c.Key = "deployerCategories." + c.Key
		out = append(out, c)
	}

	if r.FriaRequired != cmp.FriaRequired {
		out = append(out, &revisions.Change{Key: "friaRequired", Old: []any{cmp.FriaRequired}, New: []any{r.FriaRequired}})
	}

	for _, c := range r.ResourceManagement.Diff(&cmp.ResourceManagement) {
		c.Key = "resourceManagement." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *ProjectConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ProjectConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ProjectDeployerCategories) Clone() *ProjectDeployerCategories {
	dup := r
	return &dup
}

func (r ProjectDeployerCategories) Diff(cmp *ProjectDeployerCategories) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectDeployerCategories{}
	}
	if r.PublicAuthorityAnnex3 != cmp.PublicAuthorityAnnex3 {
		out = append(out, &revisions.Change{Key: "publicAuthorityAnnex3", Old: []any{cmp.PublicAuthorityAnnex3}, New: []any{r.PublicAuthorityAnnex3}})
	}

	if r.PrivateEssentialServices != cmp.PrivateEssentialServices {
		out = append(out, &revisions.Change{Key: "privateEssentialServices", Old: []any{cmp.PrivateEssentialServices}, New: []any{r.PrivateEssentialServices}})
	}

	if r.InsuranceBanking != cmp.InsuranceBanking {
		out = append(out, &revisions.Change{Key: "insuranceBanking", Old: []any{cmp.InsuranceBanking}, New: []any{r.InsuranceBanking}})
	}

	return out
}

func (r *ProjectDeployerCategories) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ProjectDeployerCategories) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ProjectResourceManagement) Clone() *ProjectResourceManagement {
	dup := r
	if r.AI != nil {
		dup.AI = make(map[string]any, len(r.AI))
		for k, v := range r.AI {
			dup.AI[k] = v
		}
	}

	if r.Infra != nil {
		dup.Infra = make(map[string]any, len(r.Infra))
		for k, v := range r.Infra {
			dup.Infra[k] = v
		}
	}

	if r.Connections != nil {
		dup.Connections = make([]*ProjectPermittedConnection, len(r.Connections))
		copy(dup.Connections, r.Connections)
	}

	return &dup
}

func (r ProjectResourceManagement) Diff(cmp *ProjectResourceManagement) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectResourceManagement{}
	}
	if !reflect.DeepEqual(r.AI, cmp.AI) {
		out = append(out, &revisions.Change{Key: "ai", Old: []any{cmp.AI}, New: []any{r.AI}})
	}

	if !reflect.DeepEqual(r.Infra, cmp.Infra) {
		out = append(out, &revisions.Change{Key: "infra", Old: []any{cmp.Infra}, New: []any{r.Infra}})
	}

	if !reflect.DeepEqual(r.Connections, cmp.Connections) {
		out = append(out, &revisions.Change{Key: "connections", Old: []any{cmp.Connections}, New: []any{r.Connections}})
	}

	return out
}

func (r *ProjectResourceManagement) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ProjectResourceManagement) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ProjectMeta) Clone() *ProjectMeta {
	dup := r
	if r.Tags != nil {
		dup.Tags = make([]string, len(r.Tags))
		copy(dup.Tags, r.Tags)
	}

	return &dup
}

func (r ProjectMeta) Diff(cmp *ProjectMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.Icon != cmp.Icon {
		out = append(out, &revisions.Change{Key: "icon", Old: []any{cmp.Icon}, New: []any{r.Icon}})
	}

	if r.Color != cmp.Color {
		out = append(out, &revisions.Change{Key: "color", Old: []any{cmp.Color}, New: []any{r.Color}})
	}

	if !reflect.DeepEqual(r.Tags, cmp.Tags) {
		out = append(out, &revisions.Change{Key: "tags", Old: []any{cmp.Tags}, New: []any{r.Tags}})
	}

	return out
}

func (r *ProjectMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ProjectMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ProjectGovernanceStep) Clone() *ProjectGovernanceStep {
	dup := r
	if r.Values != nil {
		dup.Values = make(map[string]any, len(r.Values))
		for k, v := range r.Values {
			dup.Values[k] = v
		}
	}

	return &dup
}

func (r ProjectGovernanceStep) Diff(cmp *ProjectGovernanceStep) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectGovernanceStep{}
	}
	if !reflect.DeepEqual(r.Values, cmp.Values) {
		out = append(out, &revisions.Change{Key: "values", Old: []any{cmp.Values}, New: []any{r.Values}})
	}

	if !reflect.DeepEqual(r.Status, cmp.Status) {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.ReviewNote != cmp.ReviewNote {
		out = append(out, &revisions.Change{Key: "reviewNote", Old: []any{cmp.ReviewNote}, New: []any{r.ReviewNote}})
	}

	return out
}

func (r *ProjectGovernanceStep) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ProjectGovernanceStep) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ProjectPermittedConnection) Clone() *ProjectPermittedConnection {
	dup := r
	return &dup
}

func (r ProjectPermittedConnection) Diff(cmp *ProjectPermittedConnection) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectPermittedConnection{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "id", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Connector != cmp.Connector {
		out = append(out, &revisions.Change{Key: "connector", Old: []any{cmp.Connector}, New: []any{r.Connector}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.ActionIfUnavailable != cmp.ActionIfUnavailable {
		out = append(out, &revisions.Change{Key: "actionIfUnavailable", Old: []any{cmp.ActionIfUnavailable}, New: []any{r.ActionIfUnavailable}})
	}

	if r.Replacement != cmp.Replacement {
		out = append(out, &revisions.Change{Key: "replacement", Old: []any{cmp.Replacement}, New: []any{r.Replacement}})
	}

	if r.IsAiSystem != cmp.IsAiSystem {
		out = append(out, &revisions.Change{Key: "isAiSystem", Old: []any{cmp.IsAiSystem}, New: []any{r.IsAiSystem}})
	}

	return out
}

func (r *ProjectPermittedConnection) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ProjectPermittedConnection) Value() (driver.Value, error) { return json.Marshal(r) }

const (
	ProjectStatusDraft      ProjectStatus = "draft"
	ProjectStatusActive     ProjectStatus = "active"
	ProjectStatusPublished  ProjectStatus = "published"
	ProjectStatusArchived   ProjectStatus = "archived"
	ProjectStatusSuspended  ProjectStatus = "suspended"
	ProjectStatusDeprecated ProjectStatus = "deprecated"
)

const (
	ProjectVisibilityOpen       ProjectVisibility = "open"
	ProjectVisibilityInviteOnly ProjectVisibility = "invite-only"
)

const (
	ProjectModeFree  ProjectMode = "free"
	ProjectModeGated ProjectMode = "gated"
)

const (
	ProjectMemberRoleProjectRoleGovernanceOwner             ProjectMemberRole = "governance-owner"
	ProjectMemberRoleProjectRoleSecurityOwner               ProjectMemberRole = "security-owner"
	ProjectMemberRoleProjectRoleDeveloper                   ProjectMemberRole = "developer"
	ProjectMemberRoleProjectRoleJuniorDeveloper             ProjectMemberRole = "junior-developer"
	ProjectMemberRoleProjectRoleMember                      ProjectMemberRole = "member"
	ProjectMemberRoleProjectRoleExecutiveAuthority          ProjectMemberRole = "executive-authority"
	ProjectMemberRoleProjectRoleInfrastructureAdministrator ProjectMemberRole = "infrastructure-administrator"
)

const (
	ProjectGovernanceStatusDraft            ProjectGovernanceStatus = "draft"
	ProjectGovernanceStatusSubmitted        ProjectGovernanceStatus = "submitted"
	ProjectGovernanceStatusApproved         ProjectGovernanceStatus = "approved"
	ProjectGovernanceStatusChangesRequested ProjectGovernanceStatus = "changes-requested"
)

func ParseProjectConfig(ss []string) (p ProjectConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseProjectMeta(ss []string) (p ProjectMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ProjectGovernance) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ProjectGovernance) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseProjectGovernance(ss []string) (p ProjectGovernance, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
