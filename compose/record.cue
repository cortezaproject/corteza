package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

record: {
	parents: [
		{handle: "namespace"},
		{handle: "module"},
	]

	model: {
		attributes: {
			id: schema.IdField
			module_id: { ident: "moduleID", goType: "uint64", storeIdent: "rel_module" }
			module: { goType: "*types.Module", store: false }
			values: { goType: "types.RecordValueSet", store: false }
			namespace_id: { ident: "namespaceID", goType: "uint64", storeIdent: "rel_namespace" }

			owned_by: { goType: "uint64" }
			created_at: schema.SortableTimestampField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: { goType: "uint64" }
			updated_by: { goType: "uint64" }
			deleted_by: { goType: "uint64" }
		}
	}

	filter: {
		struct: {
			module_id: { goType: "uint64" }
			namespace_id: { goType: "uint64" }
			query: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"owner.manage": {}
			"revisions.search": {}
		}
	}
}
