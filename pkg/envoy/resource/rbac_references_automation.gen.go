package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/automation/types"
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

// AutomationSessionRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func AutomationSessionRbacReferences(session string) (res *Ref, pp []*Ref, err error) {
	if session != "*" {
		res = &Ref{ResourceType: types.SessionResourceType, Identifiers: MakeIdentifiers(session)}
	}

	return
}

// AutomationTriggerRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func AutomationTriggerRbacReferences(trigger string) (res *Ref, pp []*Ref, err error) {
	if trigger != "*" {
		res = &Ref{ResourceType: types.TriggerResourceType, Identifiers: MakeIdentifiers(trigger)}
	}

	return
}
