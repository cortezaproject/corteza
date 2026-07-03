package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_tenant_membershipDefs: {
			TenantMemberRole: { name: "TenantMemberRole", values: [
				{ident: "TenantRoleAdmin", value: "admin"},
				{ident: "TenantRoleMember", value: "member"},
			]}
			TenantMemberStatus: { name: "TenantMemberStatus", values: [
				{ident: "TenantMemberStatusActive", value: "active"},
				{ident: "TenantMemberStatusSuspended", value: "suspended"},
				{ident: "TenantMemberStatusInvited", value: "invited"},
			]}
		}

tenant_membership: {
	features: {
		// hand-written TenantMembership struct has no Labels field
		labels: false
	}

	types: {
		gen: true
		defs: _tenant_membershipDefs
	}

	model: {
		attributes: {
			id: schema.IdField
			tenant_id: schema.TenantRefField & {
				// hand-written tag has no omitempty
				json: {field: "tenantID", string: true}
			}
			project_id: schema.ProjectRefField
			user_id: {
				ident:      "userID"
				goType:     "uint64"
				storeIdent: "rel_user"
				dal: {type: "ID"}
				// hand-written tag has no omitempty
				json: {field: "userID", string: true}
			}
			role: {
				type:       _tenant_membershipDefs.TenantMemberRole
				storeIdent: "role"
				dal: {length: 64}
				omitSetter: true
				omitGetter: true
			}
			status: {
				sortable:   true
				type:       _tenant_membershipDefs.TenantMemberStatus
				storeIdent: "status"
				dal: {length: 32}
				omitSetter: true
				omitGetter: true
			}
			invited_by: {
				ident:      "invitedBy"
				goType:     "uint64"
				storeIdent: "rel_invited_by"
				dal: {type: "ID", default: 0}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": {attribute: "id"}
			// hot path for token issuance: resolve tenant by user
			"index_user": {
				fields: [{attribute: "user_id"}]
			}
			"index_tenant": {
				fields: [{attribute: "tenant_id"}]
			}
			// one membership record per user per tenant
			"unique_tenant_user": {
				fields: [{attribute: "tenant_id"}, {attribute: "user_id"}]
			}
		}
	}

	filter: {
		struct: {
			tenant_membership_id: {goType: "[]uint64", ident: "tenantMembershipID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			user_id:    {goType: "uint64", ident: "userID", storeIdent: "rel_user"}
			role: {goType: "types.TenantMemberRole", storeIdent: "role"}
			status: {goType: "types.TenantMemberStatus", storeIdent: "status"}
		}

		byValue: ["tenant_membership_id", "tenant_id", "user_id"]
	}

	envoy: {
		omit: true
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: "searches for tenant membership by ID"
				}, {
					fields: ["user_id"]
					description: """
						searches for tenant membership by user

						Hot path: called at token issuance to resolve tenantID.
						"""
				}, {
					fields: ["tenant_id", "user_id"]
					description: "searches for tenant membership by tenant and user"
				},
			]
		}
	}
}
