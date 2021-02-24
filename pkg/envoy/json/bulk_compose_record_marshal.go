package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	bulkComposeRecordEncoder struct {
		encoderConfig *EncoderConfig

		res *resource.ComposeRecord

		refNs string

		refMod string
		relMod *types.Module
	}
)

func bulkComposeRecordEncoderFromResource(rec *resource.ComposeRecord, cfg *EncoderConfig) *bulkComposeRecordEncoder {
	return &bulkComposeRecordEncoder{
		encoderConfig: cfg,

		res: rec,
	}
}

func (n *bulkComposeRecordEncoder) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	_, ok := state.Res.(*resource.ComposeRecord)
	if !ok {
		return encoderErrInvalidResource(resource.COMPOSE_RECORD_RESOURCE_TYPE, state.Res.ResourceType())
	}

	return nil
}

func (n *bulkComposeRecordEncoder) Encode(ctx context.Context, w io.Writer, state *envoy.ResourceState) (err error) {
	enc := json.NewEncoder(w)
	all := len(n.encoderConfig.Fields) == 0

	err = n.res.Walker(func(r *resource.ComposeRecordRaw) error {
		m, err := makeMap(
			"id", r.ID,
		)

		r.Ts, err = r.Ts.Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
		if err != nil {
			return err
		}

		// Encode system values
		sysValues := make(map[string]interface{})
		sysValues["createdAt"] = r.Ts.CreatedAt
		sysValues["updatedAt"] = r.Ts.UpdatedAt
		sysValues["deletedAt"] = r.Ts.DeletedAt
		sysValues["createdBy"] = n.res.UserFlakes.GetByStamp(r.Us.CreatedBy)
		sysValues["updatedBy"] = n.res.UserFlakes.GetByStamp(r.Us.UpdatedBy)
		sysValues["deletedBy"] = n.res.UserFlakes.GetByStamp(r.Us.DeletedBy)
		sysValues["ownedBy"] = n.res.UserFlakes.GetByStamp(r.Us.OwnedBy)

		for k, v := range sysValues {
			// Skip fields we don't want
			if !all && !n.encoderConfig.Fields[k] {
				continue
			}
			m, err = addMap(m,
				k, v,
			)

			if err != nil {
				return err
			}
		}

		// Encode record values
		for k, v := range r.Values {
			// Skip fields we don't want
			if !all && !n.encoderConfig.Fields[k] {
				continue
			}

			f := n.res.RelMod.Fields.FindByName(k)
			if f == nil {
				return fmt.Errorf("field %s not found", k)
			}
			if f.Kind == "User" {
				m, err = addMap(m,
					k, n.res.UserFlakes.GetByKey(v),
				)
			} else {
				m, err = addMap(m,
					k, v,
				)
			}
			if err != nil {
				return err
			}
		}

		return enc.Encode(m)
	})

	if err != nil {
		return err
	}
	return nil
}
