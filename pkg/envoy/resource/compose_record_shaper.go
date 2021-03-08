package resource

type (
	// A simple wrapper struct for related resources
	composeRecordShaper struct{}
)

// ComposeRecordShaper initializes and returns a new compose record resource shaper
func ComposeRecordShaper() shaper {
	return &composeRecordShaper{}
}

// Shape shapes ResourceDatasets based on the related compose template
//
// The first step finds the matching pair;
// The second step creates an actual resource based on the two.
func (crt *composeRecordShaper) Shape(rr []Interface) ([]Interface, error) {
	// This will do for most cases
	ii := make([]Interface, 0, int(len(rr)/2)+1)

	for _, r := range rr {
		rt, ok := r.(*ComposeRecordTemplate)
		if !ok {
			continue
		}
		rd := findResourceDataset(rr, rt.Identifiers())
		if rd == nil {
			return nil, genericSourceErrUnresolved(rt.Identifiers())
		}

		ii = append(ii, crt.toResource(rt, rd))
	}

	return ii, nil
}

func (crt *composeRecordShaper) toResource(def *ComposeRecordTemplate, dt *ResourceDataset) Interface {
	w := func(f func(r *ComposeRecordRaw) error) error {
		for {
			mr, err := dt.P.Next()
			if err != nil {
				return err
			}
			if mr == nil {
				return nil
			}

			// Get the bits in order
			rRaw := &ComposeRecordRaw{
				ID:     crt.getKey(mr, def.Key),
				Values: crt.mapValues(mr, def.FieldMap),
			}

			// Process it
			err = f(rRaw)
			if err != nil {
				return err
			}
		}
	}

	return NewComposeRecordSet(w, def.NsRef.Identifiers.First(), def.ModRef.Identifiers.First())
}

// mapValues maps original values based on the provided mapping
func (crt *composeRecordShaper) mapValues(ov map[string]string, fm MappingTplSet) map[string]string {
	nv := make(map[string]string)

	// Mappings are provided as a slice since we'll want to do some additional conditional mappings.
	// We'll make an index for a nicer lookup for now.
	mx := make(map[string]*MappingTpl)
	for _, m := range fm {
		mx[m.Cell] = m
	}

	if len(mx) > 0 {
		for k, v := range ov {
			if m, has := mx[k]; has {
				nv[m.Field] = v
			}
		}
	} else {
		for k, v := range ov {
			nv[k] = v
		}
	}

	return nv
}

func (crt *composeRecordShaper) getKey(vv map[string]string, kk []string) (rtr string) {
	if len(kk) <= 0 {
		return ""
	}

	for _, k := range kk {
		rtr += k + "."
	}
	// Remove the trailing dot
	return vv[rtr[0:len(rtr)-1]]
}
