package types

import (
	"github.com/cortezaproject/corteza-server/internal/permissions"
)

const ComposePermissionResource = permissions.Resource("compose")
const NamespacePermissionResource = permissions.Resource("compose:namespace:")
const ChartPermissionResource = permissions.Resource("compose:chart:")
const ModulePermissionResource = permissions.Resource("compose:module:")
const ModuleFieldPermissionResource = permissions.Resource("compose:module-field:")
const PagePermissionResource = permissions.Resource("compose:page:")
const TriggerPermissionResource = permissions.Resource("compose:trigger:")
