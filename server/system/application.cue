package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_applicationDefs: {
			ApplicationUnify: { name: "ApplicationUnify", fields: [
				{ name: "Name", type: "string", json: "name,omitempty" },
				{ name: "Listed", type: "bool", json: "listed" },
				{ name: "Url", type: "string", json: "url" },
				{ name: "Config", type: "string", json: "config" },
				{ name: "Icon", type: "string", json: "icon,omitempty" },
				{ name: "IconID", type: "uint64", json: "iconID,string" },
				{ name: "Logo", type: "string", json: "logo,omitempty" },
				{ name: "LogoID", type: "uint64", json: "logoID,string" },
			]}
		}

application: {
	types: {
		gen: true
		defs: _applicationDefs
	}

	model: {
		attributes: {
			id: schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			name: {
				sortable: true
				dal: {}
				envoy: {
					identifier: true
				}
			}
			enabled: {
				goType: "bool"
				sortable: true,
				dal: { type: "Boolean", default: true }
			}
			weight: {
				goType: "int",
				sortable: true
				dal: { type: "Number", default: 0, meta: { "rdbms:type": "integer" } }
			}
			unify: {
				type: _applicationDefs.ApplicationUnify
				ptr: true
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				json: { field: "unify", omitEmpty: true }
			}
			owner_id:   {
				schema.AttributeUserRef,
				storeIdent: "rel_owner",
				ident: "ownerID"
				json: "ownerID"
				envoy: {
					store: {
						omitRefFilter: true
					}
					yaml: {
						identKeyAlias: ["owner"]
					}
				}
			}
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
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			name: {goType: "string"}
			// not sure about the type of flagged_ids
			flagged_ids: {goType: "[]uint64"}
			flags: {goType: "[]string"}
			inc_flags: {goType: "uint"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["name"]
		byValue: ["name"]
		byNilState: ["deleted"]
	}

	features: {
		flags: true
		labelResourceType: "application"
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Name"
			identKeyAlias: ["applications", "apps"]
		}
		store: {
			handleField: "Name"
		}
	}

	rbac: {
		operations: {
			read:
				description: "Read application"
			update:
				description: "Update application"
			delete:
				description: "Delete application"
		}
	}

	service: {

		undelete: true

		customBodyOps: ["search"]

		customFunctions: [
			{
				name: "Flag"
				cap:  "write"
				args: [
					{name: "app", goType: "*types.Application"},
					{name: "ownedBy", goType: "uint64"},
					{name: "f", goType: "string"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
			{
				name: "Unflag"
				cap:  "write"
				args: [
					{name: "app", goType: "*types.Application"},
					{name: "ownedBy", goType: "uint64"},
					{name: "f", goType: "string"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
			{
				name:   "Reorder"
				cap:    "write"
				action: "Reorder"
				args: [
					{name: "order", goType: "[]uint64"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
		]

		updateFields: ["Name", "Enabled", "Weight"]

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
					description: """
						searches for role by ID

						It returns role even if deleted or suspended
						"""
				},
			]

			functions: [
				{
					expIdent: "ApplicationMetrics"
					return: [ "*types.ApplicationMetrics"]
				}, {
					expIdent: "ReorderApplications"
					// not sure about the ident
					args: [ {ident: "order", goType: "[]uint64"}]
				},
			]
		}
	}
}
