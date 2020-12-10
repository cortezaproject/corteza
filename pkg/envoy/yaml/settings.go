package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	settings struct {
		res    map[string]interface{}
		ts     *resource.Timestamps
		config *resource.EnvoyConfig

		us *resource.Userstamps
	}
)

func (wrap *settings) UnmarshalYAML(n *yaml.Node) (err error) {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "role definition must be a map")
	}

	if wrap.res == nil {
		wrap.res = make(map[string]interface{})
	}

	if err = n.Decode(&wrap.res); err != nil {
		return
	}

	if wrap.config, err = decodeEnvoyConfig(n); err != nil {
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

func (wrap settings) MarshalEnvoy() (nn []resource.Interface, err error) {
	n := resource.NewSettings(wrap.res)
	n.SetTimestamps(wrap.ts)
	n.SetUserstamps(wrap.us)
	n.SetConfig(wrap.config)

	return []resource.Interface{n}, nil
}
