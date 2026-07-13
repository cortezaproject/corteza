package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_privacy: {
	features: {
		labels: false
	}

	types: {
		gen: true
	}

	model: {
		omitGetterSetter: true
		attributes: {
			id: schema.IdField & {json: "privacyID,string"}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField

			title: {
				goType:   "string"
				sortable: true
				dal: {type: "Text", length: 255}
			}
			description: {
				goType: "string"
				dal: {type: "Text", length: 0}
			}
			request_type: {
				goType:   "string"
				json:     "requestType,omitempty"
				sortable: true
				dal: {type: "Text", length: 128}
			}
			status: {
				goType:   "string"
				sortable: true
				dal: {type: "Text", length: 64}
			}
			severity: {
				goType: "string"
				dal: {type: "Text", length: 64}
			}
			risk: {
				goType: "string"
				dal: {type: "Text", length: 64}
			}

			request_owner: {
				goType: "uint64"
				json:   "requestOwner,string,omitempty"
				dal: {type: "Ref", refModelResType: "corteza::system:user"}
			}
			change_owner: {
				goType: "uint64"
				json:   "changeOwner,string,omitempty"
				dal: {type: "Ref", refModelResType: "corteza::system:user"}
			}
			change_approved_by: {
				goType: "uint64"
				json:   "changeApprovedBy,string,omitempty"
				dal: {type: "Ref", refModelResType: "corteza::system:user"}
			}

			risk_assessment: {
				goType: "string"
				json:   "riskAssessment,omitempty"
				dal: {type: "Text", length: 0}
			}
			change_required: {
				goType: "string"
				json:   "changeRequired,omitempty"
				dal: {type: "Text", length: 0}
			}
			risk_change: {
				goType: "string"
				json:   "riskChange,omitempty"
				dal: {type: "Text", length: 0}
			}
			// Backlog IDs, stored comma-separated for the PoC.
			backlog: {
				goType: "string"
				dal: {type: "Text", length: 0}
			}

			// Due date is stored as an ISO string (frontend treats it as a string
			// and sorts/compares client-side) — avoids timestamp param marshaling
			// in the generated REST controller.
			date_due: {
				goType: "string"
				json:   "dateDue,omitempty"
				dal: {type: "Text", length: 64}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {json: "createdBy,string"}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": {attribute: "id"}
		}
	}

	filter: {
		struct: {
			privacy_id: {goType: "[]uint64", ident: "privacyID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			status: {goType: "string"}
		}

		query: ["title", "status"]
		byValue: ["privacy_id", "project_id", "status"]
	}

	rbac: {
		operations: {
			read: description:   "Read project privacy"
			update: description: "Update project privacy"
			delete: description: "Delete project privacy"
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
						searches for project privacy by ID

						It also returns deleted project privacys.
						"""
				},
			]
		}
	}
}
