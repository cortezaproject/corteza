package auth

import (
	"github.com/cortezaproject/corteza-server/system/types"
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
			federationUser.SetRoles([]uint64{bpr.ID})
		case ServiceUserHandle:
			serviceUser = u.Clone()
			federationUser.SetRoles([]uint64{bpr.ID})
		case FederationUserHandle:
			federationUser = u.Clone()
			federationUser.SetRoles([]uint64{bpr.ID})
		}
	}
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
