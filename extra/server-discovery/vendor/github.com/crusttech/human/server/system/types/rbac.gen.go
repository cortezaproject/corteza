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
	// Component struct serves as a virtual resource type for the system component
	//
	// This struct is auto-generated
	Component struct{}
)

var (
	_ = fmt.Printf
	_ = strconv.FormatUint
)

// RbacResource returns string representation of RBAC resource for Application by calling ApplicationRbacResource fn
//
// RBAC resource is in the corteza::system:application/... format
//
// This function is auto-generated
func (r Application) RbacResource() string {
	return ApplicationRbacResource(r.ID)
}

// ApplicationRbacResource returns string representation of RBAC resource for Application
//
// RBAC resource is in the corteza::system:application/... format
//
// This function is auto-generated
func ApplicationRbacResource(id uint64) string {
	cpts := []interface{}{ApplicationResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ApplicationRbacResourceTpl(), cpts...)

}

func ApplicationRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ApigwRoute by calling ApigwRouteRbacResource fn
//
// RBAC resource is in the corteza::system:apigw-route/... format
//
// This function is auto-generated
func (r ApigwRoute) RbacResource() string {
	return ApigwRouteRbacResource(r.ID)
}

// ApigwRouteRbacResource returns string representation of RBAC resource for ApigwRoute
//
// RBAC resource is in the corteza::system:apigw-route/... format
//
// This function is auto-generated
func ApigwRouteRbacResource(id uint64) string {
	cpts := []interface{}{ApigwRouteResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ApigwRouteRbacResourceTpl(), cpts...)

}

func ApigwRouteRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for AuthClient by calling AuthClientRbacResource fn
//
// RBAC resource is in the corteza::system:auth-client/... format
//
// This function is auto-generated
func (r AuthClient) RbacResource() string {
	return AuthClientRbacResource(r.ID)
}

