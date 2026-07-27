package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_review: {
	features: {
		labels: false
	}

	types: {
		gen: true
	}

	model: {
		omitGetterSetter: true
		attributes: {
			id:         schema.IdField & {json: "reviewID,string"}
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
			review_type: {
				goType:   "string"
				json:     "reviewType,omitempty"
				sortable: true
				dal: {type: "Text", length: 128}
			}
			review_frequency: {
				goType: "string"
				json:   "reviewFrequency,omitempty"
				dal: {type: "Text", length: 64}
			}
			scope: {
				goType: "string"
				dal: {type: "Text", length: 0}
			}

			reviewer: {
				goType: "uint64"
				json:   "reviewer,string,omitempty"
				dal: {type: "Ref", refModelResType: "corteza::system:user"}
			}
			approved_by: {
				goType: "uint64"
				json:   "approvedBy,string,omitempty"
				dal: {type: "Ref", refModelResType: "corteza::system:user"}
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
			review_id: {goType: "[]uint64", ident: "reviewID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			revision_id: {goType: "uint64", ident: "revisionID", storeIdent: "rel_revision"}
			status: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["title", "status"]
		byValue: ["review_id", "project_id", "revision_id", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read project review"
			update: description: "Update project review"
			delete: description: "Delete project review"
		}
	}

	envoy: {
		omit: true
	}

	service: {
		extraServices:       false
		genAccessController: true

		updateFields: [
			"Title", "Description", "ReviewType", "ReviewFrequency", "Scope",
			"Reviewer", "ApprovedBy", "Status", "DateDue",
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
						searches for project review by ID

						It also returns deleted project reviews.
						"""
				},
			]
		}
	}
}
