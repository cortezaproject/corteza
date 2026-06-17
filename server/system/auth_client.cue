package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

auth_client: {
	features: {
		projectScoped: true
	}
	types: {
		gen: true
		// AuthClientMeta/AuthClientSecurity keep hand-written pointer-returning
		// ParseAuthClientMeta/ParseAuthClientSecurity (REST request controllers
		// depend on the *T return), incompatible with the value-returning
		// generated versions, so they are excluded from JSON helper generation.
		jsonTypesPtr: ["AuthClientMeta", "AuthClientSecurity"]
	}
	model: {
		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			meta: {
				goType: "*types.AuthClientMeta"
				json: { field: "meta", omitEmpty: true }
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'name'. Allows the FE to
			// request `?sort=name` and pages to be alphabetised by the
			// human-readable name stored inside the JSON meta column.
			name: {
				sortableJSON: { json: "meta.name", accessor: "Meta.Name" }
				store: false
				goType: "string"
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						omitEncoder: true
					}
				}
			}
			secret: {
				goType: "string"
				json: { field: "secret", omitEmpty: true }
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
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						customDecoder: true
						customEncoder: true
					}
				}
			}
			owned_by:   schema.AttributeUserRef & { json: "ownedBy" }
			created_at: schema.SortableTimestampNowField & { json: { field: "createdAt" } }
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & { json: "createdBy" }
			updated_by: schema.AttributeUserRef & { json: { field: "updatedBy", omitEmpty: true } }
			deleted_by: schema.AttributeUserRef & { json: { field: "deletedBy", omitEmpty: true } }
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			client_id: {goType: "[]uint64"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
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

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["authclients"]
		}
		store: {
			extendedRefDecoder: true
		}
	}

	rbac: {
		operations: {
			read: description:      "Read authorization client"
			update: description:    "Update authorization client"
			delete: description:    "Delete authorization client"
			authorize: description: "Authorize authorization client"
		}
	}

	service: {


		lookup:   false
		search:   false
		update:   false
		delete:   false
		undelete: true

		hooks: {
			validate:     true
			beforeCreate: true
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
