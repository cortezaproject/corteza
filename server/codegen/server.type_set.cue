package codegen

import (
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

[...schema.#codegen] &
[
	// type_set.gen.go — one per component (resources + standalone types)
	for cmp in app.human.components {
		template: "gocode/types/$component_type_set.go.tpl"
		output:   "\(cmp.ident)/types/type_set.gen.go"
		payload: {
			package: "types"
			types: [
				for res in cmp.resources {
					expIdent:          res.expIdent
					noIdField:         res.features.noIdField
					labelResourceType: res.features.labelResourceType
				},
			] + [
				for t in cmp.types {
					expIdent:          t.expIdent
					noIdField:         t.noIdField
					labelResourceType: t.labelResourceType
				},
			]
		}
	},

	// type_set.gen_test.go — one per component
	for cmp in app.human.components {
		template: "gocode/types/$component_type_set_test.go.tpl"
		output:   "\(cmp.ident)/types/type_set.gen_test.go"
		payload: {
			package: "types"
			types: [
				for res in cmp.resources {
					expIdent:  res.expIdent
					noIdField: res.features.noIdField
				},
			] + [
				for t in cmp.types {
					expIdent:  t.expIdent
					noIdField: t.noIdField
				},
			]
		}
	},

	// type_labels.gen.go — one per component (labeled types only)
	for cmp in app.human.components {
		template: "gocode/types/$component_type_labels.go.tpl"
		output:   "\(cmp.ident)/types/type_labels.gen.go"
		payload: {
			package: "types"
			types: [
				for res in cmp.resources
				if res.features.labelResourceType != "" {
					expIdent:          res.expIdent
					labelResourceType: res.features.labelResourceType
				},
			] + [
				for t in cmp.types
				if t.labelResourceType != "" {
					expIdent:          t.expIdent
					labelResourceType: t.labelResourceType
				},
			]
		}
	},

	// type_set.gen.go — one per bundle
	for b in app.human.bundles {
		template: "gocode/types/$component_type_set.go.tpl"
		output:   "\(b.outputDir)/type_set.gen.go"
		payload: {
			package: b.package
			types: [
				for t in b.types {
					expIdent:          t.expIdent
					noIdField:         t.noIdField
					labelResourceType: t.labelResourceType
				},
			]
		}
	},

	// type_set.gen_test.go — one per bundle
	for b in app.human.bundles {
		template: "gocode/types/$component_type_set_test.go.tpl"
		output:   "\(b.outputDir)/type_set.gen_test.go"
		payload: {
			package: b.package
			types: [
				for t in b.types {
					expIdent:  t.expIdent
					noIdField: t.noIdField
				},
			]
		}
	},

	// type_labels.gen.go — one per bundle (labeled types only)
	for b in app.human.bundles {
		template: "gocode/types/$component_type_labels.go.tpl"
		output:   "\(b.outputDir)/type_labels.gen.go"
		payload: {
			package: b.package
			types: [
				for t in b.types
				if t.labelResourceType != "" {
					expIdent:          t.expIdent
					labelResourceType: t.labelResourceType
				},
			]
		}
	},
]
