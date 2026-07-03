package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_recordDefs: {
			RecordValueSet: { name: "RecordValueSet", elem: _recordDefs.RecordValue, elemPtr: true }
			RecordValue: { name: "RecordValue", fields: [
				{ name: "RecordID", type: "uint64", json: "-" },
				{ name: "Name", type: "string", json: "name" },
				{ name: "Value", type: "string", json: "value,omitempty" },
				{ name: "Ref", type: "uint64", json: "-" },
				{ name: "Place", goType: "uint", json: "place,omitempty" },
				{ name: "DeletedAt", type: "time.Time", ptr: true, json: "deletedAt,omitempty" },
				{ name: "Updated", type: "bool", json: "-" },
				{ name: "OldValue", type: "string", json: "-" },
			]}
		}

record: {
	features: {
	}

	types: {
		defs: _recordDefs
	}

	parents: [
		{handle: "namespace"},
		{handle: "module"},
	]

	model: {
		ident: "compose_record"

		defaultSetter: true
		defaultGetter: true

		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			revision: {
				goType: "int"
				dal: { type: "Number", meta: { "rdbms:type": "integer" }, default: 0 }
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
				omitSetter: true
				omitGetter: true
			}
			values: {
				type: _recordDefs.RecordValueSet
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
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
			owned_by:   schema.AttributeUserRef & {
				identAlias: ["ownedBy", "OwnedBy", "owned_by"]
			}
			created_by: schema.AttributeUserRef & {
				identAlias: ["createdBy", "CreatedBy", "created_by"]
			}
			updated_by: schema.AttributeUserRef & {
				identAlias: ["updatedBy", "UpdatedBy", "updated_by"]
			}
			deleted_by: schema.AttributeUserRef & {
				identAlias: ["deletedBy", "DeletedBy", "deleted_by"]
			}
			created_by_agent: schema.AttributeAgentRef & {
				identAlias: ["createdByAgent", "CreatedByAgent", "created_by_agent"]
			}
		}

		indexes: {
			"primary": { attribute: "id" }
			"idx_compose_record_base": {
				attributes: ["module_id", "namespace_id"]
				predicate: "deleted_at IS NULL"
			}
		}
	}

	// @todo tmp
	envoy: {
		omit: true
	}

	filter: {
		struct: {
			tenant_id:    schema.TenantFilterField
			project_id:   schema.ProjectFilterField
			module_id:    { goType: "uint64" }
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
			"undelete": {}
			"owner.manage": {}
			"revisions.search": {}
		}
	}
}
