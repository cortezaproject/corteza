package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/system/types"
)

// SystemApplicationRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemApplicationRbacReferences(application string) (res *Ref, pp []*Ref, err error) {
	if application != "*" {
		res = &Ref{ResourceType: types.ApplicationResourceType, Identifiers: MakeIdentifiers(application)}
	}

	return
}

// SystemApigwRouteRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemApigwRouteRbacReferences(apigwRoute string) (res *Ref, pp []*Ref, err error) {
	if apigwRoute != "*" {
		res = &Ref{ResourceType: types.ApigwRouteResourceType, Identifiers: MakeIdentifiers(apigwRoute)}
	}

	return
}

// SystemAuthClientRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemAuthClientRbacReferences(authClient string) (res *Ref, pp []*Ref, err error) {
	if authClient != "*" {
		res = &Ref{ResourceType: types.AuthClientResourceType, Identifiers: MakeIdentifiers(authClient)}
	}

	return
}

// SystemDataPrivacyRequestRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemDataPrivacyRequestRbacReferences(dataPrivacyRequest string) (res *Ref, pp []*Ref, err error) {
	if dataPrivacyRequest != "*" {
		res = &Ref{ResourceType: types.DataPrivacyRequestResourceType, Identifiers: MakeIdentifiers(dataPrivacyRequest)}
	}

	return
}

// SystemQueueRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemQueueRbacReferences(queue string) (res *Ref, pp []*Ref, err error) {
	if queue != "*" {
		res = &Ref{ResourceType: types.QueueResourceType, Identifiers: MakeIdentifiers(queue)}
	}

	return
}

// SystemReportRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemReportRbacReferences(report string) (res *Ref, pp []*Ref, err error) {
	if report != "*" {
		res = &Ref{ResourceType: types.ReportResourceType, Identifiers: MakeIdentifiers(report)}
	}

	return
}

// SystemRoleRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemRoleRbacReferences(role string) (res *Ref, pp []*Ref, err error) {
	if role != "*" {
		res = &Ref{ResourceType: types.RoleResourceType, Identifiers: MakeIdentifiers(role)}
	}

	return
}

// SystemTemplateRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemTemplateRbacReferences(template string) (res *Ref, pp []*Ref, err error) {
	if template != "*" {
		res = &Ref{ResourceType: types.TemplateResourceType, Identifiers: MakeIdentifiers(template)}
	}

	return
}

// SystemUserRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemUserRbacReferences(user string) (res *Ref, pp []*Ref, err error) {
	if user != "*" {
		res = &Ref{ResourceType: types.UserResourceType, Identifiers: MakeIdentifiers(user)}
	}

	return
}

// SystemDalConnectionRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemDalConnectionRbacReferences(dalConnection string) (res *Ref, pp []*Ref, err error) {
	if dalConnection != "*" {
		res = &Ref{ResourceType: types.DalConnectionResourceType, Identifiers: MakeIdentifiers(dalConnection)}
	}

	return
}
