package yaml

import (
	"context"
	"fmt"

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
		return encoderErrInvalidResource(types.RecordResourceType, state.Res.ResourceType())
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
				r.Values[k], err = n.res.UserFlakes.GetByKey(v).Model()
				if err != nil {
					return err
				}
			}
		}

		cr.values = vv
		doc.addComposeRecord(cr)
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

	nn, err = encodeTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	nn, err = encodeUserstamps(nn, c.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}
