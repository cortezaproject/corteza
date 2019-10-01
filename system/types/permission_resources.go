package types

import (
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

const SystemPermissionResource = permissions.Resource("system")
const ApplicationPermissionResource = permissions.Resource("system:application:")
const OrganisationPermissionResource = permissions.Resource("system:organisation:")
const UserPermissionResource = permissions.Resource("system:user:")
const RolePermissionResource = permissions.Resource("system:role:")
const AutomationScriptPermissionResource = permissions.Resource("system:automation-script:")
const AutomationTriggerPermissionResource = permissions.Resource("system:automation-trigger:")
const ReminderPermissionResource = permissions.Resource("system:reminder:")
