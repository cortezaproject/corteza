package auth

import (
	"github.com/cortezaproject/corteza/server/system/types"
)

const (
	ProvisionUserHandle  = "corteza-provisioner"
	ServiceUserHandle    = "corteza-service"
	FederationUserHandle = "corteza-federation"

	BypassRoleHandle        = "super-admin"
	AuthenticatedRoleHandle = "authenticated"
	AnonymousRoleHandle     = "anonymous"
)

var (
	provisionUser  *types.User
	serviceUser    *types.User
	federationUser *types.User

	authenticatedRoles []*types.Role
	anonymousRoles     []*types.Role
	bypassRoles        []*types.Role
)

// SetSystemUsers takes list of users and sets/updates appropriate
// provision, service & fed users
//
// These are then accessed via special *User() fn.
func SetSystemUsers(uu types.UserSet, rr types.RoleSet) {
	bpr := rr.FindByHandle(BypassRoleHandle)

	for _, u := range uu {
		switch u.Handle {
		case ProvisionUserHandle:
			provisionUser = u.Clone()
			provisionUser.SetRoles(bpr.ID)
		case ServiceUserHandle:
			serviceUser = u.Clone()
			serviceUser.SetRoles(bpr.ID)
		case FederationUserHandle:
			federationUser = u.Clone()
			federationUser.SetRoles(bpr.ID)
		}
	}
}

func SetSystemRoles(rr types.RoleSet) {
	for _, r := range rr {
		switch r.Handle {
		case BypassRoleHandle:
			bypassRoles = append(bypassRoles, r.Clone())
		case AuthenticatedRoleHandle:
			authenticatedRoles = append(authenticatedRoles, r.Clone())
		case AnonymousRoleHandle:
			anonymousRoles = append(anonymousRoles, r.Clone())
		}
	}
}

// BypassRoles returns all bypass Roles
func BypassRoles() types.RoleSet {
	return bypassRoles
}

// AuthenticatedRoles returns all authenticated Roles
func AuthenticatedRoles() types.RoleSet {
	return authenticatedRoles
}

// AnonymousRoles returns all anonymous Roles
func AnonymousRoles() types.RoleSet {
	return anonymousRoles
}

// ProvisionUser returns clone of system provision user
func ProvisionUser() *types.User {
	return provisionUser.Clone()
}

// ServiceUser returns clone of system service user
func ServiceUser() *types.User {
	return serviceUser.Clone()
}

// FederationUser returns clone of system federation user
func FederationUser() *types.User {
	return federationUser.Clone()
}
