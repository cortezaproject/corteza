package yaml

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	composeRecord struct {
		res          *types.Record `yaml:",inline"`
		refModule    string
		refNamespace string
		refCreatedBy string
		refUpdatedBy string
		refDeletedBy string
		refOwnedBy   string

		rbac rbacRuleSet
	}
	composeRecordSet []*composeRecord

	composeRecordValues struct {
		rvs types.RecordValueSet
	}
)

// UnmarshalYAML resolves set of record definitions, either sequence or map
//
// When resolving map, key is used as module handle
//
// { module-handle: [ { ... values ... } ] }
// [ { module: module-handle, ... values ... } ]
func (wset *composeRecordSet) UnmarshalYAML(n *yaml.Node) error {
	return each(n, func(k, v *yaml.Node) (err error) {
		var (
			moduleRef string
		)

		if v == nil {
			return nodeErr(n, "malformed record definition")
		}

		if err = decodeRef(k, "module", &moduleRef); err != nil {
			return
		}

		if isKind(v, yaml.SequenceNode) {
			// multiple records defined
			return eachSeq(v, func(r *yaml.Node) error {
				var wrap = &composeRecord{refModule: moduleRef}
				if err = r.Decode(&wrap); err != nil {
					return err
				}

				*wset = append(*wset, wrap)
				return nil
			})
		}

		if isKind(v, yaml.MappingNode) {
			// one record defined
			var wrap = &composeRecord{refModule: moduleRef}
			if err = v.Decode(&wrap); err != nil {
				return
			}

			*wset = append(*wset, wrap)
		}

		return nil
	})
}

// MarshalEnvoy works a bit differenlty
func (wset composeRecordSet) MarshalEnvoy() ([]resource.Interface, error) {
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

func (wset composeRecordSet) setNamespaceRef(ref string) error {
	for _, res := range wset {
		if res.refNamespace != "" && ref != res.refNamespace {
			return fmt.Errorf("cannot override namespace reference %s with %s", res.refNamespace, ref)
		}

		res.refNamespace = ref
	}

	return nil
}

func (wrap composeRecord) MarshalEnvoy() ([]resource.Interface, error) {
	var refUsers = []string{}
	for _, u := range []string{wrap.refCreatedBy, wrap.refUpdatedBy, wrap.refDeletedBy, wrap.refOwnedBy} {
		if len(u) > 0 {
			// emails and handles
			refUsers = append(refUsers, u)
		}
	}

	return envoy.CollectNodes(
		// @todo...
		resource.NewComposeRecordSet(),
		wrap.rbac,
	)
}

func (wrap *composeRecord) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.res = &types.Record{}
	}

	// @todo enable when records are ready for RBAC
	//if wrap.rbac, err = decodeRbac(types.RecordRBACResource, n); err != nil {
	//	return
	//}

	return eachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "module":
			return decodeRef(v, "module", &wrap.refModule)

		case "values":
			// Use aux structure to decode record values into RVS
			aux := composeRecordValues{}
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
			return decodeRef(v, "createdBy user", &wrap.refCreatedBy)
		case "updatedBy":
			return decodeRef(v, "updatedBy user", &wrap.refUpdatedBy)
		case "deletedBy":
			return decodeRef(v, "deletedBy user", &wrap.refDeletedBy)
		case "ownedBy":
			return decodeRef(v, "ownedBy user", &wrap.refOwnedBy)

		}

		return nil
	})
}

// UnmarshalYAML resolves record values definitioons
//
// { <field name>: ... <scalar value>, .... }
// { <field name>: [ <scalar value> ], .... }
func (wset *composeRecordValues) UnmarshalYAML(n *yaml.Node) error {
	wset.rvs = types.RecordValueSet{}

	return eachMap(n, func(k, v *yaml.Node) error {
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