// AuthClientRbacResource returns string representation of RBAC resource for AuthClient
//
// RBAC resource is in the corteza::system:auth-client/... format
//
// This function is auto-generated
func AuthClientRbacResource(id uint64) string {
	cpts := []interface{}{AuthClientResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(AuthClientRbacResourceTpl(), cpts...)

}

func AuthClientRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for DataPrivacyRequest by calling DataPrivacyRequestRbacResource fn
//
// RBAC resource is in the corteza::system:data-privacy-request/... format
//
// This function is auto-generated
func (r DataPrivacyRequest) RbacResource() string {
	return DataPrivacyRequestRbacResource(r.ID)
}

// DataPrivacyRequestRbacResource returns string representation of RBAC resource for DataPrivacyRequest
//
// RBAC resource is in the corteza::system:data-privacy-request/... format
//
// This function is auto-generated
func DataPrivacyRequestRbacResource(id uint64) string {
	cpts := []interface{}{DataPrivacyRequestResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(DataPrivacyRequestRbacResourceTpl(), cpts...)

}

func DataPrivacyRequestRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Queue by calling QueueRbacResource fn
//
// RBAC resource is in the corteza::system:queue/... format
//
// This function is auto-generated
func (r Queue) RbacResource() string {
	return QueueRbacResource(r.ID)
}

// QueueRbacResource returns string representation of RBAC resource for Queue
//
// RBAC resource is in the corteza::system:queue/... format
//
// This function is auto-generated
func QueueRbacResource(id uint64) string {
	cpts := []interface{}{QueueResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(QueueRbacResourceTpl(), cpts...)

}

func QueueRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Report by calling ReportRbacResource fn
//
// RBAC resource is in the corteza::system:report/... format
//
// This function is auto-generated
func (r Report) RbacResource() string {
	return ReportRbacResource(r.ID)
}

// ReportRbacResource returns string representation of RBAC resource for Report
//
// RBAC resource is in the corteza::system:report/... format
//
// This function is auto-generated
func ReportRbacResource(id uint64) string {
	cpts := []interface{}{ReportResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ReportRbacResourceTpl(), cpts...)

}

func ReportRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Role by calling RoleRbacResource fn
//
// RBAC resource is in the corteza::system:role/... format
//
// This function is auto-generated
func (r Role) RbacResource() string {
	return RoleRbacResource(r.ID)
}

// RoleRbacResource returns string representation of RBAC resource for Role
//
// RBAC resource is in the corteza::system:role/... format
//
// This function is auto-generated
func RoleRbacResource(id uint64) string {
	cpts := []interface{}{RoleResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(RoleRbacResourceTpl(), cpts...)

}

func RoleRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for UserGroup by calling UserGroupRbacResource fn
//
// RBAC resource is in the corteza::system:user-group/... format
//
// This function is auto-generated
func (r UserGroup) RbacResource() string {
	return UserGroupRbacResource(r.ID)
}

// UserGroupRbacResource returns string representation of RBAC resource for UserGroup
//
// RBAC resource is in the corteza::system:user-group/... format
//
// This function is auto-generated
func UserGroupRbacResource(id uint64) string {
	cpts := []interface{}{UserGroupResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(UserGroupRbacResourceTpl(), cpts...)

}

func UserGroupRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Template by calling TemplateRbacResource fn
//
// RBAC resource is in the corteza::system:template/... format
//
// This function is auto-generated
func (r Template) RbacResource() string {
	return TemplateRbacResource(r.ID)
}

// TemplateRbacResource returns string representation of RBAC resource for Template
//
// RBAC resource is in the corteza::system:template/... format
//
// This function is auto-generated
func TemplateRbacResource(id uint64) string {
	cpts := []interface{}{TemplateResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(TemplateRbacResourceTpl(), cpts...)

}

func TemplateRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for User by calling UserRbacResource fn
//
// RBAC resource is in the corteza::system:user/... format
//
// This function is auto-generated
func (r User) RbacResource() string {
	return UserRbacResource(r.ID)
}

// UserRbacResource returns string representation of RBAC resource for User
//
// RBAC resource is in the corteza::system:user/... format
//
// This function is auto-generated
func UserRbacResource(id uint64) string {
	cpts := []interface{}{UserResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(UserRbacResourceTpl(), cpts...)

}

func UserRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for DalConnection by calling DalConnectionRbacResource fn
//
// RBAC resource is in the corteza::system:dal-connection/... format
//
// This function is auto-generated
func (r DalConnection) RbacResource() string {
	return DalConnectionRbacResource(r.ID)
}

// DalConnectionRbacResource returns string representation of RBAC resource for DalConnection
//
// RBAC resource is in the corteza::system:dal-connection/... format
//
// This function is auto-generated
func DalConnectionRbacResource(id uint64) string {
	cpts := []interface{}{DalConnectionResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(DalConnectionRbacResourceTpl(), cpts...)

}

func DalConnectionRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Connection by calling ConnectionRbacResource fn
//
// RBAC resource is in the corteza::system:connection/... format
//
// This function is auto-generated
func (r Connection) RbacResource() string {
	return ConnectionRbacResource(r.ID)
}

// ConnectionRbacResource returns string representation of RBAC resource for Connection
//
// RBAC resource is in the corteza::system:connection/... format
//
// This function is auto-generated
func ConnectionRbacResource(id uint64) string {
	cpts := []interface{}{ConnectionResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ConnectionRbacResourceTpl(), cpts...)

}

func ConnectionRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ConfiguredConnection by calling ConfiguredConnectionRbacResource fn
//
// RBAC resource is in the corteza::system:configured-connection/... format
//
// This function is auto-generated
func (r ConfiguredConnection) RbacResource() string {
	return ConfiguredConnectionRbacResource(r.ID)
}

// ConfiguredConnectionRbacResource returns string representation of RBAC resource for ConfiguredConnection
//
// RBAC resource is in the corteza::system:configured-connection/... format
//
// This function is auto-generated
func ConfiguredConnectionRbacResource(id uint64) string {
	cpts := []interface{}{ConfiguredConnectionResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ConfiguredConnectionRbacResourceTpl(), cpts...)

}

func ConfiguredConnectionRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for LlmProvider by calling LlmProviderRbacResource fn
//
// RBAC resource is in the corteza::system:llm-provider/... format
//
// This function is auto-generated
func (r LlmProvider) RbacResource() string {
	return LlmProviderRbacResource(r.ID)
}

// LlmProviderRbacResource returns string representation of RBAC resource for LlmProvider
//
// RBAC resource is in the corteza::system:llm-provider/... format
//
// This function is auto-generated
func LlmProviderRbacResource(id uint64) string {
	cpts := []interface{}{LlmProviderResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(LlmProviderRbacResourceTpl(), cpts...)

}

func LlmProviderRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Agent by calling AgentRbacResource fn
//
// RBAC resource is in the corteza::system:agent/... format
//
// This function is auto-generated
func (r Agent) RbacResource() string {
	return AgentRbacResource(r.ID)
}

// AgentRbacResource returns string representation of RBAC resource for Agent
//
// RBAC resource is in the corteza::system:agent/... format
//
// This function is auto-generated
func AgentRbacResource(id uint64) string {
	cpts := []interface{}{AgentResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(AgentRbacResourceTpl(), cpts...)

}

func AgentRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for AiConversation by calling AiConversationRbacResource fn
//
// RBAC resource is in the corteza::system:ai-conversation/... format
//
// This function is auto-generated
func (r AiConversation) RbacResource() string {
	return AiConversationRbacResource(r.ID)
}

// AiConversationRbacResource returns string representation of RBAC resource for AiConversation
//
// RBAC resource is in the corteza::system:ai-conversation/... format
//
// This function is auto-generated
func AiConversationRbacResource(id uint64) string {
	cpts := []interface{}{AiConversationResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(AiConversationRbacResourceTpl(), cpts...)

}

func AiConversationRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for KnowledgeBase by calling KnowledgeBaseRbacResource fn
//
// RBAC resource is in the corteza::system:knowledge-base/... format
//
// This function is auto-generated
func (r KnowledgeBase) RbacResource() string {
	return KnowledgeBaseRbacResource(r.ID)
}

// KnowledgeBaseRbacResource returns string representation of RBAC resource for KnowledgeBase
//
// RBAC resource is in the corteza::system:knowledge-base/... format
//
// This function is auto-generated
func KnowledgeBaseRbacResource(id uint64) string {
	cpts := []interface{}{KnowledgeBaseResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(KnowledgeBaseRbacResourceTpl(), cpts...)

}

func KnowledgeBaseRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Chatbot by calling ChatbotRbacResource fn
//
// RBAC resource is in the corteza::system:chatbot/... format
//
// This function is auto-generated
func (r Chatbot) RbacResource() string {
	return ChatbotRbacResource(r.ID)
}

// ChatbotRbacResource returns string representation of RBAC resource for Chatbot
//
// RBAC resource is in the corteza::system:chatbot/... format
//
// This function is auto-generated
func ChatbotRbacResource(id uint64) string {
	cpts := []interface{}{ChatbotResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ChatbotRbacResourceTpl(), cpts...)

}

func ChatbotRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ChatbotSession by calling ChatbotSessionRbacResource fn
//
// RBAC resource is in the corteza::system:chatbot-session/... format
//
// This function is auto-generated
func (r ChatbotSession) RbacResource() string {
	return ChatbotSessionRbacResource(r.ID)
}

// ChatbotSessionRbacResource returns string representation of RBAC resource for ChatbotSession
//
// RBAC resource is in the corteza::system:chatbot-session/... format
//
// This function is auto-generated
func ChatbotSessionRbacResource(id uint64) string {
	cpts := []interface{}{ChatbotSessionResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ChatbotSessionRbacResourceTpl(), cpts...)

}

func ChatbotSessionRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ChatbotSessionStep by calling ChatbotSessionStepRbacResource fn
//
// RBAC resource is in the corteza::system:chatbot-session-step/... format
//
// This function is auto-generated
func (r ChatbotSessionStep) RbacResource() string {
	return ChatbotSessionStepRbacResource(r.ID)
}

// ChatbotSessionStepRbacResource returns string representation of RBAC resource for ChatbotSessionStep
//
// RBAC resource is in the corteza::system:chatbot-session-step/... format
//
// This function is auto-generated
func ChatbotSessionStepRbacResource(id uint64) string {
	cpts := []interface{}{ChatbotSessionStepResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ChatbotSessionStepRbacResourceTpl(), cpts...)

}

func ChatbotSessionStepRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ChatbotSessionHandoff by calling ChatbotSessionHandoffRbacResource fn
//
// RBAC resource is in the corteza::system:chatbot-session-handoff/... format
//
// This function is auto-generated
func (r ChatbotSessionHandoff) RbacResource() string {
	return ChatbotSessionHandoffRbacResource(r.ID)
}

// ChatbotSessionHandoffRbacResource returns string representation of RBAC resource for ChatbotSessionHandoff
//
// RBAC resource is in the corteza::system:chatbot-session-handoff/... format
//
// This function is auto-generated
func ChatbotSessionHandoffRbacResource(id uint64) string {
	cpts := []interface{}{ChatbotSessionHandoffResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ChatbotSessionHandoffRbacResourceTpl(), cpts...)

}

func ChatbotSessionHandoffRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Tenant by calling TenantRbacResource fn
//
// RBAC resource is in the corteza::system:tenant/... format
//
// This function is auto-generated
func (r Tenant) RbacResource() string {
	return TenantRbacResource(r.ID)
}

// TenantRbacResource returns string representation of RBAC resource for Tenant
//
// RBAC resource is in the corteza::system:tenant/... format
//
// This function is auto-generated
func TenantRbacResource(id uint64) string {
	cpts := []interface{}{TenantResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(TenantRbacResourceTpl(), cpts...)

}

func TenantRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Project by calling ProjectRbacResource fn
//
// RBAC resource is in the corteza::system:project/... format
//
// This function is auto-generated
func (r Project) RbacResource() string {
	return ProjectRbacResource(r.ID)
}

// ProjectRbacResource returns string representation of RBAC resource for Project
//
// RBAC resource is in the corteza::system:project/... format
//
// This function is auto-generated
func ProjectRbacResource(id uint64) string {
	cpts := []interface{}{ProjectResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectRbacResourceTpl(), cpts...)

}

func ProjectRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ProjectAiSystem by calling ProjectAiSystemRbacResource fn
//
// RBAC resource is in the corteza::system:project-ai-system/... format
//
// This function is auto-generated
func (r ProjectAiSystem) RbacResource() string {
	return ProjectAiSystemRbacResource(r.ID)
}

// ProjectAiSystemRbacResource returns string representation of RBAC resource for ProjectAiSystem
//
// RBAC resource is in the corteza::system:project-ai-system/... format
//
// This function is auto-generated
func ProjectAiSystemRbacResource(id uint64) string {
	cpts := []interface{}{ProjectAiSystemResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectAiSystemRbacResourceTpl(), cpts...)

}

func ProjectAiSystemRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ProjectFriaScenario by calling ProjectFriaScenarioRbacResource fn
//
// RBAC resource is in the corteza::system:project-fria-scenario/... format
//
// This function is auto-generated
func (r ProjectFriaScenario) RbacResource() string {
	return ProjectFriaScenarioRbacResource(r.ID)
}

// ProjectFriaScenarioRbacResource returns string representation of RBAC resource for ProjectFriaScenario
//
// RBAC resource is in the corteza::system:project-fria-scenario/... format
//
// This function is auto-generated
func ProjectFriaScenarioRbacResource(id uint64) string {
	cpts := []interface{}{ProjectFriaScenarioResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectFriaScenarioRbacResourceTpl(), cpts...)

}

func ProjectFriaScenarioRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ProjectIncident by calling ProjectIncidentRbacResource fn
//
// RBAC resource is in the corteza::system:project-incident/... format
//
// This function is auto-generated
func (r ProjectIncident) RbacResource() string {
	return ProjectIncidentRbacResource(r.ID)
}

// ProjectIncidentRbacResource returns string representation of RBAC resource for ProjectIncident
//
// RBAC resource is in the corteza::system:project-incident/... format
//
// This function is auto-generated
func ProjectIncidentRbacResource(id uint64) string {
	cpts := []interface{}{ProjectIncidentResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectIncidentRbacResourceTpl(), cpts...)

}

func ProjectIncidentRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ProjectFeature by calling ProjectFeatureRbacResource fn
//
// RBAC resource is in the corteza::system:project-feature/... format
//
// This function is auto-generated
func (r ProjectFeature) RbacResource() string {
	return ProjectFeatureRbacResource(r.ID)
}

// ProjectFeatureRbacResource returns string representation of RBAC resource for ProjectFeature
//
// RBAC resource is in the corteza::system:project-feature/... format
//
// This function is auto-generated
func ProjectFeatureRbacResource(id uint64) string {
	cpts := []interface{}{ProjectFeatureResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectFeatureRbacResourceTpl(), cpts...)

}

func ProjectFeatureRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ProjectPrivacy by calling ProjectPrivacyRbacResource fn
//
// RBAC resource is in the corteza::system:project-privacy/... format
//
// This function is auto-generated
func (r ProjectPrivacy) RbacResource() string {
	return ProjectPrivacyRbacResource(r.ID)
}

// ProjectPrivacyRbacResource returns string representation of RBAC resource for ProjectPrivacy
//
// RBAC resource is in the corteza::system:project-privacy/... format
//
// This function is auto-generated
func ProjectPrivacyRbacResource(id uint64) string {
	cpts := []interface{}{ProjectPrivacyResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectPrivacyRbacResourceTpl(), cpts...)

}

func ProjectPrivacyRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ProjectTask by calling ProjectTaskRbacResource fn
//
// RBAC resource is in the corteza::system:project-task/... format
//
// This function is auto-generated
func (r ProjectTask) RbacResource() string {
	return ProjectTaskRbacResource(r.ID)
}

// ProjectTaskRbacResource returns string representation of RBAC resource for ProjectTask
//
// RBAC resource is in the corteza::system:project-task/... format
//
// This function is auto-generated
func ProjectTaskRbacResource(id uint64) string {
	cpts := []interface{}{ProjectTaskResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectTaskRbacResourceTpl(), cpts...)

}

func ProjectTaskRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ProjectReview by calling ProjectReviewRbacResource fn
//
// RBAC resource is in the corteza::system:project-review/... format
//
// This function is auto-generated
func (r ProjectReview) RbacResource() string {
	return ProjectReviewRbacResource(r.ID)
}

// ProjectReviewRbacResource returns string representation of RBAC resource for ProjectReview
//
// RBAC resource is in the corteza::system:project-review/... format
//
// This function is auto-generated
func ProjectReviewRbacResource(id uint64) string {
	cpts := []interface{}{ProjectReviewResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectReviewRbacResourceTpl(), cpts...)

}

func ProjectReviewRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ProjectBacklogItem by calling ProjectBacklogItemRbacResource fn
//
// RBAC resource is in the corteza::system:project-backlog-item/... format
//
// This function is auto-generated
func (r ProjectBacklogItem) RbacResource() string {
	return ProjectBacklogItemRbacResource(r.ID)
}

// ProjectBacklogItemRbacResource returns string representation of RBAC resource for ProjectBacklogItem
//
// RBAC resource is in the corteza::system:project-backlog-item/... format
//
// This function is auto-generated
func ProjectBacklogItemRbacResource(id uint64) string {
	cpts := []interface{}{ProjectBacklogItemResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ProjectBacklogItemRbacResourceTpl(), cpts...)

}

func ProjectBacklogItemRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::system/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::system/ format
//
// This function is auto-generated
func ComponentRbacResource() string {
	return ComponentResourceType + "/"

}

func ComponentRbacResourceTpl() string {
	return "%s"
}
