package gen

import (
	"fmt"
	"text/template"

	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/def"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/tpl"
)

func RBAC(t *template.Template, dd []*def.Document) error {
	return List{
		"type":                   rbacTypes,
		"service access control": rbacAccessControlService,
	}.Generate(t, dd)
}

// RbacTypes generates rbac definitions (one per component)
// <service>/service/rbac.gen.go
//
// Contains all RBAC related definitions
func rbacTypes(t *template.Template, dd []*def.Document) (err error) {
	const (
		templateName  = "rbac/types.go.tpl"
		outputPathTpl = "%s/types/rbac.gen.go"
	)

	for component, perComponent := range partByComponent(dd) {
		w := tpl.Wrap{
			Package:   "types",
			Component: component,
			Def:       perComponent,
		}

		err = tpl.GoTemplate(fmt.Sprintf(outputPathTpl, component), t.Lookup(templateName), w)
		if err != nil {
			return
		}
	}

	return
}

// RbacAccessControlService generates access control functions (one file per component)
// <service>/service/rbac.gen.go
//
// Contains all RBAC related definitions
func rbacAccessControlService(t *template.Template, dd []*def.Document) (err error) {
	const (
		templateName  = "rbac/access_control.go.tpl"
		outputPathTpl = "%s/service/access_control.gen.go"
	)

	for component, perComponent := range partByComponent(dd) {
		w := tpl.Wrap{
			Package:   "service",
			Component: component,
			Def:       perComponent,
			Imports:   append(collectImports(perComponent...), cImport(component, "types")),
		}

		err = tpl.GoTemplate(fmt.Sprintf(outputPathTpl, component), t.Lookup(templateName), w)
		if err != nil {
			return
		}
	}

	return
}
