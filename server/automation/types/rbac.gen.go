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

// RbacResource returns string representation of RBAC resource for NgAutomation by calling NgAutomationRbacResource fn
//
// RBAC resource is in the corteza::automation:ng-automation/... format
//
// This function is auto-generated
func (r NgAutomation) RbacResource() string {
	return NgAutomationRbacResource(r.ID)
}

// NgAutomationRbacResource returns string representation of RBAC resource for NgAutomation
//
// RBAC resource is in the corteza::automation:ng-automation/... format
//
// This function is auto-generated
func NgAutomationRbacResource(id uint64) string {
	cpts := []interface{}{NgAutomationResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(NgAutomationRbacResourceTpl(), cpts...)

}

func NgAutomationRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for TriggerDefinition by calling TriggerDefinitionRbacResource fn
//
// RBAC resource is in the corteza::automation:trigger-definition/... format
//
// This function is auto-generated
func (r TriggerDefinition) RbacResource() string {
	return TriggerDefinitionRbacResource(r.ID)
}

// TriggerDefinitionRbacResource returns string representation of RBAC resource for TriggerDefinition
//
// RBAC resource is in the corteza::automation:trigger-definition/... format
//
// This function is auto-generated
func TriggerDefinitionRbacResource(id uint64) string {
	cpts := []interface{}{TriggerDefinitionResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(TriggerDefinitionRbacResourceTpl(), cpts...)

}

func TriggerDefinitionRbacResourceTpl() string {
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
