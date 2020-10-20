package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"gopkg.in/yaml.v3"
)

type (
	ComposeModule struct {
		res  *types.Module
		Rbac `yaml:",inline"`
	}
	ComposeModuleSet []*ComposeModule

	ComposeModuleField struct {
		res  *types.ModuleField `yaml:",inline"`
		Rbac `yaml:",inline"`
	}
	ComposeModuleFieldSet []*ComposeModuleField
)

func (set *ComposeModuleSet) UnmarshalYAML(n *yaml.Node) error {
	return iterator(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &ComposeModule{}
		)

		if v == nil {
			return nodeErr(n, "malformed module definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if k != nil {
			if wrap.res.Handle != "" {
				return nodeErr(k, "cannot define handle in mapped module definition")
			}

			if !handle.IsValid(k.Value) {
				return nodeErr(n, "module reference must be a valid handle")
			}

			wrap.res.Handle = k.Value
			wrap.res.Name = k.Value
		}

		*set = append(*set, wrap)
		return
	})
}

func (wrap *ComposeModule) UnmarshalYAML(n *yaml.Node) error {
	if wrap.res == nil {
		wrap.res = &types.Module{}
	}

	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "module definition must be a map")
	}

	return iterator(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "name":
			return decodeScalar(v, &wrap.res.Name)
		case "fields":
			if !isKind(v, yaml.MappingNode) {
				return nodeErr(n, "field definition must be a map")
			}

			var (
				aux = ComposeModuleFieldSet{}
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

func (set *ComposeModuleFieldSet) UnmarshalYAML(n *yaml.Node) error {
	return iterator(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &ComposeModuleField{}
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

func (set ComposeModuleFieldSet) set() (out types.ModuleFieldSet) {
	for _, i := range set {
		out = append(out, i.res)
	}

	return out
}

func (wrap *ComposeModuleField) UnmarshalYAML(n *yaml.Node) error {
	if wrap.res == nil {
		wrap.res = &types.ModuleField{}
	}

	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "module field definition must be a map")
	}

	return iterator(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "name":
			return fmt.Errorf("name should be encoded as field definition key")

		case "label":
			return decodeScalar(v, &wrap.res.Label)

		case "kind", "type":
			return decodeScalar(v, &wrap.res.Kind)

		case "options":
			return v.Decode(&wrap.res.Options)

		case "private":
			return decodeScalar(v, &wrap.res.Private)

		case "required":
			return decodeScalar(v, &wrap.res.Required)

		case "visible":
			return decodeScalar(v, &wrap.res.Visible)

		case "multi":
			return decodeScalar(v, &wrap.res.Multi)

		case "default":
			return fmt.Errorf("field.default /// to be imple,emted")
			//wrap.res.DefaultValue = types.RecordValueSet{}
			//return deinterfacer.Each(val, func(place int, _ string, val interface{}) (err error) {
			//	field.DefaultValue = append(field.DefaultValue, &types.RecordValue{
			//		Value: deinterfacer.ToString(val),
			//		Place: uint(place),
			//	})
			//
			//	return
			//})
		}

		return nil
	})
}
