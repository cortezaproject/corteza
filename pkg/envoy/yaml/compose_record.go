package yaml

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	composeRecord struct {
		values map[string]string
		ts     *resource.Timestamps
		us     *resource.Userstamps
		config *resource.EnvoyConfig

		refModule    string
		refNamespace string
	}
	composeRecordSet []*composeRecord

	composeRecordValues struct {
		rvs types.RecordValueSet
	}
)

// UnmarshalYAML resolves set of record definitions, either sequence or map
//
// When resolving map, key is used as module handle
func (wset *composeRecordSet) UnmarshalYAML(n *yaml.Node) error {
	return Each(n, func(k, v *yaml.Node) (err error) {
		var (
			moduleRef string
		)

		if v == nil {
			return NodeErr(n, "malformed record definition")
		}

		if err = decodeRef(k, "module", &moduleRef); err != nil {
			return
		}

		if IsKind(v, yaml.SequenceNode) {
			// multiple records defined
			return EachSeq(v, func(r *yaml.Node) error {
				var wrap = &composeRecord{refModule: moduleRef}
				if err = r.Decode(&wrap); err != nil {
					return err
				}

				*wset = append(*wset, wrap)
				return nil
			})
		}

		if IsKind(v, yaml.MappingNode) {
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
			rr      resource.ComposeRecordRawSet
			nsRef   string
			modRef  string
			refUser resource.Identifiers
		}
	)

	// We'll do a list of wrappers & a map of wrappers to preserve order and keep
	// optimal lookups
	rww := make([]*rw, 0, len(wset))
	rrx := make(map[string]*rw)

	for _, res := range wset {
		// A bit stronger index just in case
		ix := res.refNamespace + "/" + res.refModule
		if _, ok := rrx[ix]; !ok {
			rrx[ix] = &rw{
				rr:      make(resource.ComposeRecordRawSet, 0, 10),
				nsRef:   res.refNamespace,
				modRef:  res.refModule,
				refUser: make(resource.Identifiers),
			}
			rww = append(rww, rrx[ix])
		}

		r := &resource.ComposeRecordRaw{
			// @todo change this probably
			ID:     res.values["id"],
			Config: res.config,
			Values: res.values,
			Ts:     res.ts,
			Us:     res.us,
		}

		rrx[ix].refUser.Add(
			res.us.CreatedBy,
			res.us.UpdatedBy,
			res.us.DeletedBy,
			res.us.OwnedBy,
		)

		rrx[ix].rr = append(rrx[ix].rr, r)
	}

	for _, w := range rww {
		cw := *w

		walker := func(f func(r *resource.ComposeRecordRaw) error) error {
			for _, r := range cw.rr {
				err := f(r)
				if err != nil {
					return err
				}
			}
			return nil
		}

		n := resource.NewComposeRecordSet(walker, w.nsRef, w.modRef)
		n.SetUserRefs(w.refUser.StringSlice())
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
	if wrap.values == nil {
		wrap.values = make(map[string]string)
	}

	// @todo enable when records are ready for RBAC
	//if wrap.rbac, err = decodeRbac(types.RecordRBACResource, n); err != nil {
	//	return
	//}

	if wrap.config, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.ts, err = decodeTimestamps(n); err != nil {
		return
	}
	if wrap.us, err = decodeUserstamps(n); err != nil {
		return
	}

	return EachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "module":
			return decodeRef(v, "module", &wrap.refModule)

		case "values":
			// Use aux structure to decode record values into RVS
			if err := v.Decode(&wrap.values); err != nil {
				return err
			}
			return nil

		}

		return nil
	})
}
