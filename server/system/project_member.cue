package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_member: {
	features: {
		labels:        false
		projectScoped: true
	}

	types: {
		gen: true
	}

	model: {
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField & {
				// hand-written tag has no omitempty (owning project ref)
				json: {field: "projectID", string: true}
			}
			user_id: {
				ident:      "userID"
				goType:     "uint64"
				storeIdent: "rel_user"
				dal: {type: "ID"}
				// hand-written tag has no omitempty
				json: {field: "userID", string: true}
			}
			role_preset: {
				ident:      "rolePreset"
				goType:     "types.ProjectMemberRole"
				storeIdent: "role_preset"
				dal: {length: 64}
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
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": {attribute: "id"}
			// one membership record per user per project
			"unique_project_user": {
				fields: [{attribute: "project_id"}, {attribute: "user_id"}]
				predicate: "deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			project_member_id: {goType: "[]uint64", ident: "projectMemberID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			user_id: {goType: "uint64", ident: "userID", storeIdent: "rel_user"}
			role_preset: {goType: "types.ProjectMemberRole", ident: "rolePreset", storeIdent: "role_preset"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["project_member_id", "project_id", "user_id"]
		byNilState: ["deleted"]
	}

	envoy: {
		omit: true
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: "searches for project member by ID"
				}, {
					fields: ["project_id", "user_id"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for project member by project and user

						It returns only valid (not deleted) membership.
						"""
				},
			]
		}
	}
}
