package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/node"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"gopkg.in/yaml.v3"
)

type (
	composeModule struct {
		res          *types.Module
		refNamespace string
		rbac         *rbacRules
	}
	composeModuleSet []*composeModule

	composeModuleField struct {
		res  *types.ModuleField `yaml:",inline"`
		rbac *rbacRules
	}
	composeModuleFieldSet []*composeModuleField
)

func (wset *composeModuleSet) UnmarshalYAML(n *yaml.Node) error {
	return each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composeModule{}
		)

		if v == nil {
			return nodeErr(n, "malformed module definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if err = decodeRef(k, "module", &wrap.res.Handle); err != nil {
			return nodeErr(n, "Chart reference must be a valid handle")
		}

		if wrap.res.Name == "" {
			// if name is not set, use handle
			wrap.res.Name = wrap.res.Handle
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wset composeModuleSet) setNamespaceRef(ref string) error {
	for _, res := range wset {
		if res.refNamespace != "" && ref != res.refNamespace {
			return fmt.Errorf("cannot override namespace reference %s with %s", res.refNamespace, ref)
		}

		res.refNamespace = ref
	}

	return nil
}

func (wset composeModuleSet) MarshalEnvoy() ([]envoy.Node, error) {
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

func (wrap *composeModule) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.rbac = &rbacRules{}
		wrap.res = &types.Module{}
	}

	if wrap.rbac, err = decodeResourceAccessControl(types.ModuleRBACResource, n); err != nil {
		return
	}

	return eachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "name":
			return decodeScalar(v, "module name", &wrap.res.Name)

		case "handle":
			return decodeScalar(v, "module handle", &wrap.res.Handle)

		case "fields":
			if !isKind(v, yaml.MappingNode) {
				return nodeErr(n, "field definition must be a map")
			}

			var (
				aux = composeModuleFieldSet{}
			)

			if err = v.Decode(&aux); err != nil {
				return err
			}

			wrap.res.Fields = aux.set()
			return nil

		}

		return nil
	})
}

func (wrap composeModule) MarshalEnvoy() ([]envoy.Node, error) {
	return envoy.CollectNodes(
		&node.ComposeModule{
			Res:          wrap.res,
			RefNamespace: wrap.refNamespace,
		},
		wrap.rbac.Ensure(),
	)
}

func (set *composeModuleFieldSet) UnmarshalYAML(n *yaml.Node) error {
	return each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composeModuleField{}
		)

		if v == nil {
			return nodeErr(n, "malformed module field definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return err
		}

		if k != nil {
			if !handle.IsValid(k.Value) {
				return nodeErr(n, "field name must be a valid handle")
			}

			wrap.res.Name = k.Value
			wrap.res.Label = k.Value
		}

		*set = append(*set, wrap)
		return
	})
}

func (set composeModuleFieldSet) set() (out types.ModuleFieldSet) {
	for _, i := range set {
		out = append(out, i.res)
	}

	return out
}

func (wrap *composeModuleField) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.res = &types.ModuleField{}
	}

	if err = n.Decode(&wrap.res); err != nil {
		return
	}

	if wrap.rbac, err = decodeResourceAccessControl(types.ModuleFieldRBACResource, n); err != nil {
		return
	}

	return eachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "name":
			return fmt.Errorf("name should be encoded as field definition key")

		case "default":
			var rvs = types.RecordValueSet{}
			switch v.Kind {
			case yaml.ScalarNode:
				rvs = rvs.Set(&types.RecordValue{Value: v.Value})

			case yaml.SequenceNode:
				_ = eachSeq(v, func(v *yaml.Node) error {
					rvs = rvs.Set(&types.RecordValue{Value: v.Value, Place: uint(len(rvs))})
					return nil
				})
			}

			wrap.res.DefaultValue = rvs
		}

		return nil
	})
}
