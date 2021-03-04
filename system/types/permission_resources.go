package types

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const SystemRBACResource = rbac.Resource("system")
const ApplicationRBACResource = rbac.Resource("system:application:")
const TemplateRBACResource = rbac.Resource("system:template:")
const UserRBACResource = rbac.Resource("system:user:")
const RoleRBACResource = rbac.Resource("system:role:")
const AuthClientRBACResource = rbac.Resource("system:auth-client:")
