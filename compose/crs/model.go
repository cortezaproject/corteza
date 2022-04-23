package crs

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
)

type (
	// cStore is a simplified interface so we can use the store.Storer to assert a valid schema
	cStore interface {
		SearchComposeModules(ctx context.Context, f types.ModuleFilter) (types.ModuleSet, types.ModuleFilter, error)
		SearchComposeModuleFields(ctx context.Context, f types.ModuleFieldFilter) (types.ModuleFieldSet, types.ModuleFieldFilter, error)
		SearchComposeNamespaces(ctx context.Context, f types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error)
	}
)

const (
	// https://www.rfc-editor.org/errata/eid1690
	emailLength = 254

	// Generally the upper most limit
	urlLength = 2048
)

// ReloadModulesFromStore resets state based on the provided cStore
func (crs *composeRecordStore) ReloadModulesFromStore(ctx context.Context, cs cStore) (err error) {
	modules, err := crs.loadModules(ctx, cs)
	if err != nil {
		return
	}

	return crs.ReloadModules(ctx, modules...)
}

// ReloadModulesFromStore resets state based on the provided set of modules
func (crs *composeRecordStore) ReloadModules(ctx context.Context, modules ...*types.Module) (err error) {
	// Clear up the old ones
	// @todo profile if manually removing nested pointers makes it faster
	crs.models = make(map[uint64]data.ModelSet)

	return crs.AddModules(ctx, modules...)
}

// AddModules adds new modules without affecting existing ones
func (crs *composeRecordStore) AddModules(ctx context.Context, modules ...*types.Module) (err error) {
	models, err := crs.modulesToModel(modules...)
	if err != nil {
		return
	}

	for storeID, models := range crs.modelByStore(models) {
		if crs.stores[storeID] == nil {
			return fmt.Errorf("can not add module to store %d: store does not exist", storeID)
		}

		err = crs.addModel(ctx, storeID, models)
		if err != nil {
			return
		}
	}

	return
}

// RemoveModules removes the specified modules
func (crs *composeRecordStore) RemoveModules(ctx context.Context, modules ...*types.Module) (err error) {
	models, err := crs.modulesToModel(modules...)
	if err != nil {
		return
	}

	// validation
	for _, model := range models {
		// Validate existence
		old := crs.getModel(model.StoreID, model.Ident)
		if old == nil {
			return fmt.Errorf("cannot remove module %s: not registered", model.Ident)
		}

		// Validate no leftover references
		// @todo we can probably expand on this quitea bit
		for _, registered := range crs.models {
			refs := registered.FilterByReferenced(model)
			if len(refs) > 0 {
				return fmt.Errorf("cannot remove module %s: referenced by other modules", model.Ident)
			}
		}
	}

	// Work
	for _, model := range models {
		oldModels := crs.models[model.StoreID]
		crs.models[model.StoreID] = make(data.ModelSet, 0, len(oldModels))
		for _, o := range oldModels {
			if o.Ident == model.Ident {
				continue
			}

			crs.models[model.StoreID] = append(crs.models[model.StoreID], o)

		}

		// @todo should the underlying store be notified about this?
	}

	return nil
}

// AlterModule updates the old module with the new one
func (crs *composeRecordStore) AlterModule(ctx context.Context, oldMod, newMod *types.Module) (err error) {
	// validation
	{
		if oldMod.Store.ComposeRecordStoreID != newMod.Store.ComposeRecordStoreID {
			return fmt.Errorf("cannot alter module stored in different record stores: old: %d, new: %d", oldMod.Store.ComposeRecordStoreID, newMod.Store.ComposeRecordStoreID)
		}
	}

	store, oldModel, err := crs.prepModuleDDL(ctx, oldMod)
	if err != nil {
		return
	}
	// store is same so we omit
	_, newModel, err := crs.prepModuleDDL(ctx, newMod)
	if err != nil {
		return
	}

	return store.AlterModel(ctx, oldModel, newModel)
}

// @todo other ddl manipupations...

// ---

func (crs *composeRecordStore) prepModuleDDL(ctx context.Context, module *types.Module) (s Store, model *data.Model, err error) {
	s, _, err = crs.getStore(ctx, module.Store.ComposeRecordStoreID)
	if err != nil {
		return
	}

	models, err := crs.modulesToModel(module)
	if err != nil {
		return
	}
	model = models[0]

	return
}

