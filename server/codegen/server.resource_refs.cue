package codegen

import (
	"strings"
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

// Builds the per-resource payload for the resource-refs template.
// Emits ResourceRefs() methods on resource types (one file per resource).
// Complex refs (nested loops, conditional logic, string parsing) are
// delegated to a hand-written resourceRefsExt companion method via extended:true.
_ResourceRefsResource: {
	res = "res": schema.#Resource

	result: {
		expIdent: res.expIdent
		refs:     res.refs.items
		extended: res.refs.extended
	}
}

[...schema.#codegen] & [
	for cmp in app.human.components
	for res in cmp.resources
	if res.refs != _|_ {
		template: "gocode/types/$component_resource_refs.go.tpl"
		output:   "\(cmp.ident)/types/\(strings.Replace(res.handle, "-", "_", -1))_refs.gen.go"
		payload: {
			package: "types"
			(_ResourceRefsResource & {"res": res}).result
		}
	},
]
