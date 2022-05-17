package dal

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"go.uber.org/zap"
)

type (
	// the core struct that outlines the DAL service facility
	service struct {
		stores map[uint64]StoreConnection

		// Indexed by corresponding storeID
		models map[uint64]ModelSet

		primary StoreConnection

		logger *zap.Logger
		inDev  bool
	}

	crsDefiner interface {
		ComposeRecordStoreID() uint64
		StoreDSN() string
		Capabilities() capabilities.Set
	}

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

	defaultStoreID uint64 = 0

	sysID          = "ID"
	sysNamespaceID = "namespaceID"
	sysModuleID    = "moduleID"
	sysCreatedAt   = "createdAt"
	sysCreatedBy   = "createdBy"
	sysUpdatedAt   = "updatedAt"
	sysUpdatedBy   = "updatedBy"
	sysDeletedAt   = "deletedAt"
	sysDeletedBy   = "deletedBy"
	sysOwnedBy     = "ownedBy"
)

// Service initializes a fresh record store where the given store serves as the default
func Service(ctx context.Context, log *zap.Logger, inDev bool, primary crsDefiner, stores ...crsDefiner) (*service, error) {
	crs := &service{
		stores:  make(map[uint64]StoreConnection),
		models:  make(map[uint64]ModelSet),
		primary: nil,

		logger: log,
		inDev:  inDev,
	}

	var err error

	crs.primary, err = connect(ctx, log, primary, inDev)
	if err != nil {
		return nil, err
	}

	return crs, crs.AddStore(ctx, stores...)
}

// ComposeRecordCreate creates the given records for the given module
func (svc *service) ComposeRecordCreate(ctx context.Context, module *types.Module, records ...ValueGetter) (err error) {
	if !module.Store.Partitioned {
		return fmt.Errorf("only partitioned modules work right now")
	}

	// Determine required capabilities
	requiredCap := capabilities.CreateCapabilities(module.Store.Capabilities...)

	// Determine store
	var s StoreConnection
	if s, _, err = svc.getStore(ctx, module.Store.ComposeRecordStoreID, requiredCap...); err != nil {
		return err
	}

	// Get model
	model := svc.lookupModel(module)
	if model == nil {
		return svc.modelNotFoundErr(module)
	}

	return s.CreateRecords(ctx, model, records...)
}

// @todo...
func (svc *service) ComposeRecordSearch(ctx context.Context, module *types.Module, filter *types.RecordFilter) (records types.RecordSet, outFilter *types.RecordFilter, err error) {
	return
	// // Determine requiredCap we'll need
	// requiredCap := capabilities.SearchCapabilities(module.Store.Capabilities...).Union(svc.recFilterCapabilities(filter))

	// // Connect to datasource
	// var s Store
	// var cc capabilities.Set
	// _ = cc
	// s, cc, err = svc.getStore(ctx, module.Store.ComposeRecordStoreID, requiredCap...)
	// if err != nil {
	// 	return
	// }

	// // Prepare data
	// model := svc.lookupModel(module)
	// if model == nil {
	// 	return nil, nil, svc.modelNotFoundErr(module)
	// }

	// loader, err := s.SearchRecords(ctx, model, nil)
	// if err != nil {
	// 	return
	// }

	// limit := int(filter.Limit)
	// if limit == 0 {
	// 	limit = 10
	// }

	// auxCC := make([]Setter, limit)
	// for i := range auxCC {
	// 	auxCC[i] = &types.Record{}
	// }

	// var ok bool
	// _ = ok
	// for loader.More() && len(records) < int(limit) {
	// 	_, err = loader.Load(model, auxCC)
	// 	if err != nil {
	// 		return
	// 	}

	// 	auxRecords, err := svc.extractRecords(model, auxCC...)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}

	// 	if !capabilities.AccessControlCapabilities().IsSubset(cc...) && filter.Check != nil {
	// 		for _, r := range auxRecords {
	// 			if r == nil {
	// 				continue
	// 			}
	// 			if ok, err = filter.Check(r); err != nil {
	// 				return nil, nil, err
	// 			} else if !ok {
	// 				continue
	// 			}

	// 			records = append(records, r)
	// 		}
	// 	} else {
	// 		for _, r := range auxRecords {
	// 			if r == nil {
	// 				break
	// 			}
	// 			records = append(records, r)
	// 		}
	// 	}
	// }

	// return
}

// ---

