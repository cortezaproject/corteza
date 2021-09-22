package gen

import (
	"fmt"
	"text/template"

	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/def"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/tpl"
)

func Locale(t *template.Template, dd []*def.Document) error {
	return List{
		"types":    localeTypes,
		"services": localeServices,
	}.Generate(t, dd)
}

func localeTypes(t *template.Template, dd []*def.Document) (err error) {
	const (
		templateName  = "locale/types.go.tpl"
		outputPathTpl = "%s/types/locale.gen.go"
	)

	dd = filter(dd, func(d *def.Document) bool { return d.Locale != nil })

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

func localeServices(t *template.Template, dd []*def.Document) (err error) {
	const (
		templateName  = "locale/service.go.tpl"
		outputPathTpl = "%s/service/locale.gen.go"
	)

	dd = filter(dd, func(d *def.Document) bool { return d.Locale != nil && !d.Locale.SkipSvc })

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

	for _, d := range dd {
		if d.Locale != nil {
		}
	}

	return
}
