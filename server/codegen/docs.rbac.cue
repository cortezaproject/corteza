package codegen

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
	"github.com/cortezaproject/corteza/server/app"
)

#_indexPayload: {
	label: string
	resources: [...string]
}

#_operationsPayload: {
	label: string

	operations: [...{
		slug: string
		label: string
		description: string
	}]
}

[...schema.#codegen] &
[
	for cmp in app.corteza.components {
		template: "docs/rbac.index.adoc.tpl"
		output:   "src/modules/generated/partials/access-control/\(cmp.handle)/index.gen.adoc"
		payload: #_indexPayload & {
			label: cmp.label
			resources: [ for res in cmp.resources { res.handle } ]
		}
	}
] +
[
	for cmp in app.corteza.components {
		template: "docs/rbac.$component.adoc.tpl"
		output:   "src/modules/generated/partials/access-control/\(cmp.handle)/component.gen.adoc"
		payload: #_operationsPayload & {
			label: cmp.label

			operations: [
				for op in cmp.rbac.operations {
					slug: "rbac-\(cmp.handle)-\(op.handle)"
					label: op.handle
					description: op.description
				}
			]
		}
	}
] +
[
	for cmp in app.corteza.components for res in cmp.resources {
		template: "docs/rbac.$resource.adoc.tpl"
		output:   "src/modules/generated/partials/access-control/\(cmp.handle)/resource.\(res.handle).gen.adoc"
		payload: #_operationsPayload & {
			label: res.handle

			operations: [
				for op in res.rbac.operations {
					slug: "rbac-\(res.handle)-\(op.handle)"
					label: op.handle
					description: op.description
				}
			]
		}
	}
]
