package compose

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

module: {
	handle: "module"
	parents: [
		{handle: "namespace"},
	]

	model: {
		ident: "compose_module"
		attributes: {
			id: schema.IdField & {
				envoy: {
					yaml: {
						identKeyEncode: "moduleID"
					}
				}
			}
			namespace_id: {
				ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }

				envoy: {
					yaml: {
						identKeyAlias: ["namespace", "namespace_id", "ns", "ns_id"]
					}
				}
			}
			handle: schema.HandleField
			name: {
				sortable: true
				dal: {}
			}
			meta: {
				goType: "rawJson"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			config: {
				goType: "types.ModuleConfig"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			fields: {
				goType: "types.ModuleFieldSet",
				store: false
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						omitEncoder: true
					}
				}
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"namespace": { attribute: "namespace_id" },
			"unique_handle": {
				fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }, { attribute: "namespace_id" }]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			module_id: { goType: "[]uint64", ident: "moduleID", storeIdent: "id" }
			namespace_id: { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
			handle: { goType: "string" }
			name: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle", "name"]
		byValue: ["handle", "module_id", "namespace_id"]
		byNilState: ["deleted"]
	}

	envoy: {
		scoped: true
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["modules", "mod"]

			extendedResourcePostProcess: true
			extendedResourceDecoders: [{
				ident: "source"
				expIdent: "Source"
				// @deprecated records is what the old version used
				identKeys: ["source", "datasource", "records"]
				supportMappedInput: false
			}]
		}
		store: {
			postSetEncoder: true
			extendedEncoder: true
			extendedFilterBuilder: true
			extendedDecoder: true
		}
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"export":        description:  "Access to export modules"
			"record.create": description:  "Create record"
			"owned-record.create": description:  "Create record with custom owner"
			"records.search": description: "List, search or filter records"
		}
	}

	store: {
		ident: "composeModule"

		api: {
			lookups: [
				{
					fields: ["namespace_id", "handle"]
					constraintCheck: true
					nullConstraint: ["deleted_at"]
					description: """
						searches for compose module by handle (case-insensitive)
						"""
				}, {
					fields: ["namespace_id", "name"]
					nullConstraint: ["deleted_at"]
					description: """
						searches for compose module by name (case-insensitive)
						"""
				}, {
					fields: ["id"]
					description: """
						searches for compose module by ID

						It returns compose module even if deleted
						"""
				},
			]
		}
	}

	locale: {
		extended: true

		keys: {
			"name": {}
		}
	}

	//locale:
	//  resource:
	//    references: [ namespace, ID ]
	//
	//  extended: true
	//  keys:
	//    - name
}
