package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

role: {
	features: {
		projectScoped: true
		labels:        true
	}

	types: {
		gen: true
		// RoleMeta keeps a hand-written pointer-returning ParseRoleMeta (REST
		// request controllers depend on the *RoleMeta return), incompatible with
		// the value-returning generated version, so it is excluded from JSON helper
		// generation.
		jsonTypesPtr: ["RoleMeta"]
	}

	model: {
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			name: {
				sortable: true
				dal: {}
			}
			handle: schema.HandleField
			meta: {
				goType: "*types.RoleMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			archived_at: schema.SortableTimestampNilField
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			role_id:    {goType: "[]uint64", ident: "roleID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			member_id:  {goType: "uint64"}
			user_group_id: {goType: "uint64" }
			resource: {goType: "string" }
			handle: {goType: "string"}
			name: {goType: "string"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			archived: {goType: "filter.State", storeIdent: "archived_at"}
		}

		query: ["handle", "name"]
		byValue: ["role_id", "name", "handle"]
		byNilState: ["deleted", "archived"]
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["roles"]
		}
		store: {}
	}

	rbac: {
		operations: {
			read: description:             "Read role"
			update: description:           "Update role"
			delete: description:           "Delete role"
			"members.manage": description: "Manage members"
		}
	}

	service: {

		undelete: true

		// lookup: original FindByID = loadRole + proc with NO read-AC and a
		//   proc-based NotFound/label post-processor whose error handling the
		//   generated load-error wrapper can't reproduce -> custom body (onLookup).
		// update/delete/undelete: dispatch eventbus events (cbEvents=true),
		//   plus IsSystem guards and the undelete-reuses-Update-events quirk -> custom bodies.
		customBodyOps: ["lookup", "update", "delete", "undelete"]

		cbEvents:       true
		templateUpdate: true
		updateFields:   ["Handle", "Name", "Meta"]
		hooks: {
			validate:     true
			beforeCreate: true
			afterCreate:  true

			beforeSearch: true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for role by ID

						It returns role even if deleted or suspended
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for role by handle

						It returns only valid role (not deleted, not suspended)
						"""
				}, {
					fields: ["name"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for role by name

						It returns only valid role (not deleted, not suspended)
						"""
				},
			]
			functions: [
				{
					expIdent: "RoleMetrics"
					return: [ "*types.RoleMetrics"]
				},
			]
		}
	}
}
