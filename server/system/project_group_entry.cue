package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_group_entry: {
	features: {
		labels:  false
		paging:  false
		sorting: false
		checkFn: false
	}

	model: {
		attributes: {
			project_group_id: {
				ident:      "projectGroupID"
				goType:     "uint64"
				storeIdent: "rel_project_group"
				dal: {type: "Ref", refModelResType: "corteza::system:project-group"}
			}
			resource_ref: {
				ident:      "resourceRef"
				goType:     "string"
				storeIdent: "resource_ref"
				dal: {type: "Text", length: 512}
			}
			created_at: schema.SortableTimestampNowField
		}

		indexes: {
			"primary": {
				fields: [{attribute: "project_group_id"}, {attribute: "resource_ref"}]
			}
		}
	}

	filter: {
		struct: {
			project_group_id: {goType: "uint64", ident: "projectGroupID", storeIdent: "rel_project_group"}
			resource_ref:     {goType: "string", ident: "resourceRef", storeIdent: "resource_ref"}
		}
		byValue: ["project_group_id", "resource_ref"]
	}

	envoy: {omit: true}

	store: {
		api: {
			lookups: [
				{
					fields: ["project_group_id", "resource_ref"]
					description: "searches for group entry by group and resource ref"
				},
			]
		}
	}
}
