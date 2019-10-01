package importer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
)

type (
	Record struct {
		imp       *Importer
		namespace *types.Namespace
		set       map[string]types.RecordSet
		dirty     map[uint64]bool
		// modRefs   []recordModuleRef
	}

	// recordModuleRef struct {
	// 	// record handle, report index, module handle
	// 	ch string
	// 	ri int
	// 	mh string
	// }
	//
	// recordFinder interface {
	// 	Find(filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)
	// }
)

func NewRecordImporter(imp *Importer, ns *types.Namespace) *Record {
	out := &Record{
		imp:       imp,
		namespace: ns,
		set:       make(map[string]types.RecordSet),
		dirty:     make(map[uint64]bool),
	}

	return out
}

func (rImp *Record) getModule(handle string) (*types.Module, error) {
	if g, ok := rImp.imp.namespaces.modules[rImp.namespace.Slug]; !ok {
		return nil, errors.Errorf("could not get modules %q from non existing namespace %q", handle, rImp.namespace.Slug)
	} else {
		return g.Get(handle)
	}
}

// CastSet
// { <module-handle>: [ { record }, ... ]
func (rImp *Record) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(_ int, moduleHandle string, def interface{}) error {
		return rImp.Cast(moduleHandle, def)
	})
}

// Cast
// { <module-handle>: [ { record }, ... ]
func (rImp *Record) Cast(moduleHandle string, def interface{}) (err error) {
	if !deinterfacer.IsSlice(def) {
		return errors.New("expecting set of records")
	}

	var (
		record *types.Record
		module *types.Module
	)

	if module, err = rImp.getModule(moduleHandle); err != nil {
		return
	} else if module == nil {
		return errors.Errorf("unexisting module %q", moduleHandle)
	}

	if _, has := rImp.set[module.Handle]; !has {
		rImp.set[module.Handle] = types.RecordSet{}
	}

	return deinterfacer.Each(def, func(i int, _ string, rdef interface{}) (err error) {
		record = &types.Record{
			ModuleID:    module.ID,
			NamespaceID: module.NamespaceID,
		}

		rImp.set[module.Handle] = append(rImp.set[module.Handle], record)

		return deinterfacer.Each(rdef, func(_ int, key string, val interface{}) (err error) {
			switch key {
			// case "handle":
			// @todo add support for record handle (only virtual, for the time of import
			//       so that we can reference one imported record to another

			case "values":
				record.Values, err = rImp.castValues(module, val)

			default:
				return fmt.Errorf("unexpected key %q for record on module %q", key, module.Handle)
			}

			return
		})
	})
}

func (rImp *Record) castValues(module *types.Module, vvDef interface{}) (types.RecordValueSet, error) {
	var rvs = types.RecordValueSet{}

	return rvs, deinterfacer.Each(vvDef, func(_ int, fieldName string, value interface{}) (err error) {
		f := module.Fields.FindByName(fieldName)
		if f == nil {
			return fmt.Errorf("unknown field %q for record on module %q", fieldName, module.Handle)
		}

		if deinterfacer.IsSlice(value) {
			if !f.Multi {
				return fmt.Errorf("field %q on module %q does not support multiple values", fieldName, module.Handle)
			}

			for p, v := range deinterfacer.ToStrings(value) {
				rvs = append(rvs, &types.RecordValue{Name: fieldName, Value: v, Place: uint(p)})
			}
		} else {
			rvs = append(rvs, &types.RecordValue{Name: fieldName, Value: deinterfacer.ToString(value)})

		}

		return
	})
}

func (rImp *Record) Store(ctx context.Context, k recordKeeper) (err error) {
	var module *types.Module

	for mh, rr := range rImp.set {

		if module, err = rImp.getModule(mh); err != nil {
			return
		}

		for _, record := range rr {
			record.NamespaceID = rImp.namespace.ID
			record.ModuleID = module.ID

			if record.ID == 0 {
				record, err = k.Create(record)
			} else if rImp.dirty[record.ID] {
				record, err = k.Update(record)
			}

			if err != nil {
				return
			}

			rImp.dirty[record.ID] = false
		}
	}

	return
}
