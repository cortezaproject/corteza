package codegen

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
	"strings"
)

rbacTypes:
	[...schema.#codegen] &
	[
		for cmp in app.corteza.components {
			template: "gocode/rbac/types.go.tpl"
			output:   "\(cmp.ident)/types/rbac.gen.go"
			payload: {
				package: "types"

				cmpIdent: cmp.ident
				// Operation/resource validators, grouped by resource
				types: [
					for res in cmp.resources {
						"const":   "\(res.expIdent)ResourceType"
						"type":    res.rbac.resource.type
						"resFunc": "\(res.expIdent)RbacResource"
						"tplFunc": "\(res.expIdent)RbacResourceTpl"
						"attFunc": "\(res.expIdent)RbacAttributes"
						"goType":  res.expIdent

						"references": [ for field in res.rbac.resource.references { strings.ToTitle(field) } ]
					},
					{
						"const":     "ComponentResourceType"
						"type":      cmp.rbac.resource.type
						"resFunc":   "ComponentRbacResource"
						"tplFunc":   "ComponentRbacResourceTpl"
						"attFunc":   "ComponentRbacAttributes"
						"goType":    "Component"
						"component": true
					},
				]
			}
		},
	]
