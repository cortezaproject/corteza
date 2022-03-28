package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

module: schema.#Resource & {
	handle: "module"
	parents: [
		{handle: "namespace"},
	]

	struct: {
		id: schema.IdField
		handle: schema.HandleField
		meta: { goType: "rawJson" }
		fields: { goType: "types.ModuleFieldSet", store: false }
		namespace_id: { ident: "namespaceID", goType: "uint64", storeIdent: "rel_namespace" }
		name: {}

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
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

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"record.create": description:  "Create record"
			"records.search": description: "List, search or filter records"
		}
	}

	store: {
		ident: "composeModule"

		settings: {
			rdbms: {
				table: "compose_module"
			}
		}

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
