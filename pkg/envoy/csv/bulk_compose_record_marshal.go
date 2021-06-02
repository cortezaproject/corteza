package csv

import (
	"context"
	"encoding/csv"
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

// Prepare prepares the composeRecord to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *bulkComposeRecordEncoder) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	_, ok := state.Res.(*resource.ComposeRecord)
	if !ok {
		return encoderErrInvalidResource(types.RecordResourceType, state.Res.ResourceType())
	}

	return nil
}

// Encode encodes the composeRecord to the document
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *bulkComposeRecordEncoder) Encode(ctx context.Context, w io.Writer, state *envoy.ResourceState) (err error) {
	enc := csv.NewWriter(w)
	all := len(n.encoderConfig.Fields) == 0

	// Generate header & cell index
	hh := make([]string, 0, 100)
	hh = append(hh, "id")

	fx := make(map[string]int)
	// 1 sys fields
	sysFields := []string{"ownedBy", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy"}
	for _, sf := range sysFields {
		if all || n.encoderConfig.Fields[sf] {
			fx[sf] = len(hh)
			hh = append(hh, sf)
		}
	}

	// 2 modlue fields
	for _, f := range n.res.RelMod.Fields {
		if all || n.encoderConfig.Fields[f.Name] {
			fx[f.Name] = len(hh)
			hh = append(hh, f.Name)
		}
	}
	enc.Write(hh)

	err = n.res.Walker(func(r *resource.ComposeRecordRaw) error {

		row := make([]string, len(hh))

		var err error
		row[0] = r.ID

		if r.Us != nil {
			if r.Us.OwnedBy != nil && (all || n.encoderConfig.Fields["ownedBy"]) {
				row[fx["ownedBy"]], err = n.res.UserFlakes.GetByStamp(r.Us.OwnedBy).Model()
			}
			if r.Us.CreatedBy != nil && (all || n.encoderConfig.Fields["createdBy"]) {
				row[fx["createdBy"]], err = n.res.UserFlakes.GetByStamp(r.Us.CreatedBy).Model()
			}
			if r.Us.UpdatedBy != nil && (all || n.encoderConfig.Fields["updatedBy"]) {
				row[fx["updatedBy"]], err = n.res.UserFlakes.GetByStamp(r.Us.UpdatedBy).Model()
			}
			if r.Us.DeletedBy != nil && (all || n.encoderConfig.Fields["deletedBy"]) {
				row[fx["deletedBy"]], err = n.res.UserFlakes.GetByStamp(r.Us.DeletedBy).Model()
			}
		}

		if r.Ts != nil {
			r.Ts, err = r.Ts.Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
			if err != nil {
				return err
			}
			if r.Ts.CreatedAt != nil && (all || n.encoderConfig.Fields["createdAt"]) {
				row[fx["createdAt"]] = r.Ts.CreatedAt.S
			}
			if r.Ts.UpdatedAt != nil && (all || n.encoderConfig.Fields["updatedAt"]) {
				row[fx["updatedAt"]] = r.Ts.UpdatedAt.S
			}
			if r.Ts.DeletedAt != nil && (all || n.encoderConfig.Fields["deletedAt"]) {
				row[fx["deletedAt"]] = r.Ts.DeletedAt.S
			}
		}
		if err != nil {
			return err
		}

		for k, v := range r.Values {
			// Skip fields we don't want
			if !all && !n.encoderConfig.Fields[k] {
				continue
			}

			cell, has := fx[k]
			if !has {
				return fmt.Errorf("unknown cell %s", k)
			}
			f := n.res.RelMod.Fields.FindByName(k)
			if f == nil {
				return fmt.Errorf("field %s not found", k)
			}

			if f.Kind == "User" {
				row[cell], err = n.res.UserFlakes.GetByKey(v).Model()
				if err != nil {
					return err
				}
			} else {
				row[cell] = v
			}
		}
		return enc.Write(row)
	})
	if err != nil {
		return err
	}

	enc.Flush()
	return nil
}
