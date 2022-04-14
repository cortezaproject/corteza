package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/system/types"
)

// SystemAttachmentRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemAttachmentRbacReferences(attachment string) (res *Ref, pp []*Ref, err error) {
	if attachment != "*" {
		res = &Ref{ResourceType: types.AttachmentResourceType, Identifiers: MakeIdentifiers(attachment)}
	}

	return
}

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

// SystemApigwFilterRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemApigwFilterRbacReferences(apigwFilter string) (res *Ref, pp []*Ref, err error) {
	if apigwFilter != "*" {
		res = &Ref{ResourceType: types.ApigwFilterResourceType, Identifiers: MakeIdentifiers(apigwFilter)}
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

// SystemAuthConfirmedClientRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemAuthConfirmedClientRbacReferences(authConfirmedClient string) (res *Ref, pp []*Ref, err error) {
	if authConfirmedClient != "*" {
		res = &Ref{ResourceType: types.AuthConfirmedClientResourceType, Identifiers: MakeIdentifiers(authConfirmedClient)}
	}

	return
}

// SystemAuthSessionRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemAuthSessionRbacReferences(authSession string) (res *Ref, pp []*Ref, err error) {
	if authSession != "*" {
		res = &Ref{ResourceType: types.AuthSessionResourceType, Identifiers: MakeIdentifiers(authSession)}
	}

	return
}

// SystemAuthOa2tokenRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemAuthOa2tokenRbacReferences(authOa2token string) (res *Ref, pp []*Ref, err error) {
	if authOa2token != "*" {
		res = &Ref{ResourceType: types.AuthOa2tokenResourceType, Identifiers: MakeIdentifiers(authOa2token)}
	}

	return
}

// SystemCredentialRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemCredentialRbacReferences(credential string) (res *Ref, pp []*Ref, err error) {
	if credential != "*" {
		res = &Ref{ResourceType: types.CredentialResourceType, Identifiers: MakeIdentifiers(credential)}
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

// SystemQueueMessageRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemQueueMessageRbacReferences(queueMessage string) (res *Ref, pp []*Ref, err error) {
	if queueMessage != "*" {
		res = &Ref{ResourceType: types.QueueMessageResourceType, Identifiers: MakeIdentifiers(queueMessage)}
	}

	return
}

// SystemReminderRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemReminderRbacReferences(reminder string) (res *Ref, pp []*Ref, err error) {
	if reminder != "*" {
		res = &Ref{ResourceType: types.ReminderResourceType, Identifiers: MakeIdentifiers(reminder)}
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

// SystemResourceTranslationRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemResourceTranslationRbacReferences(resourceTranslation string) (res *Ref, pp []*Ref, err error) {
	if resourceTranslation != "*" {
		res = &Ref{ResourceType: types.ResourceTranslationResourceType, Identifiers: MakeIdentifiers(resourceTranslation)}
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

// SystemRoleMemberRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemRoleMemberRbacReferences(roleMember string) (res *Ref, pp []*Ref, err error) {
	if roleMember != "*" {
		res = &Ref{ResourceType: types.RoleMemberResourceType, Identifiers: MakeIdentifiers(roleMember)}
	}

	return
}

// SystemSettingValueRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemSettingValueRbacReferences(settingValue string) (res *Ref, pp []*Ref, err error) {
	if settingValue != "*" {
		res = &Ref{ResourceType: types.SettingValueResourceType, Identifiers: MakeIdentifiers(settingValue)}
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
