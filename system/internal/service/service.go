package service

import (
	"context"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultRole         RoleService
	DefaultRules        RulesService
	DefaultOrganisation OrganisationService
	DefaultApplication  ApplicationService
	DefaultPermissions  PermissionsService
)

func Init() error {
	ctx := context.Background()
	DefaultRules = Rules(ctx)
	DefaultPermissions = Permissions(ctx)
	DefaultAuth = Auth(ctx)
	DefaultUser = User(ctx)
	DefaultRole = Role(ctx)
	DefaultOrganisation = Organisation(ctx)
	DefaultApplication = Application(ctx)
	return nil
}
