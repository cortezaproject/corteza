package gen

import (
	"fmt"
	"text/template"

	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/def"
)

type (
	List map[string]func(*template.Template, []*def.Document) error
)

func (gg List) Generate(tpls *template.Template, dd []*def.Document) (err error) {
	for l, g := range gg {
		if err = g(tpls, dd); err != nil {
			return fmt.Errorf("codegen for %s failed: %w", l, err)
		}
	}

	return
}

func filter(dd []*def.Document, check func(*def.Document) bool) []*def.Document {
	aux := make([]*def.Document, 0, len(dd))
	for _, d := range dd {
		if !check(d) {
			continue
		}

		aux = append(aux, d)
	}
	return aux
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

func collectImports(dd ...*def.Document) []string {
	mm := make(map[string]bool)
	for _, d := range dd {
		for _, i := range d.Imports {
			mm[i] = true
		}
	}

	ii := make([]string, 0, len(mm))
	for i := range mm {
		ii = append(ii, i)
	}

	return ii
}

// component import
func cImport(c, s string) string {
	return fmt.Sprintf(`"github.com/cortezaproject/corteza-server/%s/%s"`, c, s)
}
