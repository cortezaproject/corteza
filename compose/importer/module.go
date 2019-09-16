package importer

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	ModuleImport struct {
		namespace *types.Namespace
		set       types.ModuleSet

		pages       *PageImport
		permissions importer.PermissionImporter

		finder moduleFinder
	}

	moduleFinder interface {
		FindByHandle(uint64, string) (*types.Module, error)
	}
)

func NewModuleImporter(ns *types.Namespace, f moduleFinder, pi *PageImport, p importer.PermissionImporter) *ModuleImport {
	return &ModuleImport{
		namespace:   ns,
		set:         types.ModuleSet{},
		pages:       pi,
		permissions: p,
		finder:      f,
	}
}

// CastSet Resolves permission rules:
// { <module-handle>: { module } } or [ { module }, ... ]
func (imp *ModuleImport) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Modules defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return imp.Cast(handle, def)
	})
}

// Cast Resolves permission rules:
// { <module-handle>: { module } } or [ { module }, ... ]
func (imp *ModuleImport) Cast(handle string, def interface{}) (err error) {
	if !deinterfacer.IsMap(def) {
		return errors.New("expecting map of values for module")
	}

	var module *types.Module

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid module handle")
	}

	handle = importer.NormalizeHandle(handle)
	if module, err = imp.Get(handle); err != nil {
		return err
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "namespace":
			// namespace value sanity check
			if deinterfacer.ToString(val, imp.namespace.Slug) != imp.namespace.Slug {
				return fmt.Errorf("explicitly set namespace on module %q shadows inherited namespace", imp.namespace.Slug)
			}

		case "handle":
			// handle value sanity check
			if deinterfacer.ToString(val, handle) != handle {
				return fmt.Errorf("explicitly set handle on module %q shadows inherited handle", handle)
			}

		case "name":
			module.Name = deinterfacer.ToString(val)

		case "page":
			// Use module's handle for page
			return imp.pages.Cast(handle, val)

		case "meta":
			// @todo Module.Meta

		case "fields":
			if err = imp.castFields(module, val); err != nil {
				return err
			}

			// Stable order to prevent tests
			// from breaking up
			sort.Sort(module.Fields)

		// case "records":
		// 	return c.resolveRecords(val)

		case "allow", "deny":
			return imp.permissions.CastSet(types.ModulePermissionResource.String()+handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for module %q", key, handle)
		}

		return err
	})
}

func (imp *ModuleImport) castFields(module *types.Module, def interface{}) (err error) {
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
				return imp.permissions.CastSet(types.ModuleFieldPermissionResource.String()+fieldName, key, val)

			default:
				return fmt.Errorf("unexpected key %q for field %q on module %q", key, fieldName, module.Name)
			}

			return err
		})
	})
}

func (imp *ModuleImport) Exists(handle string) bool {
	handle = importer.NormalizeHandle(handle)

	var (
		module *types.Module
		err    error
	)

	module = imp.set.FindByHandle(handle)
	if module != nil {
		return true
	}

	if imp.namespace.ID == 0 {
		// Assuming new namespace, nothing exists yet..
		return false
	}

	if imp.finder != nil {
		module, err = imp.finder.FindByHandle(imp.namespace.ID, handle)
		if err == nil && module != nil {
			imp.set = append(imp.set, module)
			return true
		}
	}

	return false
}

// finds or makes new module
func (imp *ModuleImport) Get(handle string) (*types.Module, error) {
	handle = importer.NormalizeHandle(handle)

	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid module handle")
	}

	if !imp.Exists(handle) {
		imp.set = append(imp.set, &types.Module{
			Handle: handle,
			Name:   handle,
		})
	}

	return imp.set.FindByHandle(handle), nil
}

func (imp *ModuleImport) Store(ctx context.Context, k moduleKeeper) error {
	return imp.set.Walk(func(module *types.Module) (err error) {
		var handle = module.Handle

		if module.ID == 0 {
			module.NamespaceID = imp.namespace.ID
			module, err = k.Create(module)
		} else {
			module, err = k.Update(module)
		}

		if err != nil {
			return
		}

		imp.permissions.UpdateResources(types.ModulePermissionResource.String(), handle, module.ID)

		err = module.Fields.Walk(func(f *types.ModuleField) error {
			imp.permissions.UpdateResources(types.ModuleFieldPermissionResource.String(), f.Name, f.ID)
			return nil
		})

		return
	})
}
