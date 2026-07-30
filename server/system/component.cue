package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

component: schema.#component & {
	handle: "system"

	resources: {
		"attachment":                   attachment
		"application":                  application
		"apigw-route":                  apigw_route
		"apigw-filter":                 apigw_filter
		"auth-client":                  auth_client
		"auth-confirmed-client":        auth_confirmed_client
		"auth-session":                 auth_session
		"auth-oa2token":                auth_oa2token
		"credential":                   credential
		"data-privacy-request":         data_privacy_request
		"data-privacy-request-comment": data_privacy_request_comment
		"queue":                        queue
		"queue-message":                queue_message
		"reminder":                     reminder
		"notification":                 notification
		"report":                       report
		"resource-translation":         resource_translation
		"role":                         role
		"role-member":                  role_member
		"user-group":                   user_group
		"settings":                     settings
		"template":                     template
		"user":                         user
		"dal-connection":               dal_connection
		"dal-sensitivity-level":        dal_sensitivity_level
		"dal-schema-alteration":        dal_schema_alteration
		"connection":                   connection
		"configured-connection":        configured_connection
		"llm-provider":                 llm_provider
		"agent":                        agent
		"ai-conversation":              ai_conversation
		"knowledge-base":               knowledge_base
		"chatbot":                      chatbot
		"chatbot-session":              chatbot_session
		"chatbot-session-step":         chatbot_session_step
		"chatbot-session-handoff":      chatbot_session_handoff
		"tenant":                       tenant
		"tenant-membership":            tenant_membership
		"project":                      project
		"project-member":               project_member
		"project-ai-system":            project_ai_system
		"project-ai-system-entry":      project_ai_system_entry
		"project-fria-scenario":        project_fria_scenario
		"project-incident":             project_incident
		"project-feature":              project_feature
		"project-privacy":              project_privacy
		"project-task":                 project_task
		"project-review":               project_review
		"project-backlog-item":         project_backlog_item
		"dml-connection":               dml_connection
		"dml-mapping":                  dml_mapping
		"dml-import-run":               dml_import_run
	}

	types: {
		"apigw-profiler-hit": {noIdField: true}
		"apigw-profiler-aggregation": {noIdField: true}
		"privacy-dal-connection": {}
		"dml-column-map": {noIdField: true}
	}

	rbac: operations: {
		"action-log.read": description: "Access to action log"

		"settings.read": description:       "Read system settings"
		"settings.manage": description:     "Manage system settings"
		"auth-client.create": description:  "Create auth clients"
		"auth-clients.search": description: "List, search or filter auth clients"

		"role.create": description:  "Create roles"
		"roles.search": description: "List, search or filter roles"

		"user-group.create": description:  "Create user groups"
		"user-groups.search": description: "List, search or filter user groups"

		"user.create": description:  "Create users"
		"users.search": description: "List, search or filter users"

		"dal-connection.create": description:  "Create DAL connections"
		"dal-connections.search": description: "List, search or filter DAL connections"

		"dal-sensitivity-level.manage": description: "Can manage DAL sensitivity levels"

		"application.create": description:      "Create applications"
		"applications.search": description:     "List, search or filter auth clients"
		"application.flag.self": description:   "Manage private flags for applications"
		"application.flag.global": description: "Manage global flags for applications"

		"template.create": description:  "Create template"
		"templates.search": description: "List, search or filter templates"

		"report.create": description:  "Create report"
		"reports.search": description: "List, search or filter reports"

		"reminder.assign": description: " Assign reminders"

		"queue.create": description:  "Create messagebus queues"
		"queues.search": description: "List, search or filter messagebus queues"

		"apigw-route.create": description:  "Create API gateway route"
		"apigw-routes.search": description: "List search or filter API gateway routes"

		"resource-translations.manage": description: "List, search, create, or update resource translations"

		"dal-schema-alterations.manage": description: "List, search, apply, or dismiss DAL alterations"

		"connection.create": description:             "Create connections"
		"connections.search": description:            "List, search or filter connections"
		"configured-connection.create": description:  "Install connection connections"
		"configured-connections.search": description: "List, search or filter connection connections"

		"data-privacy-request.create": description:  "Create data privacy requests"
		"data-privacy-requests.search": description: "List, search or filter data privacy requests"

		"notification.assign": description: "Assign notifications to other users"

		"llm-provider.create": description:  "Create LLM providers"
		"llm-providers.search": description: "List, search or filter LLM providers"

		"agent.create": description:  "Create agents"
		"agents.search": description: "List, search or filter agents"

		"project-incident.create": description:  "Create project incidents"
		"project-incidents.search": description: "List, search or filter project incidents"

		"project-feature.create": description:  "Create project features"
		"project-features.search": description: "List, search or filter project features"

		"project-privacy.create": description:  "Create project privacy items"
		"project-privacys.search": description: "List, search or filter project privacy items"

		"project-task.create": description:  "Create project tasks"
		"project-tasks.search": description: "List, search or filter project tasks"

		"project-review.create": description:  "Create project reviews"
		"project-reviews.search": description: "List, search or filter project reviews"

		"project-backlog-item.create": description:  "Create project backlog items"
		"project-backlog-items.search": description: "List, search or filter project backlog items"

		"ai-conversation.create": description:  "Create AI conversations"
		"ai-conversations.search": description: "List, search or filter AI conversations"

		"knowledge-base.create": description:  "Create knowledge bases"
		"knowledge-bases.search": description: "List, search or filter knowledge bases"

		"chatbot.create": description:  "Create chatbots"
		"chatbots.search": description: "List, search or filter chatbots"

		"chatbot-session.create": description:  "Create chatbot sessions"
		"chatbot-sessions.search": description: "List, search or filter chatbot sessions"

		"chatbot-session-step.create": description:  "Create chatbot session steps"
		"chatbot-session-steps.search": description: "List, search or filter chatbot session steps"

		"chatbot-session-handoff.create": description:  "Create chatbot session handoffs"
		"chatbot-session-handoffs.search": description: "List, search or filter chatbot session handoffs"

		"tenant.create": description:  "Create tenants"
		"tenants.search": description: "List, search or filter tenants"

		"project.create": description:  "Create projects"
		"projects.search": description: "List, search or filter projects"

		"project-ai-system.create": description:  "Create AI systems"
		"project-ai-systems.search": description: "List, search or filter AI systems"

		"project-fria-scenario.create": description:  "Create FRIA risk scenarios"
		"project-fria-scenarios.search": description: "List, search or filter FRIA risk scenarios"
	}
}
