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

	// ProjectMemberRole is a fixed named preset (not freeform). Stored by value
	// on the membership record. Capabilities are derived at runtime from the
	// preset definition — never stored per-member.
	ProjectMemberRole string

	// ProjectCapabilities is derived at runtime from a member's RolePreset.
	// It is never persisted.
	ProjectCapabilities struct {
		CanRead            bool `json:"canRead"`
		CanWrite           bool `json:"canWrite"`
		CanRequestApproval bool `json:"canRequestApproval"`
		CanGrantApproval   bool `json:"canGrantApproval"`
	}
)

const (
	ProjectRoleGovernanceOwner ProjectMemberRole = "governance-owner"
	ProjectRoleSecurityOwner   ProjectMemberRole = "security-owner"
	ProjectRoleDeveloper       ProjectMemberRole = "developer"
	ProjectRoleJuniorDeveloper ProjectMemberRole = "junior-developer"
	ProjectRoleMember          ProjectMemberRole = "member"
	// ProjectRoleExecutiveAuthority reads the whole project and signs off at
	// the first governance gate; it never edits.
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
//
//	preset                 read write request grant
//	governance-owner        ✓    ✓     ✓      ✓
//	security-owner          ✓    ✓     ✓      ✓
//	developer               ✓    ✓     ✓      ✗
//	junior-developer        ✓    ✓     ✗      ✗
//	executive-authority     ✓    ✗     ✗      ✓
//	infrastructure-admin.   ✗    ✗     ✗      ✗
//	member (fallback)       ✓    ✗     ✗      ✗
//
// Unknown presets fall back to member.
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
	default: // member + unknown
		return ProjectCapabilities{CanRead: true}
	}
}
