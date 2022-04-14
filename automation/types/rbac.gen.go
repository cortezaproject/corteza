package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the automation component
	//
	// This struct is auto-generated
	Component struct{}
)

var (
	_ = fmt.Printf
	_ = strconv.FormatUint
)

const (
	WorkflowResourceType  = "corteza::automation:workflow"
	SessionResourceType   = "corteza::automation:session"
	TriggerResourceType   = "corteza::automation:trigger"
	ComponentResourceType = "corteza::automation"
)

// RbacResource returns string representation of RBAC resource for Workflow by calling WorkflowRbacResource fn
//
// RBAC resource is in the corteza::automation:workflow/... format
//
// This function is auto-generated
func (r Workflow) RbacResource() string {
	return WorkflowRbacResource(r.ID)
}

// WorkflowRbacResource returns string representation of RBAC resource for Workflow
//
// RBAC resource is in the corteza::automation:workflow/... format
//
// This function is auto-generated
func WorkflowRbacResource(id uint64) string {
	cpts := []interface{}{WorkflowResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(WorkflowRbacResourceTpl(), cpts...)

}

func WorkflowRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Session by calling SessionRbacResource fn
//
// RBAC resource is in the corteza::automation:session/... format
//
// This function is auto-generated
func (r Session) RbacResource() string {
	return SessionRbacResource(r.ID)
}

// SessionRbacResource returns string representation of RBAC resource for Session
//
// RBAC resource is in the corteza::automation:session/... format
//
// This function is auto-generated
func SessionRbacResource(id uint64) string {
	cpts := []interface{}{SessionResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(SessionRbacResourceTpl(), cpts...)

}

func SessionRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Trigger by calling TriggerRbacResource fn
//
// RBAC resource is in the corteza::automation:trigger/... format
//
// This function is auto-generated
func (r Trigger) RbacResource() string {
	return TriggerRbacResource(r.ID)
}

// TriggerRbacResource returns string representation of RBAC resource for Trigger
//
// RBAC resource is in the corteza::automation:trigger/... format
//
// This function is auto-generated
func TriggerRbacResource(id uint64) string {
	cpts := []interface{}{TriggerResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(TriggerRbacResourceTpl(), cpts...)

}

func TriggerRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::automation/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::automation/ format
//
// This function is auto-generated
func ComponentRbacResource() string {
	return ComponentResourceType + "/"

}

func ComponentRbacResourceTpl() string {
	return "%s"
}
