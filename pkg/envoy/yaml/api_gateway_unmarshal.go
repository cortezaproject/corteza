package yaml

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"github.com/cortezaproject/corteza-server/system/types"
	"gopkg.in/yaml.v3"
)

func (wset *apiGatewaySet) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &apiGateway{}
		)

		if v == nil {
			return y7s.NodeErr(n, "malformed api gateway definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wrap *apiGateway) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.rbac = make(rbacRuleSet, 0, 10)
		wrap.res = &types.ApigwRoute{}
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

	if wrap.us, err = decodeUserstamps(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "endpoint":
			return y7s.DecodeScalar(v, "api gw endpoint", &wrap.res.Endpoint)
		case "method":
			return y7s.DecodeScalar(v, "api gw method", &wrap.res.Method)
		case "enabled":
			return y7s.DecodeScalar(v, "api gw enabled", &wrap.res.Enabled)
		case "group":
			return y7s.DecodeScalar(v, "api gw group", &wrap.res.Group)
		case "meta":
			aux := &types.ApigwRouteMeta{}
			err = v.Decode(&aux)
			if err != nil {
				return err
			}
			wrap.res.Meta = *aux
			return nil

		case "filters":
			wrap.filters = make(apiGwFilterSet, 0, 10)

			err = v.Decode(&wrap.filters)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (wrap *apiGwFilter) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.res = &types.ApigwFilter{}
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.ts, err = decodeTimestamps(n); err != nil {
		return
	}

	if wrap.us, err = decodeUserstamps(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {

		case "weight":
			return y7s.DecodeScalar(v, "route filter weight", &wrap.res.Weight)

		case "ref":
			return y7s.DecodeScalar(v, "route filter ref", &wrap.res.Ref)

		case "kind":
			return y7s.DecodeScalar(v, "route filter kind", &wrap.res.Kind)

		case "params":
			return v.Decode(&wrap.res.Params)
		}

		return nil
	})
}

func (wset apiGatewaySet) MarshalEnvoy() ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, len(wset)*2)

	for _, res := range wset {
		if tmp, err := res.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

func (wrap apiGateway) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewAPIGateway(wrap.res)
	rs.SetTimestamps(wrap.ts)
	rs.SetUserstamps(wrap.us)
	rs.SetConfig(wrap.envoyConfig)

	for _, f := range wrap.filters {
		trs := rs.AddGatewayFilter(f.res)
		trs.SetTimestamps(f.ts)
		trs.SetUserstamps(f.us)
	}

	return envoy.CollectNodes(
		rs,
		wrap.rbac.bindResource(rs),
	)
}
