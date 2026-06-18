package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project: {
	features: {
		labels:       true
		tenantScoped: true
	}

	types: {
		gen: true
	}

	model: {
		attributes: {
			id:        schema.IdField
			tenant_id: schema.TenantRefField
			handle:    schema.HandleField
			status: {
				sortable: true
				goType:   "types.ProjectStatus"
				dal: {length: 32}
				omitSetter: true
				omitGetter: true
			}
			config: {
				goType: "types.ProjectConfig"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			meta: {
				goType: "types.ProjectMeta"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			governance: {
				goType: "types.ProjectGovernance"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				json: {field: "createdBy", string: true}
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": {attribute: "id"}
			"unique_handle": {
				fields: [{attribute: "handle", modifiers: ["LOWERCASE"]}]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			project_id: {goType: "[]uint64", ident: "projectID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			handle: {goType: "string"}
			status: {goType: "types.ProjectStatus"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle"]
		byValue: ["project_id", "handle"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:             "Read project"
			update: description:           "Update project"
			delete: description:           "Delete project"
			"members.manage": description: "Manage project members"
		}
	}

	envoy: {
		omit: true
	}

	service: {
		genConstructor: true

		events:   false

		lookup:   false
		update:   false
		undelete: false

		hooks: {
			beforeCreate: true
			beforeDelete: true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for project by ID

						It also returns deleted projects.
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for project by handle

						It returns only valid projects (not deleted)
						"""
				},
			]
		}
	}
}
