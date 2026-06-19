package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	// Project is the resource struct; it is generated into project.gen.go.

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

		// Mode selects the build pipeline: gated runs the full governance
		// pipeline with approval gates, free only the build steps. Immutable
		// after creation.
		Mode ProjectMode `json:"mode,omitempty"`

		// NamespaceID points to the compose namespace auto-created for this
		// project; project-scoped compose resources (modules, pages, charts)
		// live there.
		NamespaceID uint64 `json:"namespaceID,string,omitempty"`

		// DeployerCategories are the AI Act deployer-category answers gathered
		// at creation; any true value makes a Fundamental Rights Impact
		// Assessment (FRIA) step required.
		DeployerCategories ProjectDeployerCategories `json:"deployerCategories,omitempty"`
		// FriaRequired is derived from DeployerCategories at creation and
		// stored so the pipeline shape doesn't silently change later.
		FriaRequired bool `json:"friaRequired,omitempty"`

		// ResourceManagement holds the AI/infrastructure governance forms and
		// the whitelist of permitted connections the Connections step picks
		// from.
		ResourceManagement ProjectResourceManagement `json:"resourceManagement,omitempty"`
	}

	// ProjectDeployerCategories are the AI Act Deployer-category answers from
	// project creation (EU AI Act Art. 27).
	ProjectDeployerCategories struct {
		PublicAuthorityAnnex3    bool `json:"publicAuthorityAnnex3,omitempty"`
		PrivateEssentialServices bool `json:"privateEssentialServices,omitempty"`
		InsuranceBanking         bool `json:"insuranceBanking,omitempty"`
	}

	// ProjectResourceManagement is the Resource Management step state: the
	// AI/infra form values plus the permitted-connection whitelist.
	ProjectResourceManagement struct {
		AI          map[string]any                `json:"ai,omitempty"`
		Infra       map[string]any                `json:"infra,omitempty"`
		Connections []*ProjectPermittedConnection `json:"connections,omitempty"`
	}

	// ProjectPermittedConnection is one whitelisted external connection the
	// project may use; the Connections build step instantiates from this list.
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
	ProjectMode       string
)

const (
	// ProjectStatusDraft: being built in the wizard, not yet provisioned.
	ProjectStatusDraft ProjectStatus = "draft"
	// ProjectStatusActive: operational (provisioned, pre-publish).
	ProjectStatusActive ProjectStatus = "active"
	// ProjectStatusPublished: configuration locked as a version; changes go
	// through the change-request procedure.
	ProjectStatusPublished ProjectStatus = "published"
	ProjectStatusArchived  ProjectStatus = "archived"
	ProjectStatusSuspended ProjectStatus = "suspended"

	ProjectVisibilityOpen       ProjectVisibility = "open"
	ProjectVisibilityInviteOnly ProjectVisibility = "invite-only"

	// ProjectModeFree runs only the build steps; ProjectModeGated runs the full
	// governance pipeline with approval gates.
	ProjectModeFree  ProjectMode = "free"
	ProjectModeGated ProjectMode = "gated"
)

// Valid reports whether m is a known build mode.
func (m ProjectMode) Valid() bool {
	return m == ProjectModeFree || m == ProjectModeGated
}
