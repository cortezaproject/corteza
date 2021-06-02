package yaml

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

// UnmarshalYAML resolves set of record definitions, either sequence or map
//
// When resolving map, key is used as module handle
func (wset *composeRecordSet) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k, v *yaml.Node) (err error) {
		var (
			moduleRef string
		)

		if v == nil {
			return y7s.NodeErr(n, "malformed record definition")
		}

		if err = decodeRef(k, "module", &moduleRef); err != nil {
			return
		}

		if y7s.IsKind(v, yaml.SequenceNode) {
			// multiple records defined
			return y7s.EachSeq(v, func(r *yaml.Node) error {
				var wrap = &composeRecord{refModule: moduleRef}
				if err = r.Decode(&wrap); err != nil {
					return err
				}

				*wset = append(*wset, wrap)
				return nil
			})
		}

		if y7s.IsKind(v, yaml.MappingNode) {
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

func (wrap *composeRecord) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.values == nil {
		wrap.values = make(map[string]string)
	}

	// @todo enable when records are ready for RBAC
	//if wrap.rbac, err = decodeRbac(&types.Component{}, n); err != nil {
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

	return y7s.EachMap(n, func(k, v *yaml.Node) error {
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

// MarshalEnvoy works a bit differenlty
func (wset composeRecordSet) MarshalEnvoy() ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, len(wset))

	type (
		auxRecord struct {
			rr     resource.ComposeRecordRawSet
			nsRef  string
			modRef string
		}
	)

	// We'll do a list of wrappers & a map of wrappers to preserve order and keep
	// optimal lookups
	rww := make([]*auxRecord, 0, len(wset))
	rrx := make(map[string]*auxRecord)

	for _, res := range wset {
		// A bit stronger index just in case
		ix := res.refNamespace + "/" + res.refModule
		if _, ok := rrx[ix]; !ok {
			rrx[ix] = &auxRecord{
				rr:     make(resource.ComposeRecordRawSet, 0, 10),
				nsRef:  res.refNamespace,
				modRef: res.refModule,
			}
			rww = append(rww, rrx[ix])
		}

		r := &resource.ComposeRecordRaw{
			ID:     res.getID(),
			Config: res.config,
			Values: res.values,
			Ts:     res.ts,
			Us:     res.us,
		}

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
		// Empty userstamp index so the referencing will work with wildcards
		n.SetUserFlakes(make(resource.UserstampIndex))

		// Take note of the record ID's that are provided.
		// This will later let us find existing records.
		for _, r := range w.rr {
			if r.ID != "" {
				n.IDMap[r.ID] = 0
			}
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

// Utilities

func (wrap *composeRecord) getID() string {
	if wrap.values["id"] != "" {
		return wrap.values["id"]
	} else if wrap.values["ID"] != "" {
		return wrap.values["ID"]
	} else if wrap.values["recordID"] != "" {
		return wrap.values["recordID"]
	}

	return ""
}
