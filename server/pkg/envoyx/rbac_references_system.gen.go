package envoyx

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/system/types"
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

// SystemUserGroupRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemUserGroupRbacReferences(userGroup string) (res *Ref, pp []*Ref, err error) {
	if userGroup != "*" {
		res = &Ref{ResourceType: types.UserGroupResourceType, Identifiers: MakeIdentifiers(userGroup)}
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

// SystemConnectionRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemConnectionRbacReferences(connection string) (res *Ref, pp []*Ref, err error) {
	if connection != "*" {
		res = &Ref{ResourceType: types.ConnectionResourceType, Identifiers: MakeIdentifiers(connection)}
	}

	return
}

// SystemConfiguredConnectionRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemConfiguredConnectionRbacReferences(configuredConnection string) (res *Ref, pp []*Ref, err error) {
	if configuredConnection != "*" {
		res = &Ref{ResourceType: types.ConfiguredConnectionResourceType, Identifiers: MakeIdentifiers(configuredConnection)}
	}

	return
}

// SystemLlmProviderRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemLlmProviderRbacReferences(llmProvider string) (res *Ref, pp []*Ref, err error) {
	if llmProvider != "*" {
		res = &Ref{ResourceType: types.LlmProviderResourceType, Identifiers: MakeIdentifiers(llmProvider)}
	}

	return
}

// SystemAgentRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemAgentRbacReferences(agent string) (res *Ref, pp []*Ref, err error) {
	if agent != "*" {
		res = &Ref{ResourceType: types.AgentResourceType, Identifiers: MakeIdentifiers(agent)}
	}

	return
}

// SystemAiConversationRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemAiConversationRbacReferences(aiConversation string) (res *Ref, pp []*Ref, err error) {
	if aiConversation != "*" {
		res = &Ref{ResourceType: types.AiConversationResourceType, Identifiers: MakeIdentifiers(aiConversation)}
	}

	return
}

// SystemKnowledgeBaseRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemKnowledgeBaseRbacReferences(knowledgeBase string) (res *Ref, pp []*Ref, err error) {
	if knowledgeBase != "*" {
		res = &Ref{ResourceType: types.KnowledgeBaseResourceType, Identifiers: MakeIdentifiers(knowledgeBase)}
	}

	return
}

// SystemChatbotRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemChatbotRbacReferences(chatbot string) (res *Ref, pp []*Ref, err error) {
	if chatbot != "*" {
		res = &Ref{ResourceType: types.ChatbotResourceType, Identifiers: MakeIdentifiers(chatbot)}
	}

	return
}

// SystemChatbotSessionRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemChatbotSessionRbacReferences(chatbotSession string) (res *Ref, pp []*Ref, err error) {
	if chatbotSession != "*" {
		res = &Ref{ResourceType: types.ChatbotSessionResourceType, Identifiers: MakeIdentifiers(chatbotSession)}
	}

	return
}

// SystemChatbotSessionStepRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemChatbotSessionStepRbacReferences(chatbotSessionStep string) (res *Ref, pp []*Ref, err error) {
	if chatbotSessionStep != "*" {
		res = &Ref{ResourceType: types.ChatbotSessionStepResourceType, Identifiers: MakeIdentifiers(chatbotSessionStep)}
	}

	return
}

// SystemChatbotSessionHandoffRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemChatbotSessionHandoffRbacReferences(chatbotSessionHandoff string) (res *Ref, pp []*Ref, err error) {
	if chatbotSessionHandoff != "*" {
		res = &Ref{ResourceType: types.ChatbotSessionHandoffResourceType, Identifiers: MakeIdentifiers(chatbotSessionHandoff)}
	}

	return
}

// SystemTenantRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemTenantRbacReferences(tenant string) (res *Ref, pp []*Ref, err error) {
	if tenant != "*" {
		res = &Ref{ResourceType: types.TenantResourceType, Identifiers: MakeIdentifiers(tenant)}
	}

	return
}

// SystemProjectRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemProjectRbacReferences(project string) (res *Ref, pp []*Ref, err error) {
	if project != "*" {
		res = &Ref{ResourceType: types.ProjectResourceType, Identifiers: MakeIdentifiers(project)}
	}

	return
}

// SystemProjectGroupRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemProjectGroupRbacReferences(projectGroup string) (res *Ref, pp []*Ref, err error) {
	if projectGroup != "*" {
		res = &Ref{ResourceType: types.ProjectGroupResourceType, Identifiers: MakeIdentifiers(projectGroup)}
	}

	return
}

// SystemProjectIncidentRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemProjectIncidentRbacReferences(projectIncident string) (res *Ref, pp []*Ref, err error) {
	if projectIncident != "*" {
		res = &Ref{ResourceType: types.ProjectIncidentResourceType, Identifiers: MakeIdentifiers(projectIncident)}
	}

	return
}

// SystemProjectFeatureRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemProjectFeatureRbacReferences(projectFeature string) (res *Ref, pp []*Ref, err error) {
	if projectFeature != "*" {
		res = &Ref{ResourceType: types.ProjectFeatureResourceType, Identifiers: MakeIdentifiers(projectFeature)}
	}

	return
}

// SystemProjectPrivacyRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemProjectPrivacyRbacReferences(projectPrivacy string) (res *Ref, pp []*Ref, err error) {
	if projectPrivacy != "*" {
		res = &Ref{ResourceType: types.ProjectPrivacyResourceType, Identifiers: MakeIdentifiers(projectPrivacy)}
	}

	return
}

// SystemProjectTaskRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemProjectTaskRbacReferences(projectTask string) (res *Ref, pp []*Ref, err error) {
	if projectTask != "*" {
		res = &Ref{ResourceType: types.ProjectTaskResourceType, Identifiers: MakeIdentifiers(projectTask)}
	}

	return
}

// SystemProjectReviewRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemProjectReviewRbacReferences(projectReview string) (res *Ref, pp []*Ref, err error) {
	if projectReview != "*" {
		res = &Ref{ResourceType: types.ProjectReviewResourceType, Identifiers: MakeIdentifiers(projectReview)}
	}

	return
}
