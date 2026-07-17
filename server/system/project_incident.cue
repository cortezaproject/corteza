package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_incident: {
	features: {
		labels: false
	}

	types: {
		gen: true
	}

	model: {
		omitGetterSetter: true
		attributes: {
			id: schema.IdField & {json: "incidentID,string"}
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
			incident_type: {
				goType:   "string"
				json:     "incidentType,omitempty"
				sortable: true
				dal: {type: "Text", length: 128}
			}
			group_system: {
				goType: "string"
				json:   "groupSystem,omitempty"
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

			issue_owner: {
				goType: "uint64"
				json:   "issueOwner,string,omitempty"
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

			risk_issue: {
				goType: "string"
				json:   "riskIssue,omitempty"
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
			// Due/completed dates are stored as ISO strings (frontend treats them
			// as strings and sorts/compares client-side) — avoids timestamp param
			// marshaling in the generated REST controller.
			date_due: {
				goType: "string"
				json:   "dateDue,omitempty"
				dal: {type: "Text", length: 64}
			}
			completed_date: {
				goType: "string"
				json:   "completedDate,omitempty"
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
			incident_id: {goType: "[]uint64", ident: "incidentID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			status: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["title", "status"]
		byValue: ["incident_id", "project_id", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read project incident"
			update: description: "Update project incident"
			delete: description: "Delete project incident"
		}
	}

	envoy: {
		omit: true
	}

	service: {
		extraServices:       false
		genAccessController: true

		updateFields: [
			"Title", "Description", "IncidentType", "GroupSystem", "Status", "Severity", "Risk",
			"IssueOwner", "ChangeOwner", "ChangeApprovedBy", "RiskIssue", "ChangeRequired", "RiskChange",
			"DateDue", "CompletedDate",
		]

		hooks: {
			beforeCreate: true
			beforeUpdate: true
			beforeDelete: true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for project incident by ID

						It also returns deleted project incidents.
						"""
				},
			]
		}
	}
}
