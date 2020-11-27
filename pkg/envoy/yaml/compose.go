package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	compose struct {
		namespaces composeNamespaceSet
		modules    composeModuleSet
		records    composeRecordSet
		pages      composePageSet
		charts     composeChartSet
	}
)

func (c *compose) UnmarshalYAML(n *yaml.Node) error {
	var (
		nsRef string
		err   error
	)

	// 1st pass: handle doc-level references
	err = eachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "namespace":
			if def := findKeyNode(n, "namespaces"); def != nil {
				return nodeErr(def, "cannot combine namespace reference and namespaces definition")
			}

			if err := decodeRef(v, "namespace", &nsRef); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 2nd pass: handle definitions
	return eachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "namespaces":
			return v.Decode(&c.namespaces)

		case "modules":
			if err = v.Decode(&c.modules); err != nil {
				return err
			}

			return c.modules.setNamespaceRef(nsRef)

		case "pages":
			if err = v.Decode(&c.pages); err != nil {
				return err
			}

			return c.pages.setNamespaceRef(nsRef)

		case "charts":
			if err = v.Decode(&c.charts); err != nil {
				return err
			}

			return c.charts.setNamespaceRef(nsRef)

		case "records":
			if err = v.Decode(&c.records); err != nil {
				return err
			}

			return c.records.setNamespaceRef(nsRef)

		}

		return nil
	})
}

func (c compose) MarshalEnvoy() ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, 100)

	if c.namespaces != nil {
		if tmp, err := c.namespaces.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	if c.modules != nil {
		if tmp, err := c.modules.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}
