package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	Module struct {
		imp       *Importer
		namespace *types.Namespace
		set       types.ModuleSet
		dirty     map[uint64]bool
	}

	// @todo remove finder strategy, directly provide set of items
	moduleFinder interface {
		Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)
	}
)

func NewModuleImporter(imp *Importer, ns *types.Namespace) *Module {
	out := &Module{
		imp:       imp,
		namespace: ns,
		set:       types.ModuleSet{},
		dirty:     make(map[uint64]bool),
	}

	if imp.moduleFinder != nil && ns.ID > 0 {
		out.set, _, _ = imp.moduleFinder.Find(types.ModuleFilter{NamespaceID: ns.ID})
	}

	return out
}

func (pImp *Module) getPageImporter() (*Page, error) {
	if pi, ok := pImp.imp.namespaces.pages[pImp.namespace.Slug]; !ok {
		return nil, errors.Errorf("non existing namespace %q", pImp.namespace.Slug)
	} else {
		return pi, nil
	}
}

// CastSet Resolves permission rules:
// { <module-handle>: { module } } or [ { module }, ... ]
func (mImp *Module) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Modules defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return mImp.Cast(handle, def)
	})
}

// Cast Resolves permission rules:
// { <module-handle>: { module } } or [ { module }, ... ]
func (mImp *Module) Cast(handle string, def interface{}) (err error) {
	if !deinterfacer.IsMap(def) {
		return errors.New("expecting map of values for module")
	}

	var module *types.Module

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid module handle")
	}

	handle = importer.NormalizeHandle(handle)
	if module, err = mImp.Get(handle); err != nil {
		return err
	} else if module == nil {
		module = &types.Module{
			Handle: handle,
			Name:   handle,
		}

		mImp.set = append(mImp.set, module)
	} else if module.ID == 0 {
		return errors.Errorf("module handle %q already defined in this import session", module.Handle)
	}

	mImp.dirty[module.ID] = true

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "namespace":
			// namespace value sanity check
			if deinterfacer.ToString(val, mImp.namespace.Slug) != mImp.namespace.Slug {
				return fmt.Errorf("explicitly set namespace on module %q shadows inherited namespace", mImp.namespace.Slug)
			}

		case "handle":
			// handle value sanity check
			if deinterfacer.ToString(val, handle) != handle {
				return fmt.Errorf("explicitly set handle on module %q shadows inherited handle", handle)
			}

		case "name", "title", "label":
			module.Name = deinterfacer.ToString(val)

		case "page":
			if pi, err := mImp.getPageImporter(); err != nil {
				return err
			} else {
				// Use module's handle for page
				return pi.Cast(handle, val)
			}

		case "meta":
			module.Meta, err = json.Marshal(deinterfacer.Simplify(val))
			return

		case "fields":
			if err = mImp.castFields(module, val); err != nil {
				return
			}

			// Stable order to prevent tests
			// from breaking up
			sort.Sort(module.Fields)

		// case "records":
		// 	return c.resolveRecords(val)

		case "allow", "deny":
			return mImp.imp.permissions.CastSet(types.ModulePermissionResource.String()+handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for module %q", key, handle)
		}

		return err
	})
}

func (mImp *Module) castFields(module *types.Module, def interface{}) (err error) {
	return deinterfacer.Each(def, func(_ int, fieldName string, val interface{}) (err error) {
		deinterfacer.KVsetString(&fieldName, "name", val)
		field := module.Fields.FindByName(fieldName)

		if field == nil {
			field = &types.ModuleField{
				Name:  fieldName,
				Label: fieldName,
				Place: len(module.Fields),
			}

			module.Fields = append(module.Fields, field)
		}

		return deinterfacer.Each(val, func(_ int, key string, val interface{}) (err error) {
			switch key {
			case "name":
				// already handled

			case "label":
				field.Label = deinterfacer.ToString(val)

			case "kind", "type":
				field.Kind = deinterfacer.ToString(val)

			case "options":
				return mImp.castFieldOptions(field, val)

			case "private":
				field.Private = deinterfacer.ToBool(val)

			case "required":
				field.Required = deinterfacer.ToBool(val)

			case "visible":
				field.Visible = deinterfacer.ToBool(val)

			case "multi":
				field.Multi = deinterfacer.ToBool(val)

			case "default":
				field.DefaultValue = types.RecordValueSet{}
				return deinterfacer.Each(val, func(place int, _ string, val interface{}) (err error) {
					field.DefaultValue = append(field.DefaultValue, &types.RecordValue{
						Value: deinterfacer.ToString(val),
						Place: uint(place),
					})

					return
				})

			case "allow", "deny":
				return mImp.imp.permissions.CastSet(types.ModuleFieldPermissionResource.String()+fieldName, key, val)

			default:
				return fmt.Errorf("unexpected key %q for field %q on module %q", key, fieldName, module.Name)
			}

			return err
		})
	})
}

func (mImp *Module) castFieldOptions(field *types.ModuleField, def interface{}) (err error) {
	field.Options = map[string]interface{}{}
	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		field.Options[key] = deinterfacer.Simplify(val)
		return
	})
}

// Get existing modules
func (mImp *Module) Get(handle string) (*types.Module, error) {
	handle = importer.NormalizeHandle(handle)
	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid module handle")
	}

	return mImp.set.FindByHandle(handle), nil
}

func (mImp *Module) Store(ctx context.Context, k moduleKeeper) (err error) {
	// Save everything
	for _, module := range mImp.set {
		var handle = module.Handle

		if module.ID == 0 {
			module.NamespaceID = mImp.namespace.ID
			module, err = k.Create(module)
		} else if mImp.dirty[module.ID] {
			module, err = k.Update(module)
		}

		if err != nil {
			return
		}

		mImp.dirty[module.ID] = false
		mImp.imp.permissions.UpdateResources(types.ModulePermissionResource.String(), handle, module.ID)

		err = module.Fields.Walk(func(f *types.ModuleField) error {
			mImp.imp.permissions.UpdateResources(types.ModuleFieldPermissionResource.String(), f.Name, f.ID)
			return nil
		})
	}

	// Now, resolve the refs & save resolved options again
	var refs uint

	for _, module := range mImp.set {
		if refs, err = mImp.resolveRefs(module); err != nil {
			return errors.Wrap(err, "could not resolve refs")
		} else if refs >= 0 {
			module.UpdatedAt = nil
			if _, err = k.Update(module); err != nil {
				return errors.Wrap(err, "could not update resolved refs")
			}

		}
	}

	return
}

// Resolve all refs for this module
func (mImp *Module) resolveRefs(module *types.Module) (uint, error) {
	var refs uint

	return refs, func() error {
		for _, field := range module.Fields {
			if field.Options == nil {
				continue
			}

			if h, ok := field.Options["module"]; ok {
				refHandle := deinterfacer.ToString(h)

				if refHandle == "" {
					return errors.Errorf("empty module reference handle on module %q, field %q options",
						module.Handle, field.Name)

				}

				if !handle.IsValid(refHandle) {
					return errors.Errorf("invalid module handle %q used for reference on module %q, field %q options",
						refHandle, module.Handle, field.Name)
				}

				if refmod, err := mImp.Get(refHandle); err != nil || refmod == nil {
					return errors.Errorf("could not load module %q on module %q, field %q options (err: %v)",
						refHandle, module.Handle, field.Name, err)
				} else {
					refs++
					field.Options["moduleID"] = strconv.FormatUint(refmod.ID, 10)
					delete(field.Options, "module")
				}
			}
		}

		return nil
	}()
}