//  modulesByStore maps the given module set by their CRS
func (crs *composeRecordStore) modulesByStore(modules types.ModuleSet) (out map[uint64]types.ModuleSet) {
	out = make(map[uint64]types.ModuleSet)
	for _, mod := range modules {
		out[mod.Store.ComposeRecordStoreID] = append(out[mod.Store.ComposeRecordStoreID], mod)
	}

	return
}

//  modelByStore maps the given models by their CRS
func (crs *composeRecordStore) modelByStore(models data.ModelSet) (out map[uint64]data.ModelSet) {
	out = make(map[uint64]data.ModelSet)

	for _, model := range models {
		out[model.StoreID] = append(out[model.StoreID], model)
	}

	return
}

// loadModules is a utility to load available modules with all their metadata included
func (crs *composeRecordStore) loadModules(ctx context.Context, cs cStore) (mm types.ModuleSet, err error) {
	var (
		namespaces types.NamespaceSet
		modules    types.ModuleSet
		fields     types.ModuleFieldSet
	)

	namespaces, _, err = cs.SearchComposeNamespaces(ctx, types.NamespaceFilter{})
	if err != nil {
		return
	}

	for _, ns := range namespaces {
		modules, _, err = cs.SearchComposeModules(ctx, types.ModuleFilter{
			NamespaceID: ns.ID,
		})
		if err != nil {
			return
		}

		for _, mod := range modules {
			fields, _, err = cs.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{
				ModuleID: []uint64{mod.ID},
			})
			if err != nil {
				return
			}

			mod.Fields = append(mod.Fields, fields...)
		}
	}

	return
}

func moduleFieldCodec(f *types.ModuleField) (strat data.StoreCodec) {
	// Defaulting to alias
	strat = data.StoreCodecAlias{
		Ident: f.Name,
	}

	switch {
	case f.Encoding.EncodingStrategyAlias != nil:
		strat = data.StoreCodecAlias{
			Ident: f.Encoding.EncodingStrategyAlias.Ident,
		}
	case f.Encoding.EncodingStrategyJSON != nil:
		strat = data.StoreCodecJSON{
			Ident: f.Encoding.EncodingStrategyJSON.Ident,
			Path:  f.Encoding.EncodingStrategyJSON.Path,
		}
	}

	return
}

// ----------

func (crs *composeRecordStore) lookupModel(module *types.Module) (out *data.Model) {
	for _, model := range crs.models[module.Store.ComposeRecordStoreID] {
		if model.ResourceID == module.ID {
			return model
		}
	}

	return nil
}

func (crs *composeRecordStore) modulesToModel(modules ...*types.Module) (out data.ModelSet, err error) {
	refIndex := make(map[uint64]*data.Model)
	out = make(data.ModelSet, 0, len(modules))

	// Initial pass to get everything we can
	for _, module := range modules {
		model := crs.moduleModelInit(module)
		refIndex[module.ID] = model
		out = append(out, model)
	}

	// Add stuff we already have
	for _, models := range crs.models {
		for _, model := range models {
			refIndex[model.ResourceID] = model
		}
	}

	// Build up fields
	for i, mod := range modules {
		out[i].Attributes, err = crs.moduleModelAttributes(mod, refIndex)
		if err != nil {
			return
		}
	}

	return
}

func (crs *composeRecordStore) moduleModelInit(mod *types.Module) (out *data.Model) {
	return &data.Model{
		StoreID:      mod.Store.ComposeRecordStoreID,
		ResourceID:   mod.ID,
		ResourceType: types.ModuleResourceType,

		// @todo this isn't exactly this
		Ident:      mod.Handle,
		Attributes: make(data.AttributeSet, len(mod.Fields)),
	}
}

