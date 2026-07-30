package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_task: {
	features: {
		labels: false
	}

	types: {
		gen: true
	}

	model: {
		omitGetterSetter: true
		attributes: {
			id:         schema.IdField & {json: "taskID,string"}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField

			// RevisionID points at the revision (a projects row) this item is
			// assigned to. Null means unassigned.
			revision_id: {
				ident:      "revisionID"
				expIdent:   "RevisionID"
				goType:     "uint64"
				json:       "revisionID,string,omitempty"
				storeIdent: "rel_revision"
				dal: {type: "ID", default: 0}
				sortable: true
			}

			title: {
				goType:   "string"
				sortable: true
				dal: {type: "Text", length: 255}
			}
			description: {
				goType: "string"
				dal: {type: "Text", length: 0}
			}
			task_name: {
				goType:   "string"
				json:     "taskName,omitempty"
				sortable: true
				dal: {type: "Text", length: 128}
			}
			task_type: {
				goType: "string"
				json:   "taskType,omitempty"
				dal: {type: "Text", length: 128}
				sortable: true
			}
			status: {
				goType:   "string"
				sortable: true
				dal: {type: "Text", length: 64}
			}
			severity: {
				goType: "string"
				dal: {type: "Text", length: 64}
				sortable: true
			}
			risk: {
				goType: "string"
				dal: {type: "Text", length: 64}
				sortable: true
			}

			owner: {
				goType: "uint64"
				json:   "owner,string,omitempty"
				dal: {type: "Ref", refModelResType: "corteza::system:user"}
				sortable: true
			}
			change_owner: {
				goType: "uint64"
				json:   "changeOwner,string,omitempty"
				dal: {type: "Ref", refModelResType: "corteza::system:user"}
				sortable: true
			}

			// Due/completed dates are stored as ISO strings (frontend treats them
			// as strings and sorts/compares client-side) — avoids timestamp param
			// marshaling in the generated REST controller.
			date_due: {
				goType: "string"
				json:   "dateDue,omitempty"
				dal: {type: "Text", length: 64}
				sortable: true
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
			task_id: {goType: "[]uint64", ident: "taskID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			revision_id: {goType: "uint64", ident: "revisionID", storeIdent: "rel_revision"}
			status: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["title", "status"]
		byValue: ["task_id", "project_id", "revision_id", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read project task"
			update: description: "Update project task"
			delete: description: "Delete project task"
		}
	}

	envoy: {
		omit: true
	}

	service: {
		extraServices:       false
		genAccessController: true

		updateFields: [
			"RevisionID",
			"Title", "Description", "TaskName", "TaskType", "Status", "Severity", "Risk",
			"Owner", "ChangeOwner", "DateDue", "CompletedDate",
		]

		hooks: {
			beforeCreate: true
			beforeUpdate: true
			beforeDelete: true
			beforeSearch: true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for project task by ID

						It also returns deleted project tasks.
						"""
				},
			]
		}
	}
}
