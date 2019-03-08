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
	DefaultPermissions  PermissionsService
	DefaultOrganisation OrganisationService
	DefaultApplication  ApplicationService
)

func Init() {
	o.Do(func() {
		DefaultAuth = Auth()
		DefaultUser = User()
		DefaultRole = Role()
		DefaultPermissions = Permissions()
		DefaultOrganisation = Organisation()
		DefaultApplication = Application()
	})
}
