package types

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const AutomationRBACResource = rbac.Resource("automation")
const WorkflowRBACResource = rbac.Resource("automation:workflow:")
