package yaml

import (
	"github.com/cortezaproject/corteza-server/system/types"
	"gopkg.in/yaml.v3"
)

type (
	user struct {
		// when user is at least partially defined
		res *types.User `yaml:",inline"`

		// all known modules on a user
		modules ComposeModuleSet

		// module's RBAC rules
		*rbacRules
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

//func (wset userSet) MarshalEnvoy() ([]envoy.Node, error) {
//	// user usually have bunch of sub-resources defined
//	nn := make([]envoy.Node, 0, len(wset)*10)
//
//	for _, res := range wset {
//		if tmp, err := res.MarshalEnvoy(); err != nil {
//			return nil, err
//		} else {
//			nn = append(nn, tmp...)
//		}
//	}
//
//	return nn, nil
//}

func (wrap *user) UnmarshalYAML(n *yaml.Node) (err error) {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "user definition must be a map or scalar")
	}

	if wrap.res == nil {
		wrap.res = &types.User{}
	}

	if err = n.Decode(&wrap.res); err != nil {
		return
	}

	if wrap.rbacRules, err = decodeResourceAccessControl(types.UserRBACResource, n); err != nil {
		return
	}

	return nil
}

//func (wrap user) MarshalEnvoy() ([]envoy.Node, error) {
//	nn := make([]envoy.Node, 0, 1+len(wrap.modules))
//	nn = append(nn, &envoy.UserNode{Ns: wrap.res})
//
//	if tmp, err := wrap.modules.MarshalEnvoy(); err != nil {
//		return nil, err
//	} else {
//		nn = append(nn, tmp...)
//	}
//
//	// @todo rbac
//
//	//if tmp, err := wrap.rules.MarshalEnvoy(); err != nil {
//	//	return nil, err
//	//} else {
//	//	nn = append(nn, tmp...)
//	//}
//
//	return nn, nil
//}
