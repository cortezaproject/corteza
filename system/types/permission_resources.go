package types

import (
	"github.com/crusttech/crust/internal/rules"
)

const PermissionResource = rules.Resource("system")
const ApplicationPermissionResource = rules.Resource("system:application:")
const OrganisationPermissionResource = rules.Resource("system:organisation:")
const RolePermissionResource = rules.Resource("system:role:")
