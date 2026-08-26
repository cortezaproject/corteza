package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ProjectMemberFilter struct {
		ProjectMemberID []string          `json:"projectMemberID"`
		ProjectID       uint64            `json:"projectID,string,omitempty"`
		UserID          uint64            `json:"userID,string,omitempty"`
		RolePreset      ProjectMemberRole `json:"rolePreset,omitempty"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*ProjectMember) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

const (
	ProjectRoleGovernanceOwner ProjectMemberRole = "governance-owner"
	ProjectRoleSecurityOwner   ProjectMemberRole = "security-owner"
	ProjectRoleDeveloper       ProjectMemberRole = "developer"
	ProjectRoleJuniorDeveloper ProjectMemberRole = "junior-developer"
	ProjectRoleMember          ProjectMemberRole = "member"
	// ProjectRoleExecutiveAuthority reads the whole project and grants (or
	// requests changes on) the publish approval; it never edits.
	ProjectRoleExecutiveAuthority ProjectMemberRole = "executive-authority"
	// ProjectRoleInfrastructureAdministrator maintains platform infrastructure;
	// it carries no project-content access at all.
	ProjectRoleInfrastructureAdministrator ProjectMemberRole = "infrastructure-administrator"
)

// Valid reports whether r is one of the known presets.
func (r ProjectMemberRole) Valid() bool {
	switch r {
	case ProjectRoleGovernanceOwner,
		ProjectRoleSecurityOwner,
		ProjectRoleDeveloper,
		ProjectRoleJuniorDeveloper,
		ProjectRoleMember,
		ProjectRoleExecutiveAuthority,
		ProjectRoleInfrastructureAdministrator:
		return true
	}
	return false
}

// Capabilities resolves the runtime capabilities for the preset.
func (r ProjectMemberRole) Capabilities() ProjectCapabilities {
	switch r {
	case ProjectRoleGovernanceOwner, ProjectRoleSecurityOwner:
		return ProjectCapabilities{CanRead: true, CanWrite: true, CanRequestApproval: true, CanGrantApproval: true}
	case ProjectRoleDeveloper:
		return ProjectCapabilities{CanRead: true, CanWrite: true, CanRequestApproval: true}
	case ProjectRoleJuniorDeveloper:
		return ProjectCapabilities{CanRead: true, CanWrite: true}
	case ProjectRoleExecutiveAuthority:
		return ProjectCapabilities{CanRead: true, CanGrantApproval: true}
	case ProjectRoleInfrastructureAdministrator:
		return ProjectCapabilities{}
	default:
		return ProjectCapabilities{CanRead: true}
	}
}
