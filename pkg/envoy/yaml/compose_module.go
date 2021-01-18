package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	mappingTpl struct {
		resource.MappingTpl `yaml:",inline"`
	}
	mappingTplSet []*mappingTpl

	composeRecordTpl struct {
		Source string `yaml:"from"`

		Key     []string
		Mapping mappingTplSet
	}

	composeModule struct {
		res    *types.Module
		ts     *resource.Timestamps
		config *resource.EnvoyConfig

		refNamespace string
		rbac         rbacRuleSet

		recTpl *composeRecordTpl
	}
	composeModuleSet []*composeModule

	composeModuleField struct {
		res  *types.ModuleField
		rbac rbacRuleSet
	}
	composeModuleFieldSet []*composeModuleField

	// aux struct for decoding module field expressions
	composeModuleFieldExprAux types.ModuleFieldExpr
)

func (wset *composeModuleSet) UnmarshalYAML(n *yaml.Node) error {
	return Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composeModule{}
		)

		if v == nil {
			return NodeErr(n, "malformed module definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if err = decodeRef(k, "module", &wrap.res.Handle); err != nil {
			return NodeErr(n, "Chart reference must be a valid handle")
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

func (wset composeModuleSet) MarshalEnvoy() ([]resource.Interface, error) {
	// namespace usually have bunch of sub-resources defined
	nn := make([]resource.Interface, 0, len(wset)*10)

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
		wrap.rbac = make(rbacRuleSet, 0, 10)
		wrap.res = &types.Module{}
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

	return EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "name":
			return DecodeScalar(v, "module name", &wrap.res.Name)

		case "handle":
			return DecodeScalar(v, "module handle", &wrap.res.Handle)

		case "fields":
			if !IsKind(v, yaml.MappingNode) {
				return NodeErr(n, "field definition must be a map")
			}

			var (
				aux = composeModuleFieldSet{}
			)

			if err = v.Decode(&aux); err != nil {
				return err
			}

			wrap.res.Fields = aux.set()
			return nil

		case "records":
			if IsKind(v, yaml.MappingNode) {
				return v.Decode(&wrap.recTpl)
			} else {
				return NodeErr(n, "records definition must be a map")
			}

		}

		return nil
	})
}

func (wrap composeModule) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewComposeModule(wrap.res, wrap.refNamespace)
	rs.SetTimestamps(wrap.ts)
	rs.SetConfig(wrap.config)

	var crs *resource.ComposeRecordTemplate
	if wrap.recTpl != nil {
		s := wrap.recTpl

		mtt := make(resource.MappingTplSet, 0, len(s.Mapping))
		for _, m := range s.Mapping {
			mtt = append(mtt, &m.MappingTpl)
		}
		crs = resource.NewComposeRecordTemplate(rs.Identifiers().First(), wrap.refNamespace, s.Source, mtt, s.Key...)
	}

	return envoy.CollectNodes(
		rs,
		crs,
		wrap.rbac.bindResource(rs),
	)
}

func (rt *composeRecordTpl) UnmarshalYAML(n *yaml.Node) error {
	return EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "source", "origin", "from":
			rt.Source = v.Value

		case "key", "index", "pk":
			if !IsKind(v, yaml.SequenceNode) {
				rt.Key = []string{v.Value}
			} else {
				rt.Key = make([]string, 0, 3)
				EachSeq(v, func(v *yaml.Node) error {
					rt.Key = append(rt.Key, v.Value)
					return nil
				})
			}

		case "mapping", "map":
			rt.Mapping = make(mappingTplSet, 0, 20)
			// When provided as a sequence node, map the fields based on the index.
			// first cell is mapped to the first sequence value, second cell to the second, and so on.
			// Omit cells with empty, /, or - value.
			if IsKind(v, yaml.SequenceNode) {
				i := uint(0)
				EachSeq(v, func(v *yaml.Node) error {
					defer func() { i++ }()

					if v.Value == "" || v.Value == "/" || v.Value == "-" {
						return nil
					}

					rt.Mapping = append(rt.Mapping, &mappingTpl{
						MappingTpl: resource.MappingTpl{
							Index: i,
							Field: v.Value,
						},
					})
					return nil
				})
			} else if IsKind(v, yaml.MappingNode) {
				// When provided as a mapping node, it can be a simple cell: field map
				// or a more complex underlying structure.
				EachMap(v, func(k, v *yaml.Node) error {
					m := &mappingTpl{}

					if IsKind(v, yaml.MappingNode) {
						v.Decode(m)
					} else {
						m.Field = v.Value
					}

					m.Cell = k.Value
					rt.Mapping = append(rt.Mapping, m)
					return nil
				})
			}
		}

		return nil
	})
}

