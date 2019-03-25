package types

import (
	"github.com/crusttech/crust/internal/rules"
)

const PermissionResource = rules.Resource("system")
const ApplicationPermissionResource = rules.Resource("system:application:")
const OrganisationPermissionResource = rules.Resource("system:organisation:")
const UserPermissionResource = rules.Resource("system:user:")
const RolePermissionResource = rules.Resource("system:role:")
