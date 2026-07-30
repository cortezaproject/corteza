package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_project_ai_systemDefs: {
			ProjectAiSystemMeta: { name: "ProjectAiSystemMeta", fields: [
				{name: "Short", type: "string", json: "short"},
				{name: "Description", type: "string", json: "description,omitempty"},
				// IntendedPurpose is free prose (EU AI Act Art. 11 / Annex IV
				// technical documentation), so it lives in the JSON meta blob
				// rather than in a column of its own.
				{name: "IntendedPurpose", type: "string", json: "intendedPurpose,omitempty"},
			]}
		}

project_ai_system: {
	features: {
		labels:        false
	}

	types: {
		gen: true
		defs: _project_ai_systemDefs
	}

	model: {
		attributes: {
			id:        schema.IdField
			tenant_id: schema.TenantRefField
			project_id: schema.ProjectRefField & {
				// hand-written tag has no omitempty (owning project ref)
				json: {field: "projectID", string: true}
			}
			handle: schema.HandleField
			// EU AI Act Art. 6 risk classification of the AI system:
			// prohibited | high | limited | minimal. Deliberately a plain
			// string column — no DB-level enum constraint — so the vocabulary
			// can follow the regulation without a schema migration.
			risk_class: {
				goType: "string"
				dal: {type: "Text", length: 32}
			}
			meta: {
				type: _project_ai_systemDefs.ProjectAiSystemMeta
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": {attribute: "id"}
			"unique_handle_per_project": {
				fields: [{attribute: "project_id"}, {attribute: "handle", modifiers: ["LOWERCASE"]}]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			project_ai_system_id: {goType: "[]uint64", ident: "projectAiSystemID", storeIdent: "id"}
			tenant_id:            schema.TenantFilterField
			project_id:           schema.ProjectFilterField
			handle:               {goType: "string"}
			risk_class:           {goType: "string"}
			deleted:              {goType: "filter.State", storeIdent: "deleted_at"}
		}
		query: ["handle"]
		byValue: ["project_ai_system_id", "project_id", "handle", "risk_class"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:               "Read AI system"
			update: description:             "Update AI system"
			delete: description:             "Delete AI system"
			// the members of an AI system are project RESOURCES, not people
			"resources.manage": description: "Manage AI system resources"
		}
	}

	envoy: {omit: true}

	service: {
		extraServices: false

		customFunctions: [
			{
				name: "MemberList"
				cap:  "read"
				args: [{name: "projectAiSystemID", goType: "uint64"}]
				results: [
					{name: "set", goType: "types.ProjectAiSystemEntrySet"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "MemberAdd"
				cap:  "write"
				args: [
					{name: "projectAiSystemID", goType: "uint64"},
					{name: "resourceRef", goType: "string"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
			{
				name: "MemberRemove"
				cap:  "write"
				args: [
					{name: "projectAiSystemID", goType: "uint64"},
					{name: "resourceRef", goType: "string"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
		]


		filterProp: "search"

		undelete: false

		omitUpdateFields: ["handle"]

		hooks: {
			beforeCreate: true
			beforeUpdate: true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: "searches for AI system by ID"
				},
				{
					fields: ["project_id", "handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: "searches for AI system by project and handle; returns only non-deleted"
				},
			]
		}
	}
}
