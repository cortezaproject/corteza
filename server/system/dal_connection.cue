package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

dal_connection: {
	types: {
		// generate the DalConnection struct into types/dal_connection.gen.go
		gen: true
		// []dal.Issue (struct-only Issues field) needs the dal package qualifier
		imports: ["github.com/crusttech/human/server/pkg/dal"]
	}
	model: {
		attributes: {
			// hand-written tag is connectionID,string; convention would emit
			// dalConnectionID,string (derived from the resource ident).
			id:     schema.IdField & { json: { field: "connectionID", string: true } }
			handle: schema.HandleField
			type: {
				sortable: true
				dal: {}
			}

			config: {
				goType: "types.DalConnectionConfig"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			meta: {
				goType: "types.DalConnectionMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			// Struct-only: relational/computed set of DAL issues, not stored.
			issues: {
				goType: "[]dal.Issue"
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "issues", omitEmpty: true }
			}
			// Struct-only: features.labels is disabled (custom label type), so
			// the generator does not auto-append a Labels field. Reproduce the
			// hand-written map[string]string field here.
			labels: {
				goType: "map[string]string"
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "labels", omitEmpty: true }
			}
			// Virtual sortable mapped onto meta->>'name'.
			name: {
				sortableJSON: { json: "meta.name", accessor: "Meta.Name", nullable: false }
				store: false
				goType: "string"
				omitSetter: true
				omitGetter: true
				envoy: { yaml: { omitEncoder: true } }
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			// hand-written tag is createdBy,string (no omitempty); convention
			// for a uint64 Ref would add ,omitempty.
			created_by: schema.AttributeUserRef & { json: { field: "createdBy", string: true } }
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			dal_connection_id: {goType: "[]uint64", ident: "dalConnectionID", storeIdent: "id"}
			handle: {goType: "string"}
			type: {goType: "string"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["dal_connection_id", "handle", "type"]
		byNilState: ["deleted"]
	}

	features: {
		labels: false
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["connection", "connections"]
		}
		store: {
			extendedRefDecoder: true
		}
	}

	rbac: {
		operations: {
			"read": description:         "Read connection"
			"update": description:       "Update connection"
			"delete": description:       "Delete connection"
			"dal-config.manage": description: "Manage DAL configuration"
		}
	}

	service: {
		events: false

		// No undelete in the generated public contract; the hand-written
		// UndeleteByID lives in the companion file as a custom method.
		undelete: false

		// Action-log prop is named "connection" (not the resource ident).
		actionProp: "connection"
		filterProp: "search"

		// Every CRUD op has a bespoke body (proc enrichment, dal manager
		// side-effects, primary-connection handling, type/name validation)
		// so each is delegated to its on<Op> handler.
		customBodyOps: ["lookup", "search", "create", "update", "delete"]

		// lookup/update/delete already delegate access through customBodyOps.
		// search and create also need their access check moved into the body
		// (search builds a filter.Check first; create runs the missing-name
		// guard before the access check), so they opt out of the generated
		// access scaffold too.
		customAccessOps: ["search", "create"]
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for connection by ID

						It returns connection even if deleted or suspended
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for connection by handle

						It returns only valid connection (not deleted)
						"""
				},
			]
		}
	}
}
