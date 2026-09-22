package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_applicationDefs: {
			ApplicationMeta: { name: "ApplicationMeta", fields: [
				{ name: "Description", type: "string", json: "description,omitempty" },
			]}
			ApplicationUnify: { name: "ApplicationUnify", fields: [
				{ name: "Name", type: "string", json: "name,omitempty" },
				{ name: "Listed", type: "bool", json: "listed" },
				{ name: "Url", type: "string", json: "url" },
				{ name: "Config", type: "string", json: "config" },
				{ name: "Icon", type: "string", json: "icon,omitempty" },
				{ name: "IconID", type: "uint64", json: "iconID,string" },
				{ name: "Logo", type: "string", json: "logo,omitempty" },
				{ name: "LogoID", type: "uint64", json: "logoID,string" },
				// "" or "section" opens a shell section or a link; "custom" opens
				// the application's own HTML source in the sandboxed app view.
				{ name: "Kind", type: "string", json: "kind,omitempty" },
			]}
			// What the app view needs to know about the source without loading
			// it: the launcher payload carries this, never the HTML itself.
			ApplicationSourceMeta: { name: "ApplicationSourceMeta", fields: [
				{ name: "Hash", type: "string", json: "hash,omitempty" },
				{ name: "Size", type: "int", json: "size" },
				{ name: "Namespace", type: "string", json: "namespace,omitempty" },
				{ name: "Modules", goType: "[]string", json: "modules,omitempty" },
				// The declared modules the app may also create and change records in.
				{ name: "Writes", goType: "[]string", json: "writes,omitempty" },
				{ name: "UpdatedAt", type: "time.Time", ptr: true, json: "updatedAt,omitempty" },
				{ name: "UpdatedBy", type: "uint64", json: "updatedBy,string,omitempty" },
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
			meta: {
				type: _applicationDefs.ApplicationMeta
				ptr: true
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				json: { field: "meta", omitEmpty: true }
			}
			unify: {
				type: _applicationDefs.ApplicationUnify
				ptr: true
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				json: { field: "unify", omitEmpty: true }
			}
			// The custom app's HTML. Never serialized with the resource: it is
			// read through its own endpoint so list responses stay small.
			source: {
				dal: { type: "Text" }
				json: "-"
			}
			source_meta: {
				type: _applicationDefs.ApplicationSourceMeta
				ptr: true
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				json: { field: "sourceMeta", omitEmpty: true }
			}
			owner_id:   {
				schema.AttributeUserRef,
				storeIdent: "rel_owner",
				ident: "ownerID"
				json: { field: "ownerID", string: true }
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
			access:
				description: "Access application"
			update:
				description: "Update application"
			delete:
				description: "Delete application"
			"source.manage":
				description: "Replace the HTML source of a custom application"
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
				name:   "SetSource"
				cap:    "write"
				action: "SourceSet"
				args: [
					{name: "app", goType: "*types.Application"},
					{name: "source", goType: "string"},
					{name: "meta", goType: "*types.ApplicationSourceMeta"},
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
