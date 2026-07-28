package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_backlog_item: {
	features: {
		labels: false
	}

	types: {
		gen: true
	}

	model: {
		omitGetterSetter: true
		attributes: {
			id:         schema.IdField & {json: "backlogItemID,string"}
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

			// Category is the parent category table this backlog item is
			// filed against (incident|task|feature|privacy|review).
			category: {
				goType:   "string"
				sortable: true
				dal: {type: "Text", length: 32}
			}

			// EventID references the parent category item's ID. It is a
			// polymorphic reference across the 5 category tables (keyed
			// together with category above) — no FK, plain column.
			event_id: {
				ident:    "eventID"
				expIdent: "EventID"
				goType:   "uint64"
				json:     "eventID,string,omitempty"
				dal: {type: "ID", default: 0}
			}

			assignee: {
				goType: "uint64"
				json:   "assignee,string,omitempty"
				dal: {type: "Ref", refModelResType: "corteza::system:user"}
			}

			priority: {
				goType: "string"
				dal: {type: "Text", length: 64}
			}

			status: {
				goType:   "string"
				sortable: true
				dal: {type: "Text", length: 64}
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
			backlog_item_id: {goType: "[]uint64", ident: "backlogItemID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			revision_id: {goType: "uint64", ident: "revisionID", storeIdent: "rel_revision"}
			event_id: {ident: "eventID", expIdent: "EventID", goType: "uint64"}
			category: {goType: "string"}
			status: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["title"]
		byValue: ["project_id", "revision_id", "event_id", "category", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read project backlog item"
			update: description: "Update project backlog item"
			delete: description: "Delete project backlog item"
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
			"Title", "Description", "Category", "EventID", "Assignee", "Priority", "Status", "DateDue",
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
						searches for project backlog item by ID

						It also returns deleted project backlog items.
						"""
				},
			]
		}
	}
}
