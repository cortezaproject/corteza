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

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			owned_by:   schema.AttributeUserRef
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
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
