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
		ident: "compose_record"

		attributes: {
			id: schema.IdField
			revision: {
				goType: "uint"
				dal: { type: "Number" }
			}
			module_id: {
			 	ident: "moduleID",
				goType: "uint64",
				storeIdent: "rel_module"
				dal: { type: "Ref", refModelResType: "corteza::compose:module" }
			}
			module: {
				goType: "*types.Module",
				store: false
			}
			values: {
				goType: "types.RecordValueSet",
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			meta: {
				goType: "map[string]any",
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			namespace_id: {
				ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			owned_by:   schema.AttributeUserRef
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
			"idx_compose_record_base": {
				attributes: ["module_id", "namespace_id"]
				predicate: "deleted_at IS NULL"
			}
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
