package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - automation.workflow.yaml
// - automation.yaml

import (
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the automation component
	//
	// This struct is auto-generated
	Component struct{}
)

const (
	WorkflowRbacResourceSchema  = "corteza+automation.workflow"
	ComponentRbacResourceSchema = "corteza+automation"
)

// RbacResource returns string representation of RBAC resource for Workflow by calling WorkflowRbacResource fn
//
// RBAC resource is in the corteza+automation.workflow:/... format
//
// This function is auto-generated
func (r Workflow) RbacResource() string {
	return WorkflowRbacResource(r.ID)
}

// WorkflowRbacResource returns string representation of RBAC resource for Workflow
//
// RBAC resource is in the corteza+automation.workflow:/... format
//
// This function is auto-generated
func WorkflowRbacResource(ID uint64) string {
	out := WorkflowRbacResourceSchema + ":"
	out += "/"

	if ID != 0 {
		out += strconv.FormatUint(ID, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza+automation:/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza+automation:/... format
//
// This function is auto-generated
func ComponentRbacResource() string {
	out := ComponentRbacResourceSchema + ":"
	return out
}
