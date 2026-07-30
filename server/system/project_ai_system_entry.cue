package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_ai_system_entry: {
	features: {
		labels:  false
		paging:  false
		sorting: false
		checkFn: false
	}

	types: {
		gen: true
	}

	model: {
		attributes: {
			// struct-only field: there is no `id` store column (primary key is
			// the composite project_ai_system_id + resource_ref). expIdent MUST
			// be "ID" — the default ToTitle(ident) would yield "Id".
			id: {
				goType:     "uint64"
				expIdent:   "ID"
				store:      false
				omitGetter: true
				omitSetter: true
				json:       "-"
			}
			project_ai_system_id: {
				ident:      "projectAiSystemID"
				goType:     "uint64"
				storeIdent: "rel_project_ai_system"
				dal: {type: "Ref", refModelResType: "corteza::system:project-ai-system"}
				// hand-written tag has no omitempty
				json: {field: "projectAiSystemID", string: true}
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
				fields: [{attribute: "project_ai_system_id"}, {attribute: "resource_ref"}]
			}
		}
	}

	filter: {
		struct: {
			project_ai_system_id: {goType: "uint64", ident: "projectAiSystemID", storeIdent: "rel_project_ai_system"}
			resource_ref:         {goType: "string", ident: "resourceRef", storeIdent: "resource_ref"}
		}
		byValue: ["project_ai_system_id", "resource_ref"]
	}

	envoy: {omit: true}

	store: {
		api: {
			lookups: [
				{
					fields: ["project_ai_system_id", "resource_ref"]
					description: "searches for AI system entry by AI system and resource ref"
				},
			]
		}
	}
}
