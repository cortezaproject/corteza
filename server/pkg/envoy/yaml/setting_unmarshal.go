package yaml

import (
	"encoding/json"

	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
	sqlt "github.com/jmoiron/sqlx/types"
	"gopkg.in/yaml.v3"
)

func (wset *settingSet) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &setting{}
		)

		if v == nil {
			return y7s.NodeErr(n, "malformed setting definition")
		}

		wrap.res = &types.SettingValue{}

		switch v.Kind {
		case yaml.MappingNode:
			if err = v.Decode(&wrap); err != nil {
				return
			}

		default:
			if err = y7s.DecodeScalar(k, "setting name", &wrap.res.Name); err != nil {
				return err
			}

			if y7s.IsKind(v, yaml.SequenceNode) {
				aux := make([]interface{}, 0, 10)

				y7s.EachSeq(v, func(n *yaml.Node) error {
					var vx interface{}
					err := y7s.DecodeScalar(n, "setting value", &vx)
					if err != nil {
						return err
					}

					aux = append(aux, vx)
					return nil
				})

				m, err := json.Marshal(aux)
				if err != nil {
					return err
				}
				wrap.res.Value = sqlt.JSONText(m)
			} else if y7s.IsKind(v, yaml.MappingNode) {
				aux := make(map[string]interface{})

				y7s.EachMap(v, func(k, v *yaml.Node) error {
					var vx interface{}
					err := y7s.DecodeScalar(v, "setting value", &vx)
					if err != nil {
						return err
					}

					aux[k.Value] = vx
					return nil
				})

				m, err := json.Marshal(aux)
				if err != nil {
					return err
				}
				wrap.res.Value = sqlt.JSONText(m)
			} else {
				var aux interface{}
				err = y7s.DecodeScalar(v, "setting value", &aux)
				if err != nil {
					return err
				}
				m, err := json.Marshal(aux)
				if err != nil {
					return err
				}
				wrap.res.Value = sqlt.JSONText(m)
			}
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wrap *setting) UnmarshalYAML(n *yaml.Node) (err error) {
	if !y7s.IsKind(n, yaml.MappingNode) {
		return y7s.NodeErr(n, "setting definition must be a map")
	}

	if wrap.res == nil {
		wrap.res = &types.SettingValue{}
	}

	err = y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "name":
			return y7s.DecodeScalar(v, "setting name", &wrap.res.Name)

		case "value":
			var aux interface{}
			err = y7s.DecodeScalar(v, "setting value", &aux)
			if err != nil {
				return err
			}
			m, err := json.Marshal(aux)
			if err != nil {
				return err
			}
			wrap.res.Value = sqlt.JSONText(m)
		}

		return nil
	})

	if err != nil {
		return err
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

func (wset settingSet) MarshalEnvoy() ([]resource.Interface, error) {
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

func (wrap *setting) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewSetting(wrap.res)
	rs.SetTimestamps(wrap.ts)
	rs.SetUserstamps(wrap.us)
	rs.SetConfig(wrap.envoyConfig)

	return envoy.CollectNodes(
		rs,
	)
}
