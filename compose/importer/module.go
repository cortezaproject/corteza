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
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	Module struct {
		imp       *Importer
		namespace *types.Namespace
		set       types.ModuleSet
	}

	moduleFinder interface {
		FindByHandle(uint64, string) (*types.Module, error)
	}
)

func NewModuleImporter(imp *Importer, ns *types.Namespace) *Module {
	return &Module{
		imp:       imp,
		namespace: ns,
		set:       types.ModuleSet{},
	}
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
	if module, err = mImp.GetOrMake(handle); err != nil {
		return err
	}

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
		if fieldKind, ok := val.(string); ok && fieldName != "" {
			// Not much more to do here
			field := module.Fields.FindByName(fieldName)
			if field == nil {
				module.Fields = append(module.Fields, &types.ModuleField{
					Kind:  fieldKind,
					Name:  fieldName,
					Label: fieldName,
				})
			} else {
				field.Kind = fieldKind
			}

			return
		}

		if !deinterfacer.IsMap(def) {
			return errors.New("expecting map of values for module field")
		}

		deinterfacer.KVsetString(&fieldName, "name", def)
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
			case "label":
				field.Label = deinterfacer.ToString(val)
			case "kind", "type":
				field.Kind = deinterfacer.ToString(val)

			case "options":
				// @todo ModuleField.Options
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
	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		field.Options[key] = deinterfacer.Simplify(val)
		return
	})
}

func (mImp *Module) Exists(handle string) bool {
	handle = importer.NormalizeHandle(handle)

	var (
		module *types.Module
		err    error
	)

	module = mImp.set.FindByHandle(handle)
	if module != nil {
		return true
	}

	if mImp.namespace.ID == 0 {
		// Assuming new namespace, nothing exists yet..
		return false
	}

	if mImp.imp.moduleFinder != nil {
		module, err = mImp.imp.moduleFinder.FindByHandle(mImp.namespace.ID, handle)
		if err == nil && module != nil {
			mImp.set = append(mImp.set, module)
			return true
		}
	}

	return false
}

// Get finds or makes a new module
func (mImp *Module) GetOrMake(handle string) (module *types.Module, err error) {
	if module, err = mImp.Get(handle); err != nil {
		return nil, err
	} else if module == nil {
		module = &types.Module{
			Handle: handle,
			Name:   handle,
		}

		mImp.set = append(mImp.set, module)
	}

	return module, nil
}

// Get existing modules
func (mImp *Module) Get(handle string) (*types.Module, error) {
	handle = importer.NormalizeHandle(handle)
	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid module handle")
	}

	if mImp.Exists(handle) {
		return mImp.set.FindByHandle(handle), nil
	} else {
		return nil, nil
	}
}

func (mImp *Module) Store(ctx context.Context, k moduleKeeper) error {
	return mImp.set.Walk(func(module *types.Module) (err error) {
		var handle = module.Handle

		if err = mImp.resolveRefs(module); err != nil {
			return
		}

		if module.ID == 0 {
			module.NamespaceID = mImp.namespace.ID
			module, err = k.Create(module)
		} else {
			module, err = k.Update(module)
		}

		if err != nil {
			return
		}

		mImp.imp.permissions.UpdateResources(types.ModulePermissionResource.String(), handle, module.ID)

		err = module.Fields.Walk(func(f *types.ModuleField) error {
			mImp.imp.permissions.UpdateResources(types.ModuleFieldPermissionResource.String(), f.Name, f.ID)
			return nil
		})

		return
	})
}

// Resolve all refs for this page (page module, inside block)
func (mImp *Module) resolveRefs(module *types.Module) error {
	for i, field := range module.Fields {
		if field.Options == nil {
			continue
		}

		if h, ok := field.Options["module"]; ok {
			if refmod, err := mImp.Get(deinterfacer.ToString(h)); err != nil || refmod == nil {
				return errors.Wrapf(err, "could not load module %q for page %q block #%d",
					h, module.Handle, i+1)
			} else {
				field.Options["moduleID"] = strconv.FormatUint(refmod.ID, 10)
				delete(field.Options, "module")
			}
		}
	}

	return nil
}