// recFilterCapabilities utility helps construct required filter capabilities based on the provided record filter
func (svc *service) recFilterCapabilities(f *types.RecordFilter) (out capabilities.Set) {
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

func (svc service) modelNotFoundErr(module *types.Module) error {
	return fmt.Errorf("cannot create records for module %d: module not registered to svc", module.ID)
}

// AddStore registers the given store definitions as compose record stores
func (svc *service) AddStore(ctx context.Context, definers ...crsDefiner) (err error) {
	for _, definer := range definers {
		svc.stores[definer.ComposeRecordStoreID()], err = connect(ctx, svc.logger, definer, svc.inDev)
		if err != nil {
			return
		}
	}

	return nil
}

// RemoveStore removes the given store definition as a compose record store
func (svc *service) RemoveStore(ctx context.Context, storeID uint64, storeIDs ...uint64) (err error) {
	for _, storeID := range append(storeIDs, storeID) {
		s := svc.stores[storeID]
		if s == nil {
			return fmt.Errorf("can not remove compose record store %d: store does not exist", storeID)
		}

		// Potential cleanups
		if err = s.Close(ctx); err != nil {
			return
		}

		// Remove from registry
		delete(svc.stores, storeID)
	}

	return nil
}

// ---
// Utilities

func (svc *service) getModel(store uint64, ident string) *Model {
	for _, model := range svc.models[store] {
		if model.Ident == ident {
			return model
		}
	}

	return nil
}

// getStore returns a store for the given identifier/capabilities combination
func (svc *service) getStore(ctx context.Context, storeID uint64, cc ...capabilities.Capability) (store StoreConnection, can capabilities.Set, err error) {
	err = func() error {
		// get the requested store
		if storeID == defaultStoreID {
			store = svc.primary
		} else {
			store = svc.stores[storeID]
		}
		if store == nil {
			return fmt.Errorf("could not get store %d: store does not exist", storeID)
		}

		// check if store supports requested capabilities
		if !store.Can(cc...) {
			return fmt.Errorf("store does not support requested capabilities: %v", capabilities.Set(cc).Diff(store.Capabilities()))
		}
		can = store.Capabilities()
		return nil
	}()

	if err != nil {
		err = fmt.Errorf("could not connect to store %d: %v", storeID, err)
		return
	}

	return
}

// ReloadModulesFromStore resets state based on the provided cStore
func (svc *service) ReloadModulesFromStore(ctx context.Context, cs cStore) (err error) {
	modules, err := svc.loadModules(ctx, cs)
	if err != nil {
		return
	}

	return svc.ReloadModules(ctx, modules...)
}

// ReloadModulesFromStore resets state based on the provided set of modules
func (svc *service) ReloadModules(ctx context.Context, modules ...*types.Module) (err error) {
	// Clear up the old ones
	// @todo profile if manually removing nested pointers makes it faster
	svc.models = make(map[uint64]ModelSet)

	return svc.AddModules(ctx, modules...)
}

// AddModules adds new modules without affecting existing ones
func (svc *service) AddModules(ctx context.Context, modules ...*types.Module) (err error) {
	models, err := svc.modulesToModel(modules...)
	if err != nil {
		return
	}

	var (
		s StoreConnection
	)

	for storeID, models := range svc.modelByStore(models) {
		s, _, err = svc.getStore(ctx, storeID)
		if err != nil {
			return err
		}

		err = svc.addModel(ctx, s, storeID, models)
		if err != nil {
			return
		}
	}

	return
}

// RemoveModules removes the specified modules
func (svc *service) RemoveModules(ctx context.Context, modules ...*types.Module) (err error) {
	models, err := svc.modulesToModel(modules...)
	if err != nil {
		return
	}

	// validation
	for _, model := range models {
		// Validate existence
		old := svc.getModel(model.StoreID, model.Ident)
		if old == nil {
			return fmt.Errorf("cannot remove module %s: not registered", model.Ident)
		}

		// Validate no leftover references
		// @todo we can probably expand on this quitea bit
		for _, registered := range svc.models {
			refs := registered.FilterByReferenced(model)
			if len(refs) > 0 {
				return fmt.Errorf("cannot remove module %s: referenced by other modules", model.Ident)
			}
		}
	}

	// Work
	for _, model := range models {
		oldModels := svc.models[model.StoreID]
		svc.models[model.StoreID] = make(ModelSet, 0, len(oldModels))
		for _, o := range oldModels {
			if o.Ident == model.Ident {
				continue
			}

			svc.models[model.StoreID] = append(svc.models[model.StoreID], o)

		}

		// @todo should the underlying store be notified about this?
	}

	return nil
}

// AlterModule updates the old module with the new one
func (svc *service) AlterModule(ctx context.Context, oldMod, newMod *types.Module) (err error) {
	return
	// // validation
	// {
	// 	if oldMod.Store.ComposeRecordStoreID != newMod.Store.ComposeRecordStoreID {
	// 		return fmt.Errorf("cannot alter module stored in different record stores: old: %d, new: %d", oldMod.Store.ComposeRecordStoreID, newMod.Store.ComposeRecordStoreID)
	// 	}
	// }

	// store, oldModel, err := svc.prepModuleDDL(ctx, oldMod)
	// if err != nil {
	// 	return
	// }
	// // store is same so we omit
	// _, newModel, err := svc.prepModuleDDL(ctx, newMod)
	// if err != nil {
	// 	return
	// }

	// return store.AlterModel(ctx, oldModel, newModel)
}

// @todo other ddl manipupations...

// ---

// func (crs *service) prepModuleDDL(ctx context.Context, module *types.Module) (s Store, model *data.Model, err error) {
// 	s, _, err = crs.getStore(ctx, module.Store.ComposeRecordStoreID)
// 	if err != nil {
// 		return
// 	}

// 	models, err := crs.modulesToModel(module)
// 	if err != nil {
// 		return
// 	}
// 	model = models[0]

// 	return
// }

//  modelByStore maps the given models by their CRS
func (svc *service) modelByStore(models ModelSet) (out map[uint64]ModelSet) {
	out = make(map[uint64]ModelSet)

	for _, model := range models {
		out[model.StoreID] = append(out[model.StoreID], model)
	}

	return
}

// loadModules is a utility to load available modules with all their metadata included
func (svc *service) loadModules(ctx context.Context, cs cStore) (mm types.ModuleSet, err error) {
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

// moduleFieldCodec is a little utility to construct the store codec we need
// @todo compose/types.Module to dal.Model conversion MUST happen inside
//       compose/service package
func moduleFieldCodec(f *types.ModuleField) (strat Codec) {
	// Defaulting to alias
	strat = CodecAlias{
		Ident: f.Name,
	}

	switch {
	case f.Encoding.EncodingStrategyAlias != nil:
		strat = CodecAlias{
			Ident: f.Encoding.EncodingStrategyAlias.Ident,
		}
	case f.Encoding.EncodingStrategyJSON != nil:
		strat = CodecRecordValueSetJSON{
			Ident: f.Encoding.EncodingStrategyJSON.Ident,
		}
	}

	return
}

// ----------

// @todo compose/types.Module to dal.Model conversion MUST happen inside
//       compose/service package
func (svc *service) lookupModel(module *types.Module) (out *Model) {
	for _, model := range svc.models[module.Store.ComposeRecordStoreID] {
		if model.ResourceID == module.ID {
			return model
		}
	}

	return nil
}

// @todo compose/types.Module to dal.Model conversion MUST happen inside
//       compose/service package
func (svc *service) modulesToModel(modules ...*types.Module) (out ModelSet, err error) {
	refIndex := make(map[uint64]*Model)
	out = make(ModelSet, 0, len(modules))

	// Initial pass to get everything we can
	for _, module := range modules {
		model := svc.moduleModelInit(module)
		refIndex[module.ID] = model
		out = append(out, model)
	}

	// Add stuff we already have
	for _, models := range svc.models {
		for _, model := range models {
			refIndex[model.ResourceID] = model
		}
	}

	// Build up fields
	for i, mod := range modules {
		out[i].Attributes, err = svc.moduleModelAttributes(mod, refIndex)
		if err != nil {
			return
		}
	}

	return
}

// @todo compose/types.Module to dal.Model conversion MUST happen inside
//       compose/service package
func (svc *service) moduleModelInit(mod *types.Module) (out *Model) {
	return &Model{
		StoreID:      mod.Store.ComposeRecordStoreID,
		ResourceID:   mod.ID,
		ResourceType: types.ModuleResourceType,

		Ident:      formatPartitionIdent(mod),
		Attributes: make(AttributeSet, len(mod.Fields)),
	}
}

// @todo compose/types.Module to dal.Model conversion MUST happen inside
//       compose/service package
func (svc *service) moduleModelAttributes(mod *types.Module, refIndex map[uint64]*Model) (out AttributeSet, err error) {
	for _, f := range mod.Fields {
		attr := &Attribute{
			Ident:      f.Name,
			MultiValue: f.Multi,
			Store:      moduleFieldCodec(f),
		}
		out = append(out, attr)

		switch strings.ToLower(f.Kind) {
		case "bool":
			attr.Type = TypeBoolean{}
		case "datetime":
			switch {
			case f.IsDateOnly():
				attr.Type = TypeDate{}
			case f.IsTimeOnly():
				attr.Type = TypeTime{}
			default:
				attr.Type = TypeTimestamp{}
			}
		case "email":
			attr.Type = TypeText{Length: emailLength}
		case "file":
			attr.Type = TypeRef{
				RefModel:     &Model{Ident: "attachments"},
				RefAttribute: &Attribute{Ident: "id"},
			}
		case "number":
			attr.Type = TypeNumber{
				Precision: f.Options.Precision(),
				// Scale: ,
			}
		case "record":
			var refModel *Model
			mRefID := f.Options.UInt64("moduleID")
			if mRefID > 0 {
				refModel = refIndex[mRefID]
			}

			attr.Type = TypeRef{
				RefModel: refModel,
				RefAttribute: &Attribute{
					Ident: "id",
				},
			}
		case "select":
			attr.Type = TypeEnum{
				Values: f.SelectOptions(),
			}
		case "string":
			attr.Type = TypeText{
				Length: 0,
			}
		case "url":
			attr.Type = TypeText{
				Length: urlLength,
			}
		case "user":
			attr.Type = TypeRef{
				// @todo...

				RefAttribute: &Attribute{
					Ident: "id",
				},
			}
		}
	}

	// System attrs
	out = append(out,
		&Attribute{
			Ident: sysID,
			Store: moduleFieldCodec(&types.ModuleField{
				Name: sysID,
			}),
			Type: TypeID{},
		},
		&Attribute{
			Ident: sysCreatedAt,
			Store: moduleFieldCodec(&types.ModuleField{
				Name: sysCreatedAt,
			}),
			Type: TypeTimestamp{},
		},
		&Attribute{
			Ident: sysUpdatedAt,
			Store: moduleFieldCodec(&types.ModuleField{
				Name: sysUpdatedAt,
			}),
			Type: TypeTimestamp{},
		},
		&Attribute{
			Ident: sysDeletedAt,
			Store: moduleFieldCodec(&types.ModuleField{
				Name: sysDeletedAt,
			}),
			Type: TypeTimestamp{},
		},

		&Attribute{
			Ident: sysOwnedBy,
			Store: moduleFieldCodec(&types.ModuleField{
				Name: sysOwnedBy,
			}),
			Type: TypeRef{
				RefAttribute: &Attribute{Ident: "id"},
			},
		},
		&Attribute{
			Ident: sysCreatedBy,
			Store: moduleFieldCodec(&types.ModuleField{
				Name: sysCreatedBy,
			}),
			Type: TypeRef{
				RefAttribute: &Attribute{Ident: "id"},
			},
		},
		&Attribute{
			Ident: sysUpdatedBy,
			Store: moduleFieldCodec(&types.ModuleField{
				Name: sysUpdatedBy,
			}),
			Type: TypeRef{
				RefAttribute: &Attribute{Ident: "id"},
			},
		},
		&Attribute{
			Ident: sysDeletedBy,
			Store: moduleFieldCodec(&types.ModuleField{
				Name: sysDeletedBy,
			}),
			Type: TypeRef{
				RefAttribute: &Attribute{Ident: "id"},
			},
		},
	)

	return
}

func (svc *service) addModel(ctx context.Context, s StoreConnection, storeID uint64, models ModelSet) (err error) {
	for _, model := range models {
		existing := svc.getModel(storeID, model.Ident)
		if existing != nil {
			return fmt.Errorf("cannot add model %s to store %d: already exists", model.Ident, storeID)
		}

		err = svc.addModelToStore(ctx, s, model)
		if err != nil {
			return
		}

		svc.models[storeID] = append(svc.models[storeID], model)
	}

	return
}

func (svc *service) addModelToStore(ctx context.Context, s StoreConnection, model *Model) (err error) {
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

// ---

// @todo compose/types.Module to dal.Model conversion MUST happen inside
//       compose/service package
func formatPartitionIdent(mod *types.Module) string {
	rpl := strings.NewReplacer(
		"{{module}}", mod.Handle,
	)

	if mod.Store.PartitionFormat == "" {
		return mod.Handle
	}

	return rpl.Replace(mod.Store.PartitionFormat)
}
