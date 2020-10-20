package yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"gopkg.in/yaml.v3"
)

type (
	ComposeRecord struct {
		res          *types.Record       `yaml:",inline"`
		values       ComposeRecordValues `yaml:"values"`
		moduleRef    string
		createdByRef string
		updatedByRef string
		deletedByRef string
		ownedByRef   string
	}
	ComposeRecordSet []*ComposeRecord

	ComposeRecordValues struct {
		rvs types.RecordValueSet
	}
)

// UnmarshalYAML resolves set of record definitions, either sequence or map
//
// When resolving map, key is used as module handle
//
// { module-handle: [ { ... values ... } ] }
// [ { module: module-handle, ... values ... } ]
func (wset *ComposeRecordSet) UnmarshalYAML(n *yaml.Node) error {
	return iterator(n, func(k, v *yaml.Node) (err error) {
		var (
			moduleRef string
		)

		if k != nil {
			// processing mapping node, expecting module handle
			if !handle.IsValid(k.Value) {
				return nodeErr(k, "module reference must be a valid handle")
			}

			moduleRef = k.Value
		}

		if v == nil {
			return nodeErr(n, "malformed record definition")
		}

		if isKind(v, yaml.SequenceNode) {
			// multiple records defined
			return iterator(v, func(_, r *yaml.Node) error {
				var wrap = &ComposeRecord{moduleRef: moduleRef}
				if err = r.Decode(&wrap); err != nil {
					return err
				}

				*wset = append(*wset, wrap)
				return nil
			})
		}

		if isKind(v, yaml.MappingNode) {
			// one record defined
			var wrap = &ComposeRecord{moduleRef: moduleRef}
			if err = v.Decode(&wrap); err != nil {
				return
			}

			*wset = append(*wset, wrap)
		}

		return nil
	})
}

func (wrap *ComposeRecord) UnmarshalYAML(n *yaml.Node) error {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "expecting mapping node for record definition")
	}

	if wrap.res == nil {
		wrap.res = &types.Record{}
	}

	return iterator(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "module":
			return decodeRef(v, "module", &wrap.moduleRef)

		case "values":
			// Use aux structure to decode record values into RVS
			aux := ComposeRecordValues{}
			if err := v.Decode(&aux); err != nil {
				return err
			}

			wrap.res.Values = aux.rvs
			return nil

		case "createdAt":
			return v.Decode(&wrap.res.CreatedAt)
		case "updatedAt":
			return v.Decode(&wrap.res.UpdatedAt)
		case "deletedAt":
			return v.Decode(&wrap.res.DeletedAt)
		case "createdBy":
			return decodeRef(v, "user", &wrap.createdByRef)
		case "updatedBy":
			return decodeRef(v, "user", &wrap.updatedByRef)
		case "deletedBy":
			return decodeRef(v, "user", &wrap.deletedByRef)
		case "ownedBy":
			return decodeRef(v, "user", &wrap.ownedByRef)

		default:
			return nodeErr(k, "unsupported key %s used for record definition", k.Value)
		}

		return nil
	})
}

// UnmarshalYAML resolves record values definitioons
//
// { <field name>: ... <scalar value>, .... }
// { <field name>: [ <scalar value> ], .... }
func (wset *ComposeRecordValues) UnmarshalYAML(n *yaml.Node) error {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "expecting mapping node for record value definition")
	}

	wset.rvs = types.RecordValueSet{}

	return iterator(n, func(k, v *yaml.Node) error {
		if isKind(v, yaml.ScalarNode) {
			wset.rvs = append(wset.rvs, &types.RecordValue{
				Name:  k.Value,
				Value: v.Value,
			})

			return nil
		}

		if isKind(v, yaml.SequenceNode) {
			for i := range v.Content {
				if isKind(v, yaml.ScalarNode) {
					return nodeErr(n, "expecting scalar node for record value")
				}

				wset.rvs = append(wset.rvs, &types.RecordValue{
					Name:  k.Value,
					Value: v.Content[i].Value,
					Place: uint(i),
				})
			}

			return nil
		}

		return nodeErr(n, "expecting scalar or sequence node for record value")
	})
}
