package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_tenantDefs: {
			TenantConfig: { name: "TenantConfig", fields: [
				{name: "DALConnectionID", type: "uint64", json: "dalConnectionID,string,omitempty"},
				{name: "AuthProviders", slice: true, type: "string", json: "authProviders,omitempty"},
				{name: "FeatureFlags", goType: "map[string]bool", json: "featureFlags,omitempty"},
				{name: "Quotas", type: _tenantDefs.TenantQuotas, json: "quotas,omitempty"},
				{name: "Locale", type: "string", json: "locale,omitempty"},
				{name: "Timezone", type: "string", json: "timezone,omitempty"},
			]}
			TenantQuotas: { name: "TenantQuotas", fields: [
				{name: "MaxUsers", type: "int", json: "maxUsers,omitempty"},
				{name: "MaxProjects", type: "int", json: "maxProjects,omitempty"},
				{name: "MaxStorage", type: "int64", json: "maxStorage,omitempty"},
			]}
			TenantMeta: { name: "TenantMeta", fields: [
				{name: "Short", type: "string", json: "short,omitempty"},
				{name: "Description", type: "string", json: "description,omitempty"},
				{name: "LogoID", type: "uint64", json: "logoID,string,omitempty"},
				{name: "Color", type: "string", json: "color,omitempty"},
				{name: "Tags", slice: true, type: "string", json: "tags,omitempty"},
			]}
			TenantStatus: { name: "TenantStatus", values: [
				{ident: "TenantStatusActive", value: "active"},
				{ident: "TenantStatusSuspended", value: "suspended"},
				{ident: "TenantStatusArchived", value: "archived"},
			]}
		}

tenant: {
	features: {
		labels: true
	}

	types: {
		gen: true
		defs: _tenantDefs
	}

	model: {
		attributes: {
			id:     schema.IdField
			handle: schema.HandleField
			status: {
				sortable: true
				type:     _tenantDefs.TenantStatus
				dal: {length: 32}
				omitSetter: true
				omitGetter: true
			}
			config: {
				type: _tenantDefs.TenantConfig
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			meta: {
				type: _tenantDefs.TenantMeta
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}

			created_at:   schema.SortableTimestampNowField
			updated_at:   schema.SortableTimestampNilField
			suspended_at: schema.SortableTimestampNilField
			deleted_at:   schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				// hand-written tag has string but no omitempty
				json: {field: "createdBy", string: true}
			}
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

	service: {
		genConstructor: true
		filterProp:     "search"

		lookup:   false
		update:   false
		undelete: false

		hooks: {
			beforeCreate: true
			beforeDelete: true
		}

		customFunctions: [
			{
				name: "Suspend"
				cap:  "write"
				args: [{name: "ID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name: "Activate"
				cap:  "write"
				args: [{name: "ID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name: "Archive"
				cap:  "write"
				args: [{name: "ID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name: "SearchMembers"
				cap:  "read"
				args: [{name: "filter", goType: "types.TenantMembershipFilter"}]
				results: [{name: "set", goType: "types.TenantMembershipSet"}, {name: "f", goType: "types.TenantMembershipFilter"}, {name: "err", goType: "error"}]
			},
			{
				name: "Invite"
				cap:  "write"
				args: [{name: "m", goType: "*types.TenantMembership"}]
				results: [{name: "res", goType: "*types.TenantMembership"}, {name: "err", goType: "error"}]
			},
			{
				name: "AcceptInvite"
				cap:  "write"
				args: [{name: "tenantID", goType: "uint64"}, {name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name: "UpdateMember"
				cap:  "write"
				args: [{name: "m", goType: "*types.TenantMembership"}]
				results: [{name: "res", goType: "*types.TenantMembership"}, {name: "err", goType: "error"}]
			},
			{
				name: "RemoveMember"
				cap:  "write"
				args: [{name: "tenantID", goType: "uint64"}, {name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name: "SuspendMember"
				cap:  "write"
				args: [{name: "tenantID", goType: "uint64"}, {name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name: "ActivateMember"
				cap:  "write"
				args: [{name: "tenantID", goType: "uint64"}, {name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
		]
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
