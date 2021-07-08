package gen

import (
	"fmt"
	"text/template"

	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/def"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/tpl"
	"github.com/cortezaproject/corteza-server/pkg/slice"
)

func Envoy(t *template.Template, dd []*def.Document) error {
	return List{
		"resource rbac parse":      envoyResourceRbacUnmarshal,
		"resource rbac references": envoyResourceRbacReferences,
	}.Generate(t, dd)
}

// EnvoyResourceRbacUnmarshal envoy rbac unmarshal
// <service>/service/rbac.gen.go
//
// Contains all RBAC related definitions
func envoyResourceRbacUnmarshal(t *template.Template, dd []*def.Document) (err error) {
	const (
		templateName  = "envoy/resource-rbac_rules_parse.go.tpl"
		outputPathTpl = "pkg/envoy/resource/rbac_rules_parse.gen.go"
	)

	dd = filter(dd, func(d *def.Document) bool { return d.Envoy })

	// build list of component type imports
	ctImports := make([]string, 0)
	for _, d := range dd {
		imp := d.Component + "Types " + cImport(d.Component, "types")
		if !slice.HasString(ctImports, imp) {
			ctImports = append(ctImports, imp)
		}
	}

	w := tpl.Wrap{
		Package: "resource",
		Def:     dd,
		Imports: append(collectImports(dd...), ctImports...),
	}

	err = tpl.GoTemplate(outputPathTpl, t.Lookup(templateName), w)
	if err != nil {
		return
	}

	return
}

// EnvoyResourceRbacReferences generates one rbac definition file per component
// <service>/service/rbac.gen.go
//
// Contains all RBAC related definitions
func envoyResourceRbacReferences(t *template.Template, dd []*def.Document) (err error) {
	const (
		templateName  = "envoy/resource-rbac_references.go.tpl"
		outputPathTpl = "pkg/envoy/resource/rbac_references_%s.gen.go"
	)

	dd = filter(dd, func(d *def.Document) bool { return d.Envoy })

	for component, perComponent := range partByComponent(dd) {

		w := tpl.Wrap{
			Package:   "resource",
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
