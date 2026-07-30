package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_agentDefs: {
	AgentMeta: {name: "AgentMeta", fields: [
				{name: "Short", type:         "string", json: "short"},
				{name: "Description", type:   "string", json: "description,omitempty"},
				{name: "SidebarRoles", slice: true, type:     "string", json: "sidebarRoles,omitempty"},
	]}

	AgentBehavior: {name: "AgentBehavior", fields: [
				{name: "SystemPrompt", type:        "string", json:              "systemPrompt,omitempty"},
				{name: "Guardrails", slice:         true, type:                  "string", json: "guardrails,omitempty"},
				{name: "InjectSystemContext", type: "bool", json:                "injectSystemContext"},
				{name: "KnowledgeBases", goType:    "KnowledgeBaseIDList", json: "knowledgeBases,omitempty"},
				{name: "TreatyCLEnabled", ptr:      true, type:                  "bool", json: "treatyCLEnabled"},
				{name: "TreatyCLTemperature", type: "int", json:                 "tclTemperature,omitempty"},
				{name: "TreatyCLArticles", slice:   true, type:                  "string", json: "tclArticles,omitempty"},
	]}

	AgentExecution: {name: "AgentExecution", fields: [
				{name: "Model", type:  _agentDefs.AgentExecutionModel, json:  "model"},
				{name: "Limits", type: _agentDefs.AgentExecutionLimits, json: "limits"},
	]}

	AgentExecutionModel: {name: "AgentExecutionModel", fields: [
					{name: "LLMProviderID", type: "uint64", json: "llmProviderID,string,omitempty"},
					{name: "Model", type:         "string", json: "model,omitempty"},
					{name: "Temperature", ptr:    true, type:     "float64", json: "temperature,omitempty"},
	]}

	AgentExecutionLimits: {name: "AgentExecutionLimits", fields: [
					{name: "MaxIterations", type:  "int", json:     "maxIterations,omitempty"},
					{name: "Timeout", type:        "string", json:  "timeout"},
					{name: "SoftLimitRatio", type: "float64", json: "softLimitRatio,omitempty"},
					{name: "ContextWindow", type:  "int", json:     "contextWindow"},
					{name: "OutputTokens", type:   "int", json:     "outputTokens"},
	]}

	AgentAccess: {name: "AgentAccess", fields: [
				{name: "Context", type:     _agentDefs.AgentAccessContext, json: "context"},
				{name: "Tools", goType:     "[]AgentAccessTool", json:           "tools,omitempty"},
				{name: "TAQs", goType:      "[]AgentAccessTAQ", json:            "taqs,omitempty"},
				{name: "Workflows", goType: "[]AgentAccessWorkflow", json:       "workflows,omitempty"},
	]}

	AgentAccessContext: {name: "AgentAccessContext", fields: [
					{name: "Namespace", type:  "string", json:         "namespace,omitempty"},
					{name: "Module", type:     "string", json:         "module,omitempty"},
					{name: "Defaults", goType: "map[string]any", json: "defaults,omitempty"},
	]}

	AgentInvocation: {name: "AgentInvocation", fields: [
				{name: "User", type:   _agentDefs.AgentInvocationUser, json:   "user"},
				{name: "System", type: _agentDefs.AgentInvocationSystem, json: "system"},
	]}

	AgentInvocationUser: {name: "AgentInvocationUser", fields: [
					{name: "Enabled", type: "bool", json: "enabled"},
	]}

	AgentInvocationSystem: {name: "AgentInvocationSystem", fields: [
					{name: "Enabled", type:        "bool", json:            "enabled"},
					{name: "ServiceAccount", type: "uint64", json:          "serviceAccount,string,omitempty"},
					{name: "InputSchema", goType:  "json.RawMessage", json: "inputSchema,omitempty"},
					{name: "OutputFormat", type:   "string", json:          "outputFormat,omitempty"},
	]}

	AgentAccessTAQ: {name: "AgentAccessTAQ", fields: [
				{name: "ID", type:          "uint64", json:            "id,string"},
				{name: "Description", type: "string", json:            "description,omitempty"},
				{name: "Params", goType:    "map[string]string", json: "params,omitempty"},
	]}

	AgentAccessWorkflow: {name: "AgentAccessWorkflow", fields: [
					{name: "ID", type:          "uint64", json: "id,string"},
					{name: "Description", type: "string", json: "description,omitempty"},
	]}

	AgentAccessTool: {name: "AgentAccessTool", fields: [
				{name: "Name", type:        "string", json:                          "name"},
				{name: "Description", type: "string", json:                          "description"},
				{name: "Allow", slice:      true, type:                              _agentDefs.AgentAccessAllow, json: "allow"},
				{name: "Context", type:     _agentDefs.AgentAccessToolContext, json: "context,omitempty"},
	]}

	AgentAccessToolContext: {name: "AgentAccessToolContext", fields: [
					{name: "Defaults", goType:  "map[string]any", json: "defaults,omitempty"},
					{name: "Overrides", goType: "map[string]any", json: "overrides,omitempty"},
	]}

	AgentAccessAllow: {name: "AgentAccessAllow", fields: [
					{name: "NamespaceID", type: "uint64", json:            "namespaceID,string"},
					{name: "ModuleIDs", goType: "AgentAccessIDList", json: "moduleIDs"},
	]}
}

