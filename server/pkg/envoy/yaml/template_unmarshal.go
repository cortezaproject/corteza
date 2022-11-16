package yaml

import (
	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
	"gopkg.in/yaml.v3"
)

func (wset *templateSet) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &template{}
		)

		if v == nil {
			return y7s.NodeErr(n, "malformed template definition")
		}

		wrap.res = &types.Template{}

		switch v.Kind {
		case yaml.MappingNode:
			if err = v.Decode(&wrap); err != nil {
				return
			}

		default:
			return y7s.NodeErr(n, "expecting scalar or map with template definitions")

		}

		if err = decodeRef(k, "template handle", &wrap.res.Handle); err != nil {
			return err
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wrap *template) UnmarshalYAML(n *yaml.Node) (err error) {
	if !y7s.IsKind(n, yaml.MappingNode) {
		return y7s.NodeErr(n, "template definition must be a map")
	}

	if wrap.res == nil {
		wrap.res = &types.Template{}
	}

	if err = n.Decode(&wrap.res); err != nil {
		return
	}

	if wrap.rbac, err = decodeRbac(n); err != nil {
		return
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.ts, err = decodeTimestamps(n); err != nil {
		return
	}

	return nil
}

func (wset templateSet) MarshalEnvoy() ([]resource.Interface, error) {
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

func (wrap template) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewTemplate(wrap.res)
	rs.SetTimestamps(wrap.ts)
	rs.SetConfig(wrap.envoyConfig)

	return envoy.CollectNodes(
		rs,
		wrap.rbac.bindResource(rs),
	)
}
