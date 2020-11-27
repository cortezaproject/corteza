package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
	"gopkg.in/yaml.v3"
)

type (
	user struct {
		// when user is at least partially defined
		res *types.User `yaml:",inline"`

		// module's RBAC rules
		rbac rbacRuleSet
	}
	userSet []*user
)

// UnmarshalYAML resolves set of user definitions, either sequence or map
//
// When resolving map, key is used as handle
// Also supporting { handle: name } definitions
//
func (wset *userSet) UnmarshalYAML(n *yaml.Node) error {
	return each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &user{}
		)

		if v == nil {
			return nodeErr(n, "malformed user definition")
		}

		wrap.res = &types.User{
			EmailConfirmed: true,
		}

		switch v.Kind {
		case yaml.ScalarNode:
			if err = decodeScalar(v, "user email", &wrap.res.Email); err != nil {
				return
			}

		case yaml.MappingNode:
			if err = v.Decode(&wrap.res); err != nil {
				return
			}

		default:
			return nodeErr(n, "expecting scalar or map with user definitions")

		}

		if err = decodeRef(k, "user", &wrap.res.Handle); err != nil {
			return err
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wset userSet) MarshalEnvoy() ([]resource.Interface, error) {
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

func (wrap *user) UnmarshalYAML(n *yaml.Node) (err error) {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "user definition must be a map")
	}

	if wrap.res == nil {
		wrap.res = &types.User{}
	}

	if err = n.Decode(&wrap.res); err != nil {
		return
	}

	if wrap.rbac, err = decodeRbac(n); err != nil {
		return
	}

	return nil
}

func (wrap user) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewUser(wrap.res)
	return envoy.CollectNodes(
		rs,
		wrap.rbac.bindResource(rs),
	)
}