agent: {
	types: {
		gen: true
		idLists: ["AgentAccessIDList"]

		defs: _agentDefs
	}

	features: {
		labels:            true
		labelResourceType: "agent"
	}

	model: {
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle:     schema.HandleField
			status: {
				sortable: true
				goType:   "string"
				dal: {length: 32}
			}
			revision: {
				goType: "int"
				dal: {type: "Number", meta: {"rdbms:type": "integer"}}
			}
			meta: {
				type: _agentDefs.AgentMeta
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'short'.
			name: {
				sortableJSON: {json: "meta.short", accessor: "Meta.Short", nullable: false}
				store:      false
				goType:     "string"
				omitSetter: true
				omitGetter: true
				envoy: {yaml: {omitEncoder: true}}
			}
			behavior: {
				type: _agentDefs.AgentBehavior
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			execution: {
				type: _agentDefs.AgentExecution
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			access: {
				type: _agentDefs.AgentAccess
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			invocation: {
				type: _agentDefs.AgentInvocation
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				// hand-written tag has no omitempty (record is always created)
				json: {field: "createdBy", string: true}
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": {attribute: "id"}
			// Scoped to the project (2026-07-30), not global. A project revision
			// branch copies its agents, so the parent revision and its draft
			// necessarily hold same-handled agents at once -- a global unique
			// handle makes that impossible. It is also what the publish diff
			// already assumes: diffSources identifies resources across revisions
			// by kind + handle, which only works if a handle can repeat across
			// revisions. Widening a unique constraint keeps every currently
			// valid dataset valid.
			"unique_handle_per_project": {
				fields: [{attribute: "project_id"}, {attribute: "handle", modifiers: ["LOWERCASE"]}]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			agent_id: {goType: "[]uint64", ident: "agentID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: {goType: "string"}
			status: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle", "status"]
		byValue: ["agent_id", "project_id", "handle", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read agent"
			update: description: "Update agent"
			delete: description: "Delete agent"
		}
	}

	service: {
		// project-state lock: block writes unless the owning project is a draft
		guard: true

		genAccessController: true

		undelete: false

		customBodyOps: ["lookup", "search", "create", "update"]

		customFunctions: []
	}

	refs: {
		items: [
			{path: "Behavior.KnowledgeBases", kind:       "KindKnowledgeBase", reason:      "ReasonAgentKnowledgeBase", iter: "sliceID"},
			{path: "Execution.Model.LLMProviderID", kind: "KindLlmProvider", reason:        "ReasonAgentLlmProvider"},
			{path: "Access.TAQs", kind:                   "KindNgAutomation", reason:       "ReasonAgentAutomation", iter: "sliceField"},
			{path: "Access.Workflows", kind:              "KindAutomationWorkflow", reason: "ReasonAgentAutomation", iter: "sliceField"},
		]
		extended: true
	}

	// Envoy stays omitted -- revisited 2026-07-30, after trying the opposite.
	//
	// Un-omitting looks like the way to make a revision branch copy agents:
	// envoy exists to rewire references across an ID remap, which is the whole
	// problem a branch copy has to solve. It does not work. The system
	// component's encode cannot express a copy at all -- matchupAgents
	// (system/envoy/store_encode.gen.go) matches purely on node identifiers
	// (handle, ID) against an UNSCOPED search of every agent in the store, so a
	// copy always matches its source and becomes an UPDATE, moving the original
	// into the draft instead of duplicating it. The merge algorithms
	// (Replace/Skip/Panic) only choose what happens once a match is found; none
	// of them means "treat as new". Scope resolution is compose-only too
	// (getScopeNodes is a stub outside compose), so a scoped decode silently
	// reads every agent in the store.
	//
	// Both are fixable only by generating real scope support for system
	// resources, which rewrites the encode path of every system type -- far too
	// much blast radius for a branch-copy fix. The branch copies agents by hand
	// instead; see system/service/project_revision.go.
	envoy: {
		omit: true
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for agent by ID

						It also returns deleted agents.
						"""
				}, {
					// Project-scoped to match unique_handle_per_project above --
					// a bare handle no longer identifies one agent. Safe to
					// change: nothing outside generated store code called the
					// global LookupAgentByHandle.
					fields: ["project_id", "handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for agent by project and handle

						It returns only valid agents (not deleted)
						"""
				},
			]
		}
	}
}
