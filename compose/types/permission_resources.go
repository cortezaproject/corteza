package types

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const ComposeRBACResource = rbac.Resource("compose")
const NamespaceRBACResource = rbac.Resource("compose:namespace:")
const ChartRBACResource = rbac.Resource("compose:chart:")
const ModuleRBACResource = rbac.Resource("compose:module:")
const ModuleFieldRBACResource = rbac.Resource("compose:module-field:")
const PageRBACResource = rbac.Resource("compose:page:")
