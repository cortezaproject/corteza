package types

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const FederationRBACResource = rbac.Resource("federation")
const NodeRBACResource = rbac.Resource("federation:node:")
const ModuleRBACResource = rbac.Resource("federation:module:")
