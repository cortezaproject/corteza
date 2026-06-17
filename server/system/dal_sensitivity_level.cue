package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

dal_sensitivity_level: {
	model: {
		attributes: {
			id:     schema.IdField & {
				// hand-written tag uses "sensitivityLevelID", not the
				// "dalSensitivityLevelID" the res.ident convention would produce.
				json: "sensitivityLevelID,string"
			}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			level: {
				sortable: true,
				goType: "int"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}

			meta: {
				goType: "types.DalSensitivityLevelMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			// struct-only field: features.labels is false (no label store
			// codegen) but the hand-written struct still carries a plain
			// map[string]string Labels field, distinct from the feature
			// field's map[string]labelTypes.LabelValue type.
			labels: {
				goType: "map[string]string"
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "labels", omitEmpty: true }
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				// hand-written tag has no omitempty (record is always created)
				json: { field: "createdBy", string: true }
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			dal_sensitivity_level_id: {goType: "[]uint64", ident: "dalSensitivityLevelID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: { goType: "string" }

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["dal_sensitivity_level_id", "handle"]
		byNilState: ["deleted"]
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["sensitivity_level"]
		}
		store: {}
	}

	features: {
		labels: false
		projectScoped: true
	}

	types: {
		gen: true
	}

	service: {
		// No eventbus events are emitted by this service.
		events: false

		// Single "Manage" permission gates every op, the action-log message
		// templates all reference {{sensitivityLevel}} and the filter prop is
		// named "search" (see dal_sensitivity_level_actions.yaml).
		actionProp: "sensitivityLevel"
		filterProp: "search"

		// UndeleteByID exists in the companion file.
		undelete: true

		// Every CRUD body is bespoke (store.Tx wrapper, svc.prepare normalization,
		// DAL ReplaceSensitivityLevel / RemoveSensitivityLevel side-effects), so all
		// bodies delegate to hand-written on<Op> handlers.
		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

		// Access is a single CanManageDalSensitivityLevel check (not the standard
		// CanSearch* / CanCreate* names), performed inside the on<Op> handlers --
		// suppress the standard search/create access checks in the scaffold.
		customAccessOps: ["search", "create"]
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for user by ID

						It returns user even if deleted or suspended
						"""
				}
			]
		}
	}
}
