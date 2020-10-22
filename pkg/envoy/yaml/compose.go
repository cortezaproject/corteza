package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"gopkg.in/yaml.v3"
)

type (
	compose struct {
		namespaces ComposeNamespaceSet
		modules    ComposeModuleSet
		records    ComposeRecordSet
		// pages ComposePagesSet
		// charts ComposeChartsSet
	}
)

func (c *compose) UnmarshalYAML(n *yaml.Node) error {
	if !isKind(n, yaml.MappingNode) {
		// root node kind be mapping
		return nodeErr(n, "expecting mapping node")
	}

	var (
		nsRef string
		err   error
	)

	// 1st pass: handle doc-level references
	err = iterator(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "namespace":
			if def := findKeyNode(n, "namespaces"); def != nil {
				return nodeErr(def, "cannot combine namespace reference and namespaces definition")
			}

			if err := decodeScalar(v, "namespace ref", &nsRef); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 2nd pass: handle definitions
	return iterator(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "namespaces":
			return v.Decode(&c.namespaces)

		case "modules":
			if err = v.Decode(&c.modules); err != nil {
				return err
			}

			return c.modules.setNamespaceRef(nsRef)

		case "records":
			if err = v.Decode(&c.records); err != nil {
				return err
			}

			return c.records.setNamespaceRef(nsRef)
		}

		return nil
	})
}

func (wrap compose) MarshalEnvoy() ([]envoy.Node, error) {
	nn := make([]envoy.Node, 0, 100)

	if wrap.namespaces != nil {
		if tmp, err := wrap.namespaces.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	if wrap.modules != nil {
		if tmp, err := wrap.modules.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}
