package crs

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
)

func (crs *composeRecordStore) ComposeRecordCreate(ctx context.Context, module *types.Module, records ...*types.Record) (err error) {
	// Determine required capabilities
	requiredCap := capabilities.CreateCapabilities(module.Store.Capabilities...)

	// Determine store
	var s Store
	if s, _, err = crs.getStore(ctx, module.Store.ComposeRecordStoreID, requiredCap...); err != nil {
		return err
	}

	// Prepare data
	model := crs.lookupModel(module)
	if model == nil {
		return crs.modelNotFoundErr(module)
	}

	aux := make([]Getter, len(records))
	for i, r := range records {
		aux[i] = Getter(r)
	}
	return s.CreateRecords(ctx, model, aux...)
}

func (crs *composeRecordStore) ComposeRecordSearch(ctx context.Context, module *types.Module, filter *types.RecordFilter) (records types.RecordSet, outFilter *types.RecordFilter, err error) {
	// Determine requiredCap we'll need
	requiredCap := capabilities.SearchCapabilities(module.Store.Capabilities...).Union(crs.recFilterCapabilities(filter))

	// Connect to datasource
	var s Store
	var cc capabilities.Set
	_ = cc
	s, cc, err = crs.getStore(ctx, module.Store.ComposeRecordStoreID, requiredCap...)
	if err != nil {
		return
	}

	// Prepare data
	model := crs.lookupModel(module)
	if model == nil {
		return nil, nil, crs.modelNotFoundErr(module)
	}

	loader, err := s.SearchRecords(ctx, model, nil)
	if err != nil {
		return
	}

	limit := int(filter.Limit)
	if limit == 0 {
		limit = 10
	}

	auxCC := make([]Setter, limit)
	for i := range auxCC {
		auxCC[i] = &types.Record{}
	}

	var ok bool
	_ = ok
	for loader.More() && len(records) < int(limit) {
		_, err = loader.Load(model, auxCC)
		if err != nil {
			return
		}

		auxRecords, err := crs.extractRecords(model, auxCC...)
		if err != nil {
			return nil, nil, err
		}

		if !capabilities.AccessControlCapabilities().IsSubset(cc...) && filter.Check != nil {
			for _, r := range auxRecords {
				if r == nil {
					continue
				}
				if ok, err = filter.Check(r); err != nil {
					return nil, nil, err
				} else if !ok {
					continue
				}

				records = append(records, r)
			}
		} else {
			for _, r := range auxRecords {
				if r == nil {
					break
				}
				records = append(records, r)
			}
		}
	}

	return
}

// ---

func (crs *composeRecordStore) recFilterCapabilities(f *types.RecordFilter) (out capabilities.Set) {
	if f == nil {
		return
	}
	if f.PageCursor != nil {
		out = append(out, capabilities.Paging)
	}

	if f.IncPageNavigation {
		out = append(out, capabilities.Paging)
	}

	if f.IncTotal {
		out = append(out, capabilities.Stats)
	}

	if f.Sort != nil {
		out = append(out, capabilities.Sorting)
	}

	return
}

func (crs composeRecordStore) modelNotFoundErr(module *types.Module) error {
	return fmt.Errorf("cannot create records for module %d: module not registered to crs", module.ID)
}

func (crs composeRecordStore) extractRecords(model *data.Model, gg ...Setter) (out types.RecordSet, err error) {
	out = make(types.RecordSet, len(gg))
	for i, g := range gg {
		out[i] = g.(*types.Record)
	}

	return
}
