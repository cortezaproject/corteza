package yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"gopkg.in/yaml.v3"
)

type (
	ComposeNamespace struct {
		res     *types.Namespace `yaml:",inline"`
		ref     string
		Modules ComposeModuleSet
		Rbac    `yaml:",inline"`
	}
	ComposeNamespaceSet []*ComposeNamespace
)

// UnmarshalYAML resolves set of namespace definitions, either sequence or map
//
// When resolving map, key is used as slug
//
//
func (wset *ComposeNamespaceSet) UnmarshalYAML(n *yaml.Node) error {
	return iterator(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &ComposeNamespace{}
		)

		if v == nil {
			return nodeErr(n, "malformed namespace definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if k != nil {
			if wrap.res.Slug != "" {
				return nodeErr(k, "cannot define slug in mapped namespace definition")
			}

			if !handle.IsValid(k.Value) {
				return nodeErr(n, "namespace reference must be a valid handle")
			}

			// set namespace slug from map key value
			wrap.res.Slug = k.Value
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wrap *ComposeNamespace) UnmarshalYAML(n *yaml.Node) error {
	if isKind(n, yaml.ScalarNode) {
		wrap.ref = n.Value
		return nil
	}

	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "namespace definition must be a map or scalar")
	}

	if wrap.res == nil {
		wrap.res = &types.Namespace{
			// namespaces are enabled by default
			Enabled: true,
		}
	}

	return iterator(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "slug":
			return decodeScalar(v, &wrap.res.Slug)

		case "name":
			return decodeScalar(v, &wrap.res.Name)

		case "enabled":
			return decodeScalar(v, &wrap.res.Enabled)

		}

		return nil
	})
}