func (crs *composeRecordStore) moduleModelAttributes(mod *types.Module, refIndex map[uint64]*data.Model) (out data.AttributeSet, err error) {
	for _, f := range mod.Fields {
		attr := &data.Attribute{
			Ident:      f.Name,
			MultiValue: f.Multi,
			Store:      moduleFieldCodec(f),
		}
		out = append(out, attr)

		switch strings.ToLower(f.Kind) {
		case "bool":
			attr.Type = data.TypeBoolean{}
		case "datetime":
			switch {
			case f.IsDateOnly():
				attr.Type = data.TypeDate{}
			case f.IsTimeOnly():
				attr.Type = data.TypeTime{}
			default:
				attr.Type = data.TypeTimestamp{}
			}
		case "email":
			attr.Type = data.TypeText{Length: emailLength}
		case "file":
			attr.Type = data.TypeRef{
				RefModel:     &data.Model{Ident: "attachments"},
				RefAttribute: &data.Attribute{Ident: "id"},
			}
		case "number":
			attr.Type = data.TypeNumber{
				Precision: f.Options.Precision(),
				// Scale: ,
			}
		case "record":
			var refModel *data.Model
			mRefID := f.Options.UInt64("moduleID")
			if mRefID > 0 {
				refModel = refIndex[mRefID]
			}

			attr.Type = data.TypeRef{
				RefModel: refModel,
				RefAttribute: &data.Attribute{
					Ident: "id",
				},
			}
		case "select":
			attr.Type = data.TypeEnum{
				Values: f.SelectOptions(),
			}
		case "string":
			attr.Type = data.TypeText{
				Length: 0,
			}
		case "url":
			attr.Type = data.TypeText{
				Length: urlLength,
			}
		case "user":
			attr.Type = data.TypeRef{
				// @todo...

				RefAttribute: &data.Attribute{
					Ident: "id",
				},
			}
		}
	}

	// System attrs
	out = append(out,
		&data.Attribute{
			Ident: "recordID",
			Store: moduleFieldCodec(&types.ModuleField{
				Name: "recordID",
			}),
			Type: data.TypeID{},
		},
		&data.Attribute{
			Ident: "createdAt",
			Store: moduleFieldCodec(&types.ModuleField{
				Name: "createdAt",
			}),
			Type: data.TypeTimestamp{},
		},
		&data.Attribute{
			Ident: "updatedAt",
			Store: moduleFieldCodec(&types.ModuleField{
				Name: "updatedAt",
			}),
			Type: data.TypeTimestamp{},
		},
		&data.Attribute{
			Ident: "deletedAt",
			Store: moduleFieldCodec(&types.ModuleField{
				Name: "deletedAt",
			}),
			Type: data.TypeTimestamp{},
		},

		&data.Attribute{
			Ident: "ownedBy",
			Store: moduleFieldCodec(&types.ModuleField{
				Name: "ownedBy",
			}),
			Type: data.TypeRef{
				RefAttribute: &data.Attribute{Ident: "id"},
			},
		},
		&data.Attribute{
			Ident: "createdBy",
			Store: moduleFieldCodec(&types.ModuleField{
				Name: "createdBy",
			}),
			Type: data.TypeRef{
				RefAttribute: &data.Attribute{Ident: "id"},
			},
		},
		&data.Attribute{
			Ident: "updatedBy",
			Store: moduleFieldCodec(&types.ModuleField{
				Name: "updatedBy",
			}),
			Type: data.TypeRef{
				RefAttribute: &data.Attribute{Ident: "id"},
			},
		},
		&data.Attribute{
			Ident: "deletedBy",
			Store: moduleFieldCodec(&types.ModuleField{
				Name: "deletedBy",
			}),
			Type: data.TypeRef{
				RefAttribute: &data.Attribute{Ident: "id"},
			},
		},
	)

	return
}

func (crs *composeRecordStore) addModel(ctx context.Context, storeID uint64, models data.ModelSet) (err error) {
	for _, model := range models {
		existing := crs.getModel(storeID, model.Ident)
		if existing != nil {
			return fmt.Errorf("cannot add model %s to store %d: already exists", model.Ident, storeID)
		}

		err = crs.addModelToStore(ctx, storeID, model)
		if err != nil {
			return
		}

		crs.models[storeID] = append(crs.models[storeID], model)
	}

	return
}

func (crs *composeRecordStore) addModelToStore(ctx context.Context, storeID uint64, model *data.Model) (err error) {
	s, _, err := crs.getStore(ctx, storeID)
	if err != nil {
		return err
	}

	available, err := s.Models(ctx)
	if err != nil {
		return err
	}

	// Check if already in there
	if existing := available.FindByIdent(model.Ident); existing != nil {
		// Assert validity
		diff := existing.Diff(model)
		if len(diff) > 0 {
			return fmt.Errorf("model %s exists: model not compatible: %v", existing.Ident, diff)
		}

		return nil
	}

	// Try to add to store
	err = s.AddModel(ctx, model)
	if err != nil {
		return
	}

	return nil
}
