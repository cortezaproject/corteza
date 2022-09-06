package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

auth_client: {
	model: {
		attributes: {
				id:     schema.IdField
				handle: schema.HandleField
				meta: {goType: "*types.AuthClientMeta"}
				secret: {goType: "string"}
				scope: {goType: "string"}
				valid_grant: {goType: "string"}
				redirect_uri: {goType: "string", ident: "redirectURI"}
				enabled: {sortable: true, goType: "bool"}
				trusted: {sortable: true, goType: "bool"}
				valid_from: { goType: "*time.Time" }
				expires_at: schema.SortableTimestampNilField
				security: {goType: "*types.AuthClientSecurity"}

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
			client_id: {goType: "[]uint64"}
			handle: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["handle"]
		byNilState: ["deleted"]
	}

	confirmed_client: {
		user_id: {goType: "uint64"}
		client_id: {goType: "uint64"}
		confirmed_at: {goType: "schema.OptTimestamp"}
	}

	confirmed_client_filter: {
		user_id: {goType: "uint64"}
	}

	rbac: {
		operations: {
			read: description:      "Read authorization client"
			update: description:    "Update authorization client"
			delete: description:    "Delete authorization client"
			authorize: description: "Authorize authorization client"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
							searches for auth client by ID

							It returns auth clint even if deleted
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for auth client by ID

						It returns auth clint even if deleted
						"""
				},
			]
		}
	}
}
