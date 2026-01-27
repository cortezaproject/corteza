package envoyx

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza/server/automation/types"
)

// AutomationWorkflowRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func AutomationWorkflowRbacReferences(workflow string) (res *Ref, pp []*Ref, err error) {
	if workflow != "*" {
		res = &Ref{ResourceType: types.WorkflowResourceType, Identifiers: MakeIdentifiers(workflow)}
	}

	return
}

// AutomationNgAutomationRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func AutomationNgAutomationRbacReferences(ngAutomation string) (res *Ref, pp []*Ref, err error) {
	if ngAutomation != "*" {
		res = &Ref{ResourceType: types.NgAutomationResourceType, Identifiers: MakeIdentifiers(ngAutomation)}
	}

	return
}
