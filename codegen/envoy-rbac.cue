package codegen

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
	// "strings"
)

envoyRBAC:
	[...schema.#codegen] &
	[
		for cmp in app.corteza.components {
			template: "gocode/envoy/rbac-references.go.tpl"
			output:   "pkg/envoy/resource/rbac_references_\(cmp.ident).gen.go"
			payload: {
				package: "resource"
				imports: [
					"\"github.com/cortezaproject/corteza-server/\(cmp.ident)/types\"",
				]

				//    cmpIdent: cmp.ident
				// Operation/resource validators, grouped by resource
				resources: [
					for res in cmp.resources {
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
			template: "gocode/envoy/rbac-rules-parse.go.tpl"
			output:   "pkg/envoy/resource/rbac_rules_parse.gen.go"
			payload: {
				package: "resource"
				imports: [
					for cmp in app.corteza.components {
						"\(cmp.ident)Types \"github.com/cortezaproject/corteza-server/\(cmp.ident)/types\""
					},
				]

				//    cmpIdent: cmp.ident
				// Operation/resource validators, grouped by resource
				resources: [
					for cmp in app.corteza.components for res in cmp.resources {
						importAlias: "\(cmp.ident)Types"
						expIdent:    res.expIdent

						typeConst: "\(importAlias).\(expIdent)ResourceType"
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
	]
