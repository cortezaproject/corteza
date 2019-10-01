package types

import (
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

const ComposePermissionResource = permissions.Resource("compose")
const NamespacePermissionResource = permissions.Resource("compose:namespace:")
const ChartPermissionResource = permissions.Resource("compose:chart:")
const ModulePermissionResource = permissions.Resource("compose:module:")
const ModuleFieldPermissionResource = permissions.Resource("compose:module-field:")
const PagePermissionResource = permissions.Resource("compose:page:")
const AutomationScriptPermissionResource = permissions.Resource("compose:automation-script:")
const AutomationTriggerPermissionResource = permissions.Resource("compose:automation-trigger:")
