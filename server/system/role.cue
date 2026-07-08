package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_roleDefs: {
	RoleMeta: {name: "RoleMeta", fields: [
				{name: "Description", type: "string", json:             "description,omitempty"},
				{name: "Context", type:     _roleDefs.RoleContext, ptr: true, json: "context,omitempty"},
	]}
	RoleContext: {name: "RoleContext", fields: [
				{name: "Resource", type: "string", slice: true, json: "resourceTypes,omitempty"},
				{name: "Expr", type:     "string", json:  "expr,omitempty"},
	]}
	RoleMetrics: {name: "RoleMetrics", fields: [
				{name: "Total", type:         "uint", json:  "total"},
				{name: "Valid", type:         "uint", json:  "valid"},
				{name: "Deleted", type:       "uint", json:  "deleted"},
				{name: "Archived", type:      "uint", json:  "archived"},
				{name: "DailyCreated", type:  "uint", slice: true, json: "dailyCreated"},
				{name: "DailyDeleted", type:  "uint", slice: true, json: "dailyDeleted"},
				{name: "DailyUpdated", type:  "uint", slice: true, json: "dailyUpdated"},
				{name: "DailyArchived", type: "uint", slice: true, json: "dailyArchived"},
	]}
}

role: {
	features: {
		labels:            true
		labelResourceType: "role"
	}

	types: {
		gen: true
		// RoleMeta keeps a hand-written pointer-returning ParseRoleMeta (REST
		// request controllers depend on the *RoleMeta return), incompatible with
		// the value-returning generated version, so it is excluded from JSON helper
		// generation.
		jsonTypesPtr: ["RoleMeta"]

		defs: _roleDefs
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
				type: _roleDefs.RoleMeta
				ptr:  true
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}

			archived_at: schema.SortableTimestampNilField
			created_at:  schema.SortableTimestampNowField
			updated_at:  schema.SortableTimestampNilField
			deleted_at:  schema.SortableTimestampNilField
		}

		indexes: {
			"primary": {attribute: "id"}
		}
	}

	filter: {
		struct: {
			role_id: {goType: "[]uint64", ident: "roleID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			member_id: {goType: "uint64"}
			user_group_id: {goType: "uint64"}
			resource: {goType: "string"}
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
			mappedField:        "Handle"
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

		customFunctions: [
			{
				name:   "Archive"
				cap:    "write"
				action: "Archive"
				args: [{name: "roleID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "Unarchive"
				cap:    "write"
				action: "Unarchive"
				args: [{name: "roleID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:  "CloneRules"
				cap:   "write"
				ac:    "CanGrant"
				acErr: "ErrNotAllowedToCloneRules"
				args: [{name: "roleID", goType: "uint64"}, {name: "cloneToRoleID", goType: "...uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name: "Membership"
				cap:  "read"
				args: [{name: "userID", goType: "uint64"}]
				results: [{name: "mm", goType: "types.RoleMemberSet"}, {name: "err", goType: "error"}]
			},
			{
				name:   "MemberList"
				cap:    "read"
				action: "Members"
				args: [{name: "roleID", goType: "uint64"}]
				results: [{name: "mm", goType: "types.RoleMemberSet"}, {name: "err", goType: "error"}]
			},
			{
				name:   "MemberAdd"
				cap:    "write"
				action: "MemberAdd"
				args: [{name: "roleID", goType: "uint64"}, {name: "memberID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "MemberAddGroup"
				cap:    "write"
				action: "MemberAdd"
				args: [{name: "roleID", goType: "uint64"}, {name: "userGroupID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "MemberRemove"
				cap:    "write"
				action: "MemberRemove"
				args: [{name: "roleID", goType: "uint64"}, {name: "memberID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "MemberRemoveGroup"
				cap:    "write"
				action: "MemberRemove"
				args: [{name: "roleID", goType: "uint64"}, {name: "userGroupID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
		]

		cbEvents:       true
		templateUpdate: true
		updateFields: ["Handle", "Name", "Meta"]
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