func (set *composeModuleFieldSet) UnmarshalYAML(n *yaml.Node) error {
	return Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composeModuleField{}
		)

		if v == nil {
			return NodeErr(n, "malformed module field definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return err
		}

		if k != nil {
			if !handle.IsValid(k.Value) {
				return NodeErr(n, "field name must be a valid handle")
			}

			wrap.res.Name = k.Value
			if wrap.res.Label == "" {
				wrap.res.Label = k.Value
			}
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
		wrap.rbac = make(rbacRuleSet, 0, 10)
		wrap.res = &types.ModuleField{}
	}

	if wrap.rbac, err = decodeRbac(n); err != nil {
		return
	}

	return EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "name":
			return fmt.Errorf("name should be encoded as field definition key")

		case "place":
			return DecodeScalar(v, "module place", &wrap.res.Place)

		case "kind":
			return DecodeScalar(v, "module kind", &wrap.res.Kind)

		case "label":
			return DecodeScalar(v, "module label", &wrap.res.Label)

		case "private":
			return DecodeScalar(v, "module private", &wrap.res.Private)

		case "required":
			return DecodeScalar(v, "module required", &wrap.res.Required)

		case "visible":
			return DecodeScalar(v, "module visible", &wrap.res.Visible)

		case "multi":
			return DecodeScalar(v, "module multi", &wrap.res.Multi)

		case "options":
			if err = v.Decode(&wrap.res.Options); err != nil {
				return err
			}

		case "expressions":
			ea := composeModuleFieldExprAux{}
			if err = v.Decode(&ea); err != nil {
				return err
			}
			wrap.res.Expressions = types.ModuleFieldExpr(ea)

		case "default":
			var rvs = types.RecordValueSet{}
			switch v.Kind {
			case yaml.ScalarNode:
				rvs = rvs.Set(&types.RecordValue{Value: v.Value})

			case yaml.SequenceNode:
				_ = EachSeq(v, func(v *yaml.Node) error {
					rvs = rvs.Set(&types.RecordValue{Value: v.Value, Place: uint(len(rvs))})
					return nil
				})
			}

			wrap.res.DefaultValue = rvs
		}

		return nil
	})
}

func (aux *composeModuleFieldExprAux) UnmarshalYAML(n *yaml.Node) (err error) {
	return EachMap(n, func(k *yaml.Node, v *yaml.Node) error {
		switch k.Value {
		case "valueExpr", "value":
			aux.ValueExpr = v.Value
			return nil
		case "sanitizer":
			aux.Sanitizers = []string{v.Value}
			return nil
		case "sanitizers":
			return EachSeq(v, func(san *yaml.Node) error {
				aux.Sanitizers = append(aux.Sanitizers, san.Value)
				return nil
			})
		case "validator", "validators":
			return Each(v, func(k *yaml.Node, v *yaml.Node) error {
				vld := &types.ModuleFieldValidator{}
				if IsKind(v, yaml.MappingNode) {
					if err := v.Decode(vld); err != nil {
						return err
					}
				} else {
					vld.Test = k.Value
					vld.Error = v.Value
				}

				aux.Validators = append(aux.Validators, *vld)
				return nil
			})
		case "disableDefaultValidators":
			return v.Decode(&aux.DisableDefaultValidators)
		case "formatter":
			aux.Formatters = []string{v.Value}
			return nil
		case "formatters":
			return EachSeq(v, func(fmt *yaml.Node) error {
				aux.Formatters = append(aux.Formatters, fmt.Value)
				return nil
			})
		case "disableDefaultFormatters":
			return v.Decode(&aux.DisableDefaultFormatters)
		}

		return nil
	})
}
