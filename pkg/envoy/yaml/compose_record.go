package yaml

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	composeRecord struct {
		// res *types.Record `yaml:",inline"`
		values    map[string]string
		sysValues map[string]string

		refModule    string
		refNamespace string
		// createdBy, updatedBy, deletedBy, ownedBy
		refUser map[string]string
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

	type (
		rw struct {
			rr     resource.ComposeRecordRawSet
			nsRef  string
			modRef string
		}
	)

	// moduleRef to values set
	recMap := make(map[string]*rw)

	for _, res := range wset {
		if recMap[res.refModule] == nil {
			recMap[res.refModule] = &rw{
				rr:     make(resource.ComposeRecordRawSet, 0, 10),
				nsRef:  res.refNamespace,
				modRef: res.refModule,
			}
		}

		r := &resource.ComposeRecordRaw{
			// @todo change this probably
			ID: res.values["id"],

			Values:    res.values,
			RefUser:   res.refUser,
			SysValues: res.sysValues,
		}
		recMap[res.refModule].rr = append(recMap[res.refModule].rr, r)
	}

	for _, w := range recMap {
		walker := func(f func(r *resource.ComposeRecordRaw) error) error {
			for _, r := range w.rr {
				err := f(r)
				if err != nil {
					return err
				}
			}
			return nil
		}

		n := resource.NewComposeRecordSet(walker, w.nsRef, w.modRef)
		for _, r := range w.rr {
			n.IDMap[r.ID] = 0
		}

		nn = append(nn, n)
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

func (wrap *composeRecord) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.refUser == nil {
		wrap.refUser = make(map[string]string)
	}
	if wrap.values == nil {
		wrap.values = make(map[string]string)
	}
	if wrap.sysValues == nil {
		wrap.sysValues = make(map[string]string)
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
			if err := v.Decode(&wrap.values); err != nil {
				return err
			}
			return nil

		case "createdAt":
			return v.Decode(wrap.sysValues["createdAt"])
		case "updatedAt":
			return v.Decode(wrap.sysValues["updatedAt"])
		case "deletedAt":
			return v.Decode(wrap.sysValues["deletedAt"])
		case "createdBy":
			return v.Decode(wrap.refUser["createdBy"])
		case "updatedBy":
			return v.Decode(wrap.refUser["updatedBy"])
		case "deletedBy":
			return v.Decode(wrap.refUser["deletedBy"])
		case "ownedBy":
			return v.Decode(wrap.refUser["ownedBy"])

		}

		return nil
	})
}
