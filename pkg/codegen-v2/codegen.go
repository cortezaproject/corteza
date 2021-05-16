package main

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v2/internal/def"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v2/internal/tpl"
	"text/template"
)

// rbac generates one rbac definition file per service
// <service>/service/rbac.gen.go
//
// Contains all RBAC related definitions
func accessControlService(t *template.Template, dd []*def.Document) (err error) {
	const (
		template      = "access_control.go.tpl"
		outputPathTpl = "%s/service/access_control.gen.go"
	)

	for component, perComponent := range partByComponent(dd) {
		w := tpl.Wrap{
			Package:   "service",
			Component: component,
			Def:       perComponent,
		}

		w.Imports = append(w.Imports, cImport(component, "types"))

		err = tpl.GoTemplate(fmt.Sprintf(outputPathTpl, component), t.Lookup(template), w)
		if err != nil {
			return
		}
	}

	return
}

// rbac generates one rbac definition file per service
// <service>/service/rbac.gen.go
//
// Contains all RBAC related definitions
func rbacTypes(t *template.Template, dd []*def.Document) (err error) {
	const (
		template      = "rbac.go.tpl"
		outputPathTpl = "%s/types/rbac.gen.go"
	)

	for component, perComponent := range partByComponent(dd) {
		w := tpl.Wrap{
			Package:   "types",
			Component: component,
			Def:       perComponent,
		}

		err = tpl.GoTemplate(fmt.Sprintf(outputPathTpl, component), t.Lookup(template), w)
		if err != nil {
			return
		}
	}

	return
}

func partByComponent(dd []*def.Document) map[string][]*def.Document {
	var (
		parted = make(map[string][]*def.Document)
	)

	for _, d := range dd {
		parted[d.Component] = append(parted[d.Component], d)
	}

	return parted
}

func cImport(c, s string) string {
	return fmt.Sprintf("github.com/cortezaproject/corteza-server/%s/%s", c, s)
}
