package types

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const SystemRBACResource = rbac.Resource("system")
const ApplicationRBACResource = rbac.Resource("system:application:")
const OrganisationRBACResource = rbac.Resource("system:organisation:")
const UserRBACResource = rbac.Resource("system:user:")
const RoleRBACResource = rbac.Resource("system:role:")
const WorkflowRBACResource = rbac.Resource("system:workflow:")
