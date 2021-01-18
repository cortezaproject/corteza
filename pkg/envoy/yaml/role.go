package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"
	"github.com/cortezaproject/corteza-server/system/types"
	"gopkg.in/yaml.v3"
)

type (
	role struct {
		// when role is at least partially defined
		res    *types.Role
		ts     *resource.Timestamps
		config *resource.EnvoyConfig

		// all known modules on a role
		modules composeModuleSet

		// module's RBAC rules
		rbac rbacRuleSet
	}
	roleSet []*role
)

// UnmarshalYAML resolves set of role definitions, either sequence or map
//
// When resolving map, key is used as handle
// Also supporting { handle: name } definitions
//
func (wset *roleSet) UnmarshalYAML(n *yaml.Node) error {
	return Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &role{}
		)

		if v == nil {
			return NodeErr(n, "malformed role definition")
		}

		wrap.res = &types.Role{
			// no special defaults
		}

		switch v.Kind {
		case yaml.ScalarNode:
			if err = DecodeScalar(v, "role name", &wrap.res.Name); err != nil {
				return
			}

		case yaml.MappingNode:
			if err = v.Decode(&wrap); err != nil {
				return
			}
		}

		if err = decodeRef(k, "role", &wrap.res.Handle); err != nil {
			return err
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wset roleSet) MarshalEnvoy() ([]resource.Interface, error) {
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

func (wrap *role) UnmarshalYAML(n *yaml.Node) (err error) {
	if !IsKind(n, yaml.MappingNode) {
		return NodeErr(n, "role definition must be a map")
	}

	if wrap.res == nil {
		wrap.res = &types.Role{}
	}

	if err = n.Decode(&wrap.res); err != nil {
		return
	}

	if wrap.rbac, err = decodeRbac(n); err != nil {
		return
	}

	if wrap.config, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.ts, err = decodeTimestamps(n); err != nil {
		return
	}

	return nil
}

func (wrap role) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewRole(wrap.res)
	rs.SetTimestamps(wrap.ts)
	rs.SetConfig(wrap.config)

	return envoy.CollectNodes(
		rs,
		wrap.rbac.bindResource(rs),
	)
}
