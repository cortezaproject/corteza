package yaml

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	composeRecordEncoder struct {
		encoderConfig *EncoderConfig

		res *resource.ComposeRecord

		refNamespace string
		relNs        *types.Namespace
		refModule    string
		relMod       *types.Module
	}
)

func composeRecordSetFromResource(rec *resource.ComposeRecord, cfg *EncoderConfig) *composeRecordEncoder {
	return &composeRecordEncoder{
		res:           rec,
		encoderConfig: cfg,
	}
}

func (n *composeRecordEncoder) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	rec, ok := state.Res.(*resource.ComposeRecord)
	if !ok {
		return errors.New("invalid type; @todo error")
	}

	// Get the related namespace
	n.relNs = resource.FindComposeNamespace(state.ParentResources, rec.RefNs.Identifiers)
	if n.relNs == nil {
		return resource.ComposeNamespaceErrUnresolved(rec.RefNs.Identifiers)
	}

	// Get the related module
	n.relMod = resource.FindComposeModule(state.ParentResources, rec.RefMod.Identifiers)
	if n.relNs == nil {
		return resource.ComposeNamespaceErrUnresolved(rec.RefMod.Identifiers)
	}

	n.refNamespace = relNsToRef(n.relNs)
	n.refModule = relModToRef(n.relMod)

	return nil
}

func (n *composeRecordEncoder) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	return n.res.Walker(func(r *resource.ComposeRecordRaw) error {
		cr := &composeRecord{
			us:           r.Us,
			config:       r.Config,
			refModule:    n.refModule,
			refNamespace: n.refNamespace,
		}

		cr.us, err = resolveUserstamps(state.ParentResources, r.Us)
		if err != nil {
			return err
		}
		cr.ts, err = r.Ts.Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
		if err != nil {
			return err
		}

		vv := r.Values
		for k, v := range vv {
			f := n.res.RelMod.Fields.FindByName(k)
			if f == nil {
				return fmt.Errorf("field %s not found", k)
			}

			// Resolve user references
			if f.Kind == "User" {
				r.Values[k], err = n.res.UserFlakes.GetByKey(v).Stringify()
				if err != nil {
					return err
				}
			}
		}

		cr.values = vv
		doc.AddComposeRecord(cr)
		return nil
	})
}

func (rr composeRecordSet) MarshalYAML() (interface{}, error) {
	byMod := make(map[string]composeRecordSet)

	for _, r := range rr {
		if _, has := byMod[r.refModule]; !has {
			byMod[r.refModule] = make(composeRecordSet, 0, 100)
		}
		byMod[r.refModule] = append(byMod[r.refModule], r)
	}

	nn, _ := makeMap()
	for refMod, rr := range byMod {
		ss, err := makeSeq()
		for _, r := range rr {
			ss, err = addSeq(ss, r)
			if err != nil {
				return nil, err
			}
		}
		nn, err = addMap(nn,
			refMod, ss,
		)
	}

	return nn, nil
}

func (c *composeRecord) MarshalYAML() (interface{}, error) {
	nn, err := makeMap(
		"values", c.values,
	)
	if err != nil {
		return nil, err
	}

	nn, err = mapTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	nn, err = mapUserstamps(nn, c.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func resolveUserstamps(rr []resource.Interface, us *resource.Userstamps) (*resource.Userstamps, error) {
	if us == nil {
		return nil, nil
	}

	fetch := func(us *resource.Userstamp) (*resource.Userstamp, error) {
		if us == nil || us.UserID == 0 {
			return nil, nil
		}

		// This one can be considered as valid
		if us.Ref != "" && us.UserID > 0 && us.U != nil {
			return us, nil
		}

		ii := resource.MakeIdentifiers()

		if us.UserID > 0 {
			ii = ii.Add(strconv.FormatUint(us.UserID, 10))
		}
		if us.Ref != "" {
			ii = ii.Add(us.Ref)
		}
		if us.U != nil {
			if us.U.Handle != "" {
				ii = ii.Add(us.U.Handle)
			}
			if us.U.Username != "" {
				ii = ii.Add(us.U.Username)
			}
			if us.U.Email != "" {
				ii = ii.Add(us.U.Email)
			}
		}

		u := resource.FindUser(rr, ii)
		if u == nil {
			return nil, resource.UserErrUnresolved(ii)
		}

		return resource.MakeUserstamp(u), nil
	}
	var err error
	us.CreatedBy, err = fetch(us.CreatedBy)
	us.UpdatedBy, err = fetch(us.UpdatedBy)
	us.DeletedBy, err = fetch(us.DeletedBy)
	us.OwnedBy, err = fetch(us.OwnedBy)
	us.RunAs, err = fetch(us.RunAs)

	if err != nil {
		return nil, err
	}

	return us, nil
}
