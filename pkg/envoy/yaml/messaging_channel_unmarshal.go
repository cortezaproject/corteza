package yaml

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

func (wset *messagingChannelSet) UnmarshalYAML(n *yaml.Node) error {
	return y7s.EachSeq(n, func(v *yaml.Node) (err error) {
		var (
			wrap = &messagingChannel{}
		)

		if v == nil || !y7s.IsKind(v, yaml.MappingNode) {
			return y7s.NodeErr(n, "malformed messagingChannel definition")
		}

		wrap.res = &types.Channel{}
		if err = v.Decode(&wrap); err != nil {
			return
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wset messagingChannelSet) MarshalEnvoy() ([]resource.Interface, error) {
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

func (wrap *messagingChannel) UnmarshalYAML(n *yaml.Node) (err error) {
	if !y7s.IsKind(n, yaml.MappingNode) {
		return y7s.NodeErr(n, "messagingChannel definition must be a map")
	}

	if wrap.res == nil {
		wrap.res = &types.Channel{}
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
	if wrap.us, err = decodeUserstamps(n); err != nil {
		return
	}

	return nil
}

func (wrap messagingChannel) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewMessagingChannel(wrap.res)
	rs.SetTimestamps(wrap.ts)
	rs.SetUserstamps(wrap.us)
	rs.SetConfig(wrap.envoyConfig)
	return envoy.CollectNodes(
		rs,
		wrap.rbac.bindResource(rs),
	)
}
