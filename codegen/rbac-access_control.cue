package codegen

import (
  "github.com/cortezaproject/corteza-server/codegen/schema"
	 "github.com/cortezaproject/corteza-server/app"
)

rbacAccessControl:
	[...schema.#codegen] &
	[
		for cmp in app.corteza.components {
			template: "gocode/rbac/access_control.go.tpl"
			output:   "\(cmp.ident)/service/access_control.gen.go"
			payload: {
				imports: [
					"github.com/cortezaproject/corteza-server/\(cmp.ident)/types",
				]
				package: "service"

				// All possible RBAC operations on component and resources
				// flattened
				operations: [
					for res in cmp.resources for op in res.rbac.operations {
						"op":            op.handle
						"const":         "types.\(res.expIdent)ResourceType"
						"ctor":          "types.\(res.expIdent)RbacResource(\(len(res.rbac.resource.references)*"0,"))"
						"goType":        res.goType
						"description":   op.description
						"checkFuncName": op.checkFuncName
					},
					for op in cmp.rbac.operations {
						"op":            op.handle
						"const":         "types.ComponentResourceType"
						"ctor":          "types.ComponentRbacResource()"
						"goType":        "types.Component"
						"description":   op.description
						"checkFuncName": op.checkFuncName
						"component":     true
					},
				]

				// Operation/resource validators, grouped by resource
				validation: [
					for res in cmp.resources {
						"label":      res.ident
						"const":      "types.\(res.expIdent)ResourceType"
						"funcName":   "rbac\(res.expIdent)ResourceValidator"
						"references": res.rbac.resource.references
						"operations": [ for op in res.rbac.operations {op.handle}]
					},
					{
						"label":    "\(cmp.ident) component"
						"const":    "types.ComponentResourceType"
						"funcName": "rbacComponentResourceValidator"
						"operations": [ for op in cmp.rbac.operations {op.handle}]
					},
				]
			}
		},
	]
