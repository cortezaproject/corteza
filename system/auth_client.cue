package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

auth_client: {
	model: {
		attributes: {
			id:     schema.IdField
			handle: schema.HandleField
			meta: {
				goType: "*types.AuthClientMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			secret: {
				goType: "string"
				dal: { type: "Text", length: 64 }
			}
			scope: {
				goType: "string"
				dal: { type: "Text", length: 512 }
			}
			valid_grant: {
				goType: "string"
				dal: { type: "Text", length: 32 }
			}
			redirect_uri: {
				goType: "string",
				ident: "redirectURI"
				dal: {}
			}
			enabled: {
				sortable: true,
				goType: "bool"
				dal: { type: "Boolean", default: false }
			}
			trusted: {
				sortable: true,
				goType: "bool"
				dal: { type: "Boolean", default: false }
			}
			valid_from: schema.SortableTimestampNilField
			expires_at: schema.SortableTimestampNilField
			security: {
				goType: "*types.AuthClientSecurity"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			owned_by:   schema.AttributeUserRef
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
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
