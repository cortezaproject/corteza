package types

import (
	"github.com/crusttech/crust/internal/rules"
)

const PermissionResource = rules.Resource("compose")
const NamespacePermissionResource = rules.Resource("compose:namespace:")
const ChartPermissionResource = rules.Resource("compose:chart:")
const ModulePermissionResource = rules.Resource("compose:module:")
const PagePermissionResource = rules.Resource("compose:page:")
const TriggerPermissionResource = rules.Resource("compose:trigger:")
