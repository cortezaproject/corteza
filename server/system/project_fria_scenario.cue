package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_project_fria_scenarioDefs: {
	// Everything on a FRIA risk scenario that is prose or a taxonomy key list
	// lives here rather than in columns of its own.
	//
	// The []string fields hold TAXONOMY KEYS from the frontend's
	// sections/project/config/friaTaxonomies.js (TRIGGER_CONDITIONS,
	// IMPACTED_PARTIES, VULNERABLE_GROUPS, FUNDAMENTAL_RIGHTS,
	// AI_HARM_VECTORS). Those keys are a PERSISTED-DATA INTERFACE: append-only,
	// never renamed. Renaming a key silently orphans every stored scenario that
	// selected it — the row keeps a string nothing in the taxonomy answers to,
	// and the selection disappears from the UI without an error anywhere.
	//
	// The vocabularies themselves will move with EU AI Act Art. 27 guidance, so
	// they are deliberately a JSON blob: adding a trigger type or a rights
	// chapter must not cost a schema migration.
	ProjectFriaScenarioMeta: {name: "ProjectFriaScenarioMeta", fields: [
		// 1. harm narrative (title and severity are real columns)
		{name: "Description", type: "string", json: "description,omitempty"},

		// 2. trigger conditions
		{name: "TriggerTypes", slice: true, type: "string", json: "triggerTypes,omitempty"},
		{name: "TriggerDescription", type: "string", json: "triggerDescription,omitempty"},

		// 3. impacted parties + vulnerable groups
		{name: "ImpactedParties", slice: true, type: "string", json: "impactedParties,omitempty"},
		{name: "VulnerableGroups", slice: true, type: "string", json: "vulnerableGroups,omitempty"},
		{name: "VulnerableGroupsNotes", type: "string", json: "vulnerableGroupsNotes,omitempty"},

		// 4. fundamental rights
		{name: "Rights", slice: true, type: "string", json: "rights,omitempty"},

		// 5. AI harm vectors
		{name: "HarmVectors", slice: true, type: "string", json: "harmVectors,omitempty"},
		{name: "HarmVectorsDescription", type: "string", json: "harmVectorsDescription,omitempty"},
	]}
}

project_fria_scenario: {
	features: {
		labels: false
	}

	types: {
		gen:  true
		defs: _project_fria_scenarioDefs
	}

	model: {
		attributes: {
			id:        schema.IdField
			tenant_id: schema.TenantRefField
			project_id: schema.ProjectRefField & {
				// hand-written tag has no omitempty (owning project ref)
				json: {field: "projectID", string: true}
			}
			// The AI system this scenario assesses. Art. 27 attaches the FRIA
			// to a specific high-risk AI SYSTEM, not to a project, and a project
			// can hold several — so this is a real, filterable, sortable column
			// rather than a meta key: every scenario list is scoped by it.
			ai_system_id: {
				ident:      "aiSystemID"
				goType:     "uint64"
				storeIdent: "rel_ai_system"
				dal: {type: "Ref", refModelResType: "corteza::system:project-ai-system"}
				// hand-written tag has no omitempty
				json: {field: "aiSystemID", string: true}
				sortable: true
			}
			title: {
				goType:   "string"
				sortable: true
				dal: {type: "Text", length: 255}
			}
			// Severity of the assessed harm: low | medium | high | critical.
			// Deliberately a plain string column — no DB-level enum constraint —
			// so the scale can follow EU guidance without a schema migration,
			// same call as project_ai_system's risk_class.
			severity: {
				goType:   "string"
				sortable: true
				dal: {type: "Text", length: 32}
			}
			meta: {
				type: _project_fria_scenarioDefs.ProjectFriaScenarioMeta
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {json: "createdBy,string"}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		// No unique-handle index: a scenario has no handle. It is identified by
		// ID and titled in prose — two scenarios on the same AI system may
		// legitimately carry the same title.
		indexes: {
			"primary": {attribute: "id"}
		}
	}

	filter: {
		struct: {
			project_fria_scenario_id: {goType: "[]uint64", ident: "projectFriaScenarioID", storeIdent: "id"}
			tenant_id:                schema.TenantFilterField
			project_id:               schema.ProjectFilterField
			ai_system_id:             {goType: "uint64", ident: "aiSystemID", storeIdent: "rel_ai_system"}
			severity:                 {goType: "string"}
			deleted:                  {goType: "filter.State", storeIdent: "deleted_at"}
		}
		query: ["title"]
		byValue: ["project_fria_scenario_id", "project_id", "ai_system_id", "severity"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read FRIA risk scenario"
			update: description: "Update FRIA risk scenario"
			delete: description: "Delete FRIA risk scenario"
		}
	}

	envoy: {omit: true}

	service: {
		extraServices:       false
		genAccessController: true

		filterProp: "search"

		undelete: false

		updateFields: ["AiSystemID", "Title", "Severity", "UpdatedBy"]

		hooks: {
			beforeCreate: true
			beforeUpdate: true
			beforeDelete: true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for FRIA risk scenario by ID

						It also returns deleted FRIA risk scenarios.
						"""
				},
			]
		}
	}
}
