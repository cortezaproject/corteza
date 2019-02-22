package service

import (
	"sync"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	o                   sync.Once
	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultRole         RoleService
	DefaultPermission   PermissionsService
	DefaultOrganisation OrganisationService
)

func Init() {
	o.Do(func() {
		DefaultAuth = Auth()
		DefaultUser = User()
		DefaultRole = Role()
		DefaultPermission = Permission()
		DefaultOrganisation = Organisation()
	})
}
