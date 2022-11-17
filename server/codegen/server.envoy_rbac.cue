package codegen

import (
	"github.com/cortezaproject/corteza/server/app"
	"github.com/cortezaproject/corteza/server/codegen/schema"
	// "strings"
)

[...schema.#codegen] &
[
	for cmp in app.corteza.components {
		template: "gocode/envoy/rbac_references_$component.go.tpl"
		output:   "pkg/envoy/resource/rbac_references_\(cmp.ident).gen.go"
		payload: {
			package: "resource"
			imports: [
				"\"github.com/cortezaproject/corteza/server/\(cmp.ident)/types\"",
			]

			resources: [
				for res in cmp.resources if res.rbac != _|_ {
					rbacRefFunc: "\(cmp.expIdent)\(res.expIdent)RbacReferences"
					references: [
						for p in res.parents {p},
						{param: res.ident, refField: "ID", expIdent: res.expIdent},
					]
				},
			]
		}
	},
]+
[
	{
		template: "gocode/envoy/rbac_rules_parse.go.tpl"
		output:   "pkg/envoy/resource/rbac_rules_parse.gen.go"
		payload: {
			package: "resource"
			imports: [
				for cmp in app.corteza.components {
					"\(cmp.ident)Types \"github.com/cortezaproject/corteza/server/\(cmp.ident)/types\""
				},
			]

			resources: [
				for cmp in app.corteza.components for res in cmp.resources if res.rbac != _|_ {
					importAlias: "\(cmp.ident)Types"
					expIdent:    res.expIdent

					typeConst:   "\(importAlias).\(expIdent)ResourceType"
					rbacRefFunc: "\(cmp.expIdent)\(res.expIdent)RbacReferences"
					references: [
						for p in res.parents {p},
						{param: res.ident, refField: "ID", expIdent: res.expIdent},
					]
				},

				for cmp in app.corteza.components {
					importAlias: "\(cmp.ident)Types"
					expIdent:    cmp.expIdent

					typeConst: "\(importAlias).ComponentResourceType"
					references: []
				},
			]
		}
	},
]+
[
	{
		template: "gocode/envoy/resource_translation.go.tpl"
		output:   "pkg/envoy/resource/resource_translation.gen.go"
		payload: {
			package: "resource"
			resources: [
				for cmp in app.corteza.components for res in cmp.resources if res.locale != _|_ {
					expIdent: "\(cmp.expIdent)\(res.expIdent)"
					extended: res.locale.extended
				},
			]
		}
	},
]+
[
	{
		template: "gocode/envoy/resource_translation_parse.go.tpl"
		output:   "pkg/envoy/resource/resource_translation_parse.gen.go"
		payload: {
			package: "resource"
			imports: [
				for cmp in app.corteza.components for res in cmp.resources if res.locale != _|_ {
					"\(cmp.ident)Types \"github.com/cortezaproject/corteza/server/\(cmp.ident)/types\""
				},
			]

			resources: [
				for cmp in app.corteza.components for res in cmp.resources if res.locale != _|_ {
					importAlias:  "\(cmp.ident)Types"
					typeConst:    "\(importAlias).\(res.expIdent)ResourceTranslationType"
					resTrRefFunc: "\(cmp.expIdent)\(res.expIdent)ResourceTranslationReferences"
					references: [
						for p in res.parents {p},
						{param: res.ident, refField: "ID", expIdent: res.expIdent},
					]
				},
			]
		}
	},
]+
[
	// wrapped with additional for loop to trim out templates with empty types list
	for tpl in [
		for cmp in app.corteza.components {
			template: "gocode/envoy/resource_translation_references_$component.go.tpl"
			output:   "pkg/envoy/resource/resource_translation_references_\(cmp.ident).gen.go"
			payload: {
				package: "resource"
				imports: [
					"\"github.com/cortezaproject/corteza/server/\(cmp.ident)/types\"",
				]

				resources: [
					for res in cmp.resources if res.locale != _|_ {
						resTrRefFunc: "\(cmp.expIdent)\(res.expIdent)ResourceTranslationReferences"
						expIdent:     res.expIdent
						references: [
							for p in res.parents {p},
						]
					},
				]
			}
		},
	] if len(tpl.payload.resources) > 0 {tpl}
]
