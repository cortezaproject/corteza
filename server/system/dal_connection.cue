package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_dal_connectionDefs: {
			DalConnectionConfig: { name: "DalConnectionConfig", fields: [
				{ name: "DAL", ptr: true, type: _dal_connectionDefs.DalConnectionConfigDAL, json: "dal,omitempty" },
				{ name: "Privacy", type: _dal_connectionDefs.DalConnectionConfigPrivacy, json: "privacy" },
			]}
			DalConnectionConfigPrivacy: { name: "DalConnectionConfigPrivacy", fields: [
				{ name: "SensitivityLevelID", type: "uint64", json: "sensitivityLevelID,string,omitempty" },
			]}
			DalConnectionConfigDAL: { name: "DalConnectionConfigDAL", fields: [
				{ name: "Type", type: "string", json: "type" },
				{ name: "Params", goType: "map[string]any", json: "params" },
				{ name: "ModelIdent", type: "string", json: "modelIdent" },
				{ name: "ModelIdentCheck", slice: true, type: "string", json: "modelIdentCheck" },
			]}
			DalConnectionMeta: { name: "DalConnectionMeta", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Ownership", type: "string", json: "ownership" },
				{ name: "Location", goType: "geolocation.Full", json: "location" },
				{ name: "Properties", type: _dal_connectionDefs.DalConnectionMetaProperties, json: "properties" },
			]}
			DalConnectionMetaProperties: { name: "DalConnectionMetaProperties", fields: [
				{ name: "DataAtRestEncryption", type: _dal_connectionDefs.DalConnectionMetaProperty, json: "dataAtRestEncryption" },
				{ name: "DataAtRestProtection", type: _dal_connectionDefs.DalConnectionMetaProperty, json: "dataAtRestProtection" },
				{ name: "DataAtTransitEncryption", type: _dal_connectionDefs.DalConnectionMetaProperty, json: "dataAtTransitEncryption" },
				{ name: "DataRestoration", type: _dal_connectionDefs.DalConnectionMetaProperty, json: "dataRestoration" },
			]}
			DalConnectionMetaProperty: { name: "DalConnectionMetaProperty", fields: [
				{ name: "Enabled", type: "bool", json: "enabled" },
				{ name: "Notes", type: "string", json: "notes" },
			]}
		}

dal_connection: {
	types: {
		// generate the DalConnection struct into types/dal_connection.gen.go
		gen: true
		// []dal.Issue (struct-only Issues field) needs the dal package qualifier
		imports: ["github.com/crusttech/human/server/pkg/dal"]
		defs: _dal_connectionDefs
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
				type: _dal_connectionDefs.DalConnectionConfig
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			meta: {
				type: _dal_connectionDefs.DalConnectionMeta
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
		// No undelete in the generated public contract; the hand-written
		// UndeleteByID lives in the companion file as a custom method.
		undelete: false

		actionProp: "connection"
		filterProp: "search"

		customBodyOps: ["lookup", "search", "create", "update", "delete"]

		customFunctions: [
			{name: "ReloadConnections", cap: "write"},
		]

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
