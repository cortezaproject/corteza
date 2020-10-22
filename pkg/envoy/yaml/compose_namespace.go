package yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"gopkg.in/yaml.v3"
)

type (
	ComposeNamespace struct {
		// when namespace is at least partially defined
		res *types.Namespace `yaml:",inline"`

		// all known modules on a namespace
		modules ComposeModuleSet

		// module's RBAC rules
		*rbacRules
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

func (wset ComposeNamespaceSet) MarshalEnvoy() ([]envoy.Node, error) {
	// namespace usually have bunch of sub-resources defined
	nn := make([]envoy.Node, 0, len(wset)*10)

	for _, res := range wset {
		if tmp, err := res.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

func (wrap *ComposeNamespace) UnmarshalYAML(n *yaml.Node) error {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "namespace definition must be a map or scalar")
	}

	if wrap.res == nil {
		wrap.rbacRules = &rbacRules{}
		wrap.res = &types.Namespace{
			// namespaces are enabled by default
			Enabled: true,
		}
	}

	return iterator(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "slug":
			return decodeScalar(v, "namespace slug", &wrap.res.Slug)

		case "name":
			return decodeScalar(v, "namespace name", &wrap.res.Name)

		case "enabled":
			return decodeScalar(v, "namespace enabled", &wrap.res.Enabled)

		case "modules":
			return v.Decode(&wrap.modules)

		case "allow", "deny":
			return wrap.rbacRules.DecodeResourceRules(types.NamespaceRBACResource, k, v)

		}

		return nil
	})
}

func (wrap ComposeNamespace) MarshalEnvoy() ([]envoy.Node, error) {
	nn := make([]envoy.Node, 0, 1+len(wrap.modules))
	nn = append(nn, &envoy.ComposeNamespaceNode{Ns: wrap.res})

	if tmp, err := wrap.modules.MarshalEnvoy(); err != nil {
		return nil, err
	} else {
		nn = append(nn, tmp...)
	}

	// @todo rbac

	//if tmp, err := wrap.rules.MarshalEnvoy(); err != nil {
	//	return nil, err
	//} else {
	//	nn = append(nn, tmp...)
	//}

	return nn, nil
}
