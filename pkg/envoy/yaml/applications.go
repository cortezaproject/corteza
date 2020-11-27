package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
	"gopkg.in/yaml.v3"
)

type (
	application struct {
		// when application is at least partially defined
		res *types.Application `yaml:",inline"`

		// module's RBAC rules
		rbac rbacRuleSet
	}
	applicationSet []*application
)

// UnmarshalYAML resolves set of application definitions, either sequence or map
//
// When resolving map, key is used as handle
// Also supporting { handle: name } definitions
func (wset *applicationSet) UnmarshalYAML(n *yaml.Node) error {
	return eachSeq(n, func(v *yaml.Node) (err error) {
		var (
			wrap = &application{}
		)

		if v == nil {
			return nodeErr(n, "malformed application definition")
		}

		wrap.res = &types.Application{
			Enabled: true,
		}

		switch v.Kind {
		case yaml.ScalarNode:
			if err = decodeScalar(v, "application name", &wrap.res.Name); err != nil {
				return
			}

		case yaml.MappingNode:
			if err = v.Decode(&wrap.res); err != nil {
				return
			}

		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wset applicationSet) MarshalEnvoy() ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, len(wset))

	for _, res := range wset {
		if tmp, err := res.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}

	}

	return nn, nil
}

func (wrap *application) UnmarshalYAML(n *yaml.Node) (err error) {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "application definition must be a map")
	}

	if wrap.res == nil {
		wrap.res = &types.Application{}
	}

	if err = n.Decode(&wrap.res); err != nil {
		return
	}

	if wrap.rbac, err = decodeRbac(n); err != nil {
		return
	}

	return nil
}

func (wrap application) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewApplication(wrap.res)
	return envoy.CollectNodes(
		rs,
		wrap.rbac.bindResource(rs),
	)
}
