package codegen

import (
	"github.com/crusttech/human/server/codegen/schema"
	"github.com/crusttech/human/server/app"
)

#_operation: {
	handle:      string
	description: string
}

[...schema.#codegen] &
[
	for cmp in app.human.components {
		template: "docs/permissions.md.tpl"
		output:   "reference/permissions/\(cmp.handle).gen.md"
		payload: {
			handle: cmp.handle
			label:  cmp.label

			operations: [ for op in cmp.rbac.operations {#_operation & {handle: op.handle, description: op.description}}]

			resources: [
				for res in cmp.resources if res.rbac != _|_ {
					handle: res.handle
					operations: [ for op in res.rbac.operations {#_operation & {handle: op.handle, description: op.description}}]
				},
			]
		}
	},
]+
[
	{
		template: "docs/permissions.index.md.tpl"
		output:   "reference/permissions/index.gen.md"
		payload: {
			components: [ for cmp in app.human.components {handle: cmp.handle, label: cmp.label}]
		}
	},
]
