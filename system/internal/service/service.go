package service

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
	DefaultRules = Rules()
	DefaultPermissions = Permissions()
	DefaultAuth = Auth()
	DefaultUser = User()
	DefaultRole = Role()
	DefaultOrganisation = Organisation()
	DefaultApplication = Application()
	return nil
}
