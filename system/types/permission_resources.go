package types

import (
	"github.com/cortezaproject/corteza-server/internal/permissions"
)

const SystemPermissionResource = permissions.Resource("system")
const ApplicationPermissionResource = permissions.Resource("system:application:")
const OrganisationPermissionResource = permissions.Resource("system:organisation:")
const UserPermissionResource = permissions.Resource("system:user:")
const RolePermissionResource = permissions.Resource("system:role:")
