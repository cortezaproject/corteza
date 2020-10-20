package yaml

import (
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

	return iterator(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "namespace":
			c.namespaces = ComposeNamespaceSet{&ComposeNamespace{}}
			return v.Decode(&c.namespaces[0])
		case "namespaces":
			return v.Decode(&c.namespaces)
		case "modules":
			return v.Decode(&c.modules)
		case "records":
			return v.Decode(&c.records)
		}

		return nil
	})
}
