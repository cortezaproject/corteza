package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
)

type (
	// Project is a secondary isolation unit that belongs to exactly one tenant.
	//
	// TenantID is carried for the multi-tenancy model but is not yet wired into
	// scope routing; it defaults to 0 until tenant scoping lands.
	Project struct {
		ID       uint64        `json:"projectID,string"`
		TenantID uint64        `json:"tenantID,string,omitempty"`
		Handle   string        `json:"handle"`
		Status   ProjectStatus `json:"status"`

		Config ProjectConfig `json:"config"`
		Meta   ProjectMeta   `json:"meta"`

		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	ProjectFilter struct {
		ProjectID []string      `json:"projectID"`
		TenantID  uint64        `json:"tenantID,string,omitempty"`
		Handle    string        `json:"handle"`
		Status    ProjectStatus `json:"status"`
		Query     string        `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*Project) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	// ProjectConfig holds behavioral settings. System acts on these.
	ProjectConfig struct {
		// Visibility controls whether all tenant members can access the project
		// (open) or only explicitly added members (invite-only).
		Visibility ProjectVisibility `json:"visibility,omitempty"`
		// DefaultMemberRole is the preset assigned to new members automatically.
		DefaultMemberRole ProjectMemberRole `json:"defaultMemberRole,omitempty"`
		// FeatureFlags are project-specific feature overrides.
		FeatureFlags map[string]bool `json:"featureFlags,omitempty"`
	}

	// ProjectMeta is display-only. System does not act on these.
	ProjectMeta struct {
		Short       string   `json:"short,omitempty"`
		Description string   `json:"description,omitempty"`
		Icon        string   `json:"icon,omitempty"`
		Color       string   `json:"color,omitempty"`
		Tags        []string `json:"tags,omitempty"`
	}

	ProjectStatus     string
	ProjectVisibility string
)

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusArchived  ProjectStatus = "archived"
	ProjectStatusSuspended ProjectStatus = "suspended"

	ProjectVisibilityOpen       ProjectVisibility = "open"
	ProjectVisibilityInviteOnly ProjectVisibility = "invite-only"
)

func (m *ProjectConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ProjectConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *ProjectMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ProjectMeta) Value() (driver.Value, error) { return json.Marshal(m) }

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
