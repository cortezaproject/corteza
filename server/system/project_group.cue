package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_project_groupDefs: {
			ProjectGroupMeta: { name: "ProjectGroupMeta", fields: [
				{name: "Short", type: "string", json: "short"},
				{name: "Description", type: "string", json: "description,omitempty"},
			]}
		}

project_group: {
	features: {
		labels:        false
	}

	types: {
		gen: true
		defs: _project_groupDefs
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
			meta: {
				type: _project_groupDefs.ProjectGroupMeta
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
			project_group_id: {goType: "[]uint64", ident: "projectGroupID", storeIdent: "id"}
			tenant_id:        schema.TenantFilterField
			project_id:       schema.ProjectFilterField
			handle:           {goType: "string"}
			deleted:          {goType: "filter.State", storeIdent: "deleted_at"}
		}
		query: ["handle"]
		byValue: ["project_group_id", "project_id", "handle"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:             "Read project group"
			update: description:           "Update project group"
			delete: description:           "Delete project group"
			"members.manage": description: "Manage project group members"
		}
	}

	envoy: {omit: true}

	service: {

		customFunctions: [
			{
				name: "MemberList"
				cap:  "read"
				args: [{name: "projectGroupID", goType: "uint64"}]
				results: [
					{name: "set", goType: "types.ProjectGroupEntrySet"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "MemberAdd"
				cap:  "write"
				args: [
					{name: "projectGroupID", goType: "uint64"},
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
					{name: "projectGroupID", goType: "uint64"},
					{name: "resourceRef", goType: "string"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
		]

		genConstructor: true

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
					description: "searches for project group by ID"
				},
				{
					fields: ["project_id", "handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: "searches for project group by project and handle; returns only non-deleted"
				},
			]
		}
	}
}
