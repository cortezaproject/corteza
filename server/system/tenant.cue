package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

tenant: {
	features: {
		labels: true
	}

	model: {
		attributes: {
			id:     schema.IdField
			handle: schema.HandleField
			status: {
				sortable: true
				goType:   "types.TenantStatus"
				dal: {length: 32}
				omitSetter: true
				omitGetter: true
			}
			config: {
				goType: "types.TenantConfig"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			meta: {
				goType: "types.TenantMeta"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}

			created_at:   schema.SortableTimestampNowField
			updated_at:   schema.SortableTimestampNilField
			suspended_at: schema.SortableTimestampNilField
			deleted_at:   schema.SortableTimestampNilField
			created_by:   schema.AttributeUserRef
			updated_by:   schema.AttributeUserRef
			deleted_by:   schema.AttributeUserRef
		}

		indexes: {
			"primary": {attribute: "id"}
			"unique_handle": {
				fields: [{attribute: "handle", modifiers: ["LOWERCASE"]}]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			tenant_id: {goType: "[]uint64", ident: "tenantID", storeIdent: "id"}
			handle: {goType: "string"}
			status: {goType: "types.TenantStatus"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle"]
		byValue: ["tenant_id", "handle"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:             "Read tenant"
			update: description:           "Update tenant"
			delete: description:           "Delete tenant"
			suspend: description:          "Suspend tenant"
			"members.manage": description: "Manage tenant members"
		}
	}

	envoy: {
		omit: true
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for tenant by ID

						It also returns deleted tenants.
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for tenant by handle

						It returns only valid tenants (not deleted)
						"""
				},
			]
		}
	}
}
