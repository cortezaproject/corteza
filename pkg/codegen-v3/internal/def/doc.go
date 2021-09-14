package def

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/tpl"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	Document struct {
		Skip                bool `yaml:"(skip)"`
		Imports             []string
		Component           string
		IsComponentResource bool `yaml:"-"`
		Resource            string
		Source              string
		RBAC                *rbac
		Locale              *locale
		Envoy               bool `yaml:"envoy"`
	}
)

func (set *rbacOperations) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k *yaml.Node, v *yaml.Node) (err error) {
		def := rbacOperation{}
		if k != nil {
			def.Operation = k.Value
		}

		*set = append(*set, &def)
		return v.Decode(&def)
	})
}

// Preproc preprocesses the document and sets defaults
func (doc *Document) Proc(filename string) error {
	doc.Source = filename

	// filename parts
	fp := strings.Split(filename, ".")
	// trim extension
	fp = fp[:len(fp)-1]
	if len(fp) > 0 && doc.Component == "" {
		// set component from the 1st part
		// component is system, compose, ...
		doc.Component = fp[0]
	}

	if len(fp) > 1 && doc.Resource == "" {
		// if there are more parts, set resource
		// resource is user, module, record, workflow, ...
		doc.Resource = fp[1]
	}

	if strings.ToLower(doc.Resource) == "component" {
		return fmt.Errorf("can not use 'component' as a resource name, this is done automatically")
	} else if doc.Resource == "" {
		doc.Resource = "component"
		doc.IsComponentResource = true
	}

	doc.Imports = normalizeImport(doc.Imports...)

	if err := doc.RBAC.proc(doc.Component, doc.Resource); err != nil {
		return err
	}

	if err := doc.Locale.proc(doc.Component, doc.Resource); err != nil {
		return err
	}

	doc.Resource = tpl.Export(doc.Resource)

	return nil
}

func normalizeImport(ii ...string) []string {
	for i := range ii {
		if strings.Contains(ii[i], " ") {
			p := strings.SplitN(ii[i], " ", 2)
			ii[i] = fmt.Sprintf(`%s "%s"`, p[0], strings.Trim(p[1], `"`))
		} else {
			ii[i] = fmt.Sprintf(`"%s"`, strings.Trim(ii[i], `"'`+"`"))
		}
	}

	return ii
}
