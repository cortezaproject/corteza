package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/modern-go/reflect2"

	"github.com/crusttech/human/server/pkg/revisions"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/filter"

	"github.com/crusttech/human/server/compose/service/values"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/pkg/slice"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
)

type (
	moduleServices struct {
		locale           ResourceTranslationsManagerService
		dal              dal.FullService
		schemaAltManager schemaAltManager
	}

	moduleAccessController interface {
		CanManageResourceTranslations(ctx context.Context) bool
		CanSearchModulesOnNamespace(context.Context, *types.Namespace) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreateModuleOnNamespace(context.Context, *types.Namespace) bool
		CanReadModule(context.Context, *types.Module) bool
		CanUpdateModule(context.Context, *types.Module) bool
		CanDeleteModule(context.Context, *types.Module) bool
	}

	ModuleService interface {
		FindByID(ctx context.Context, namespaceID, moduleID uint64) (*types.Module, error)
		FindByName(ctx context.Context, namespaceID uint64, name string) (*types.Module, error)
		FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*types.Module, error)
		FindByAny(ctx context.Context, namespaceID uint64, identifier interface{}) (*types.Module, error)
		Find(ctx context.Context, filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)
		Search(ctx context.Context, filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)
		SearchSensitive(ctx context.Context, filter types.PrivacyModuleFilter) (set []types.PrivacyModule, f types.PrivacyModuleFilter, err error)

		Create(ctx context.Context, module *types.Module) (*types.Module, error)
		Update(ctx context.Context, module *types.Module) (*types.Module, error)
		DeleteByID(ctx context.Context, namespaceID, moduleID uint64) error

		// @note probably temporary just so tests are easier
		ReloadDALModels(ctx context.Context) error
	}

	moduleUpdateHandler func(ctx context.Context, ns *types.Namespace, c *types.Module) (moduleChanges, error)

	moduleChanges uint8

	// Model management on DAL Service
	dalModelManager interface {
		GetConnectionByID(ID uint64) *dal.ConnectionWrap
		Search(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (dal.Iterator, error)

		ReplaceModel(context.Context, []*dal.Alteration, *dal.Model) (newAlts []*dal.Alteration, err error)
		RemoveModel(ctx context.Context, connectionID, ID uint64) error
		SearchModelIssues(ID uint64) []dal.Issue
	}
)

const (
	moduleUnchanged     moduleChanges = 0
	moduleChanged       moduleChanges = 1
	moduleLabelsChanged moduleChanges = 2
	moduleFieldsChanged moduleChanges = 4

	recordTable            = "compose_record"
	recordFieldID          = "ID"
	recordFieldModuleID    = "moduleID"
	recordFieldNamespaceID = "namespaceID"
)

const (
	// https://www.rfc-editor.org/errata/eid1690
	emailLength = 254

	// Generally the upper most limit
	urlLength = 2048

	// Fixed total precision for DECIMAL number-field columns; the field's
	// configured Precision() is the scale (decimal places).
	maxPrecisionLength = 15

	sysID          = "ID"
	sysNamespaceID = "namespaceID"
	sysModuleID    = "moduleID"
	sysRevision    = "revision"
	sysMeta        = "meta"
	sysCreatedAt   = "createdAt"
	sysCreatedBy   = "createdBy"
	sysUpdatedAt   = "updatedAt"
	sysUpdatedBy   = "updatedBy"
	sysDeletedAt   = "deletedAt"
	sysDeletedBy   = "deletedBy"
	sysOwnedBy     = "ownedBy"

	colSysID          = "id"
	colSysNamespaceID = "rel_namespace"
	colSysModuleID    = "rel_module"
	colSysRevision    = "revision"
	colSysMeta        = "meta"
	colSysCreatedAt   = "created_at"
	colSysCreatedBy   = "created_by"
	colSysUpdatedAt   = "updated_at"
	colSysUpdatedBy   = "updated_by"
	colSysDeletedAt   = "deleted_at"
	colSysDeletedBy   = "deleted_by"
	colSysOwnedBy     = "owned_by"
)

var (
	systemFields = slice.ToStringBoolMap([]string{
		"recordID",
		"ownedBy",
		"revision",
		"meta",
		"createdBy",
		"createdAt",
		"updatedBy",
		"updatedAt",
		"deletedBy",
		"deletedAt",
	})
)

func Module(am schemaAltManager) *module {
	return &module{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		services: &moduleServices{
			locale:           DefaultResourceTranslation,
			dal:              dal.Service(),
			schemaAltManager: am,
		},
	}
}

func (svc *module) onLookup(ctx context.Context, namespaceID uint64, ID uint64, aProps *moduleActionProps) (*types.Module, error) {
	return svc.lookup(ctx, namespaceID, func(p *moduleActionProps) (*types.Module, error) {
		if ID == 0 {
			return nil, ModuleErrInvalidID()
		}
		p.module.ID = ID
		return store.LookupComposeModuleByID(ctx, svc.store, ID)
	})
}

func (svc *module) onSearch(ctx context.Context, filter types.ModuleFilter, aProps *moduleActionProps) (types.ModuleSet, types.ModuleFilter, error) {
	return svc.Find(ctx, filter)
}

func (svc *module) onCreate(ctx context.Context, new *types.Module) error {
	_, err := svc.createModule(ctx, new)
	return err
}

func (svc *module) onUpdate(ctx context.Context, s store.Storer, upd *types.Module, res *types.Module, aProps *moduleActionProps, before, after func() error) error {
	if !svc.ac.CanUpdateModule(ctx, res) {
		return ModuleErrNotAllowedToUpdate()
	}
	if err := svc.uniqueCheck(ctx, upd); err != nil {
		return err
	}
	if err := validateModuleDedupRules(ctx, upd); err != nil {
		return ModuleErrDedupConfigurationInvalidMissingConstraint()
	}

	// The generated Update wrapper loads only the module row, not its fields.
	// Load the stored fields so updateModuleFields can diff the incoming set
	// against them (deciding what to create/update/delete) and backfill each
	// field's namespace, module and project IDs.
	//
	// Skipping this has two effects: field add/edit/delete silently never
	// persist, and the returned fields carry zero namespace/module IDs — which
	// collapses their RBAC resource to a wildcard, so CanReadRecordValue /
	// CanUpdateRecordValue report false for every field (even for a super-admin,
	// because wildcard resources are rejected before the bypass check).
	if err := loadModuleFields(ctx, s, res); err != nil {
		return err
	}

	ns, err := loadNamespace(ctx, s, res.NamespaceID)
	if err != nil {
		return err
	}

	old := res.Clone()

	res.Name = upd.Name
	res.Handle = upd.Handle
	res.Meta = upd.Meta
	res.Config = upd.Config
	res.Fields = upd.Fields

	// hasRecords protects field name/kind changes once a module holds data; it
	// stays false here (the DAL model is the source of truth for materialized
	// records), matching the behaviour before the codegen refactor.
	hasRecords := false
	if err := updateModuleFields(ctx, s, res, old, hasRecords); err != nil {
		return err
	}

	// Reflect the field changes into the DAL model so records can use the new
	// schema without a reload.
	if err := DalModelReplace(ctx, s, svc.services.schemaAltManager, svc.services.dal, ns, res); err != nil {
		return err
	}

	svc.procDal(res)
	return nil
}

func (svc *module) onDelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Module, aProps *moduleActionProps) error {
	if !svc.ac.CanDeleteModule(ctx, res) {
		return ModuleErrNotAllowedToDelete()
	}
	res.DeletedAt = now()
	if err := store.UpdateComposeModule(ctx, s, res); err != nil {
		return err
	}
	return DalModelRemove(ctx, svc.services.dal, res)
}

func (svc *module) onUndelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Module, aProps *moduleActionProps) error {
	if !svc.ac.CanDeleteModule(ctx, res) {
		return ModuleErrNotAllowedToUndelete()
	}
	res.DeletedAt = nil
	return store.UpdateComposeModule(ctx, s, res)
}

func (svc *module) onReloadDALModels(ctx context.Context, aProps *moduleActionProps) error {
	return DalModelReload(ctx, svc.store, svc.services.schemaAltManager, svc.services.dal)
}

func (svc module) Find(ctx context.Context, filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error) {
	var (
		ns     *types.Namespace
		aProps = &moduleActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Module) (bool, error) {
		if !svc.ac.CanReadModule(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		ns, err = loadNamespace(ctx, svc.store, filter.NamespaceID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		if !svc.ac.CanSearchModulesOnNamespace(ctx, ns) {
			return ModuleErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Module{}.LabelResourceKind(),
				filter.Labels,
				id.Uints(filter.ModuleID...)...,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchComposeModules(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = loadModuleLabels(ctx, svc.store, set...); err != nil {
			return err
		}

		err = loadModuleFields(ctx, svc.store, set...)
		if err != nil {
			return err
		}

		set.Walk(func(m *types.Module) error {
			svc.proc(ctx, m)
			return nil
		})
		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ModuleActionSearch, err)
}

// FindByName tries to find module by name
func (svc module) FindByName(ctx context.Context, namespaceID uint64, name string) (m *types.Module, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		aProps.module.Name = name
		return store.LookupComposeModuleByNamespaceIDName(ctx, svc.store, namespaceID, name)
	})
}

// FindByHandle tries to find module by handle
func (svc module) FindByHandle(ctx context.Context, namespaceID uint64, h string) (m *types.Module, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		if !handle.IsValid(h) {
			return nil, ModuleErrInvalidHandle()
		}

		aProps.module.Handle = h
		return store.LookupComposeModuleByNamespaceIDHandle(ctx, svc.store, namespaceID, h)
	})
}

// FindByAny tries to find module in a particular namespace by id, handle or name
func (svc module) FindByAny(ctx context.Context, namespaceID uint64, identifier interface{}) (m *types.Module, err error) {
	if ID, ok := identifier.(uint64); ok {
		m, err = svc.FindByID(ctx, namespaceID, ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			m, err = svc.FindByID(ctx, namespaceID, ID)
		} else {
			m, err = svc.FindByHandle(ctx, namespaceID, strIdentifier)
			if err == nil && m.ID == 0 {
				m, err = svc.FindByName(ctx, namespaceID, strIdentifier)
			}
		}
	} else {
		// force invalid ID error
		// we do that to wrap error with lookup action context
		_, err = svc.FindByID(ctx, namespaceID, 0)
	}

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (svc module) proc(ctx context.Context, m *types.Module) {
	svc.procLocale(ctx, m)
	svc.procDal(m)
}

func (svc module) procLocale(ctx context.Context, m *types.Module) {
	if svc.services.locale == nil || svc.services.locale.Locale() == nil {
		return
	}

	tag := locale.GetAcceptLanguageFromContext(ctx)
	m.DecodeTranslations(svc.services.locale.Locale().ResourceTranslations(tag, m.ResourceTranslation()))

	m.Fields.Walk(func(mf *types.ModuleField) error {
		mf.DecodeTranslations(svc.services.locale.Locale().ResourceTranslations(tag, mf.ResourceTranslation()))
		return nil
	})
}

func (svc module) procDal(m *types.Module) {
	if svc.services.dal == nil {
		return
	}

	m.Issues = svc.services.dal.SearchModelIssues(m.ID)
	m.Issues = append(m.Issues, svc.services.dal.SearchResourceIssues("corteza::system:revision", m.RbacResource())...)
	if len(m.Issues) == 0 {
		m.Issues = nil
	}
}

func (svc *module) createModule(ctx context.Context, new *types.Module) (*types.Module, error) {
	var (
		ns     *types.Namespace
		aProps = &moduleActionProps{module: new}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Handle) {
			return ModuleErrInvalidHandle()
		}

		for _, f := range new.Fields {
			if systemFields[f.Name] {
				return ModuleErrFieldNameReserved()
			}
		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if !svc.ac.CanCreateModuleOnNamespace(ctx, ns) {
			return ModuleErrNotAllowedToCreate()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		if new.Fields != nil {
			err = new.Fields.Walk(func(f *types.ModuleField) error {
				f.ID = nextID()
				f.ModuleID = new.ID
				f.NamespaceID = new.NamespaceID
				// A field always belongs to the same project as its module.
				f.ProjectID = new.ProjectID
				f.CreatedAt = *now()
				f.UpdatedAt = nil
				f.DeletedAt = nil

				// Assure validatorID
				for i, v := range f.Expressions.Validators {
					v.ValidatorID = uint64(i) + 1
					f.Expressions.Validators[i] = v
				}

				if !handle.IsValid(f.Name) {
					return ModuleErrInvalidHandle()
				}

				return nil
			})
			if err != nil {
				return
			}
		}

		// Verify dal system field mappings
		_ = handleDalSysFieldEncodingUpdate(new)

		aProps.setChanged(new)

		if err = store.CreateComposeModule(ctx, s, new); err != nil {
			return err
		}

		if err = store.CreateComposeModuleField(ctx, s, new.Fields...); err != nil {
			return err
		}

		tt := new.EncodeTranslations()
		for _, f := range new.Fields {
			tt = append(tt, f.EncodeTranslations()...)
		}

		if err = updateTranslations(ctx, svc.ac, svc.services.locale, tt...); err != nil {
			return
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		if err = DalModelReplace(ctx, s, svc.services.schemaAltManager, svc.services.dal, ns, new); err != nil {
			return err
		}

		svc.procDal(new)
		return nil
	})

	return new, svc.recordAction(ctx, aProps, ModuleActionCreate, err)
}

// SearchSensitive will list all module with at least one private module field
func (svc module) SearchSensitive(ctx context.Context, filter types.PrivacyModuleFilter) (set []types.PrivacyModule, f types.PrivacyModuleFilter, err error) {
	var (
		mm types.ModuleSet

		reqConnes    = make(map[uint64]bool)
		hasReqConnes = len(filter.ConnectionID) > 0
	)

	for _, connectionID := range filter.ConnectionID {
		reqConnes[id.Uint(connectionID)] = true
	}

	err = func() error {
		mm, _, err = svc.Find(ctx, types.ModuleFilter{NamespaceID: filter.NamespaceID})
		if err != nil {
			return err
		}

		for _, m := range mm {
			conn := svc.services.dal.GetConnectionByID(m.Config.DAL.ConnectionID)
			if err != nil {
				return err
			}

			connID := conn.ID
			if hasReqConnes && !reqConnes[connID] {
				continue
			}

			isSensitive := false
			for _, f := range m.Fields {
				isSensitive = isSensitive || f.IsSensitive()
			}

			tag := locale.GetAcceptLanguageFromContext(ctx)
			m.DecodeTranslations(svc.services.locale.Locale().ResourceTranslations(tag, m.ResourceTranslation()))

			if isSensitive && m != nil {
				pm := types.PrivacyModule{
					Module: types.PrivacyModuleMeta{
						ID:     m.ID,
						Name:   m.Name,
						Handle: m.Handle,
						Fields: m.Fields,
					},
					ConnectionID: connID,
				}

				set = append(set, pm)
			}
		}

		return nil
	}()

	return set, filter, err
}

// lookup fn() orchestrates module lookup, namespace preload and check, module reading...
func (svc module) lookup(ctx context.Context, namespaceID uint64, lookup func(*moduleActionProps) (*types.Module, error)) (m *types.Module, err error) {
	var aProps = &moduleActionProps{module: &types.Module{NamespaceID: namespaceID}}

	err = func() error {
		if ns, err := loadNamespace(ctx, svc.store, namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if m, err = lookup(aProps); errors.IsNotFound(err) {
			return ModuleErrNotFound()
		} else if err != nil {
			return err
		}

		aProps.setModule(m)

		if !svc.ac.CanReadModule(ctx, m) {
			return ModuleErrNotAllowedToRead()
		}

		if err = loadModuleLabels(ctx, svc.store, m); err != nil {
			return err
		}

		if err = loadModuleFields(ctx, svc.store, m); err != nil {
			return err
		}

		svc.proc(ctx, m)
		return nil
	}()

	return m, svc.recordAction(ctx, aProps, ModuleActionLookup, err)
}

func (svc module) uniqueCheck(ctx context.Context, m *types.Module) (err error) {
	if m.Handle != "" {
		if e, _ := store.LookupComposeModuleByNamespaceIDHandle(ctx, svc.store, m.NamespaceID, m.Handle); e != nil && e.ID > 0 && e.ID != m.ID {
			return ModuleErrHandleNotUnique()
		}
	}

	if m.Name != "" {
		if e, _ := store.LookupComposeModuleByNamespaceIDName(ctx, svc.store, m.NamespaceID, m.Name); e != nil && e.ID > 0 && e.ID != m.ID {
			return ModuleErrNameNotUnique()
		}
	}

	return nil
}

// updates module fields
// expecting to receive all module fields, as it deletes the rest
// also, sort order of the fields is also important as this fn stores and updates field's place as send
func updateModuleFields(ctx context.Context, s store.Storer, new, old *types.Module, hasRecords bool) (err error) {
	// Go over new to assure field integrity
	for _, f := range new.Fields {
		if f.ModuleID == 0 {
			f.ModuleID = new.ID
		}
		if f.NamespaceID == 0 {
			f.NamespaceID = new.NamespaceID
		}
		// Keep the field's project in sync with its module (set unconditionally:
		// the payload never carries it, and a field cannot outlive its module's
		// project).
		f.ProjectID = new.ProjectID

		if systemFields[f.Name] && !old.Fields.HasName(f.Name) {
			// make sure we're backward compatible, or better:
			// if, by some weird case, someone managed to get invalid field name into
			// the store, we'll turn a blind eye.
			return ModuleErrFieldNameReserved()
		}

		// backward compatible; we didn't check for valid handle.
		// if a field already existed and the handle is invalid we ignore the error.
		if !handle.IsValid(f.Name) && old.Fields.FindByName(f.Name) == nil {
			return ModuleErrInvalidHandle()
		}

		if f.ModuleID != new.ID {
			return fmt.Errorf("module id of field %q does not match the module", f.Name)
		}
	}

	// Delete any missing module fields
	n := now()
	ff := make(types.ModuleFieldSet, 0, len(old.Fields))
	for _, of := range old.Fields {
		nf := new.Fields.FindByID(of.ID)

		if nf == nil {
			of.DeletedAt = n
			ff = append(ff, of)
		} else if nf.DeletedAt != nil {
			of.DeletedAt = n
			ff = append(ff, of)
		}
	}

	if len(ff) > 0 {
		err = store.DeleteComposeModuleField(ctx, s, ff...)
		if err != nil {
			return err
		}
	}

	// Next preproc any default values
	new.Fields, err = moduleFieldDefaultPreparer(ctx, s, new, new.Fields)
	if err != nil {
		return err
	}

	// Assure; create/update remaining fields
	idx := 0
	ff = make(types.ModuleFieldSet, 0, len(old.Fields))
	for _, f := range new.Fields {
		if f.DeletedAt != nil {
			continue
		}

		f.Place = idx
		if of := old.Fields.FindByID(f.ID); of != nil {
			f.CreatedAt = of.CreatedAt

			// We do not have any other code in place that would handle changes of field name and kind, so we need
			// to reset any changes made to the field.
			// @todo remove when we are able to handle field rename & type change
			if hasRecords {
				f.Name = of.Name
				f.Kind = of.Kind
			}

			f.UpdatedAt = now()

			err = store.UpdateComposeModuleField(ctx, s, f)
			if err != nil {
				return err
			}

			if label.Changed(f.Labels, of.Labels) {
				if err = label.Update(ctx, s, f); err != nil {
					return
				}
			}

			ff = append(ff, f)
		} else {
			f.ID = nextID()
			f.CreatedAt = *now()

			if err = store.CreateComposeModuleField(ctx, s, f); err != nil {
				return err
			}
			if err = label.Update(ctx, s, f); err != nil {
				return
			}

			ff = append(ff, f)
		}

		idx++
	}

	sort.Sort(ff)
	new.Fields = ff

	return nil
}

func moduleFieldDefaultPreparer(ctx context.Context, s store.Storer, m *types.Module, newFields types.ModuleFieldSet) (types.ModuleFieldSet, error) {
	var err error

	// prepare an auxiliary module to perform isolated validations on
	auxm := &types.Module{
		Handle:      "aux_module",
		NamespaceID: m.NamespaceID,
		Fields:      types.ModuleFieldSet{nil},
	}

	for _, f := range newFields {
		if f.DefaultValue == nil || len(f.DefaultValue) == 0 {
			continue
		}
		auxm.Fields[0] = f

		vv := f.DefaultValue
		vv.SetUpdatedFlag(true)
		// Module field default values should not have a field name, so let's temporarily add it
		vv.Walk(func(rv *types.RecordValue) error {
			rv.Name = f.Name
			return nil
		})

		if err = RecordValueSanitization(auxm, vv); err != nil {
			return nil, err
		}

		vv = values.Sanitizer().Run(auxm, vv)

		r := &types.Record{
			Values: vv,
		}

		rve := defaultValidator(DefaultRecord).Run(ctx, s, auxm, r)
		if !rve.IsValid() {
			return nil, rve
		}

		vv = values.Formatter().Run(auxm, vv)

		// Module field default values should not have a field name, so let's remove it
		vv.Walk(func(rv *types.RecordValue) error {
			rv.Name = ""
			return nil
		})

		f.DefaultValue = vv
	}
	return newFields, nil
}

func loadModuleFields(ctx context.Context, s store.Storer, mm ...*types.Module) (err error) {
	if len(mm) == 0 {
		return nil
	}

	var (
		ff  types.ModuleFieldSet
		mff = types.ModuleFieldFilter{ModuleID: types.ModuleSet(mm).IDs()}
	)

	if ff, _, err = store.SearchComposeModuleFields(ctx, s, mff); err != nil {
		return
	}

	for _, m := range mm {
		m.Fields = ff.FilterByModule(m.ID)
		m.Fields.Walk(func(f *types.ModuleField) error {
			f.NamespaceID = m.NamespaceID
			return nil
		})

		sort.Sort(m.Fields)
	}

	return
}

// loads record module with fields and namespace
func loadModuleCombo(ctx context.Context, s store.Storer, namespaceID, moduleID uint64) (ns *types.Namespace, m *types.Module, err error) {
	ns, err = loadNamespace(ctx, s, namespaceID)
	if err != nil {
		return
	}

	if m, err = loadModule(ctx, s, moduleID); err != nil {
		return
	}

	if namespaceID != m.NamespaceID {
		return nil, nil, ModuleErrNotFound()
	}

	err = loadModuleFields(ctx, s, m)
	return
}

func loadModuleField(ctx context.Context, s store.Storer, namespaceID, moduleID, fieldID uint64) (res *types.ModuleField, err error) {
	if moduleID == 0 {
		return nil, ModuleErrInvalidID()
	}

	if res, err = store.LookupComposeModuleFieldByID(ctx, s, fieldID); errors.IsNotFound(err) {
		err = ModuleErrNotFound()
	}

	if err == nil && (moduleID != res.ModuleID) {
		// Make sure  module-field belongs to the right module
		return nil, ModuleErrNotFound()
	}

	// add namespace ID on the module-field
	res.NamespaceID = namespaceID

	return
}

// loadLabeledModules loads labels on one or more modules and their fields
func loadModuleLabels(ctx context.Context, s store.Labels, set ...*types.Module) error {
	if len(set) == 0 {
		return nil
	}

	mll := make([]label.LabeledResource, 0, len(set))
	fll := make([]label.LabeledResource, 0, len(set)*10)
	for i := range set {
		mll = append(mll, set[i])

		for j := range set[i].Fields {
			fll = append(fll, set[i].Fields[j])
		}
	}

	if err := label.Load(ctx, s, mll...); err != nil {
		return err
	}

	if err := label.Load(ctx, s, fll...); err != nil {
		return err
	}

	return nil
}

func validateModuleDedupRules(ctx context.Context, m *types.Module) error {
	return m.Config.RecordDeDup.Rules.Validate()
}

// DalModelReload reloads all defined compose modules into the DAL
func DalModelReload(ctx context.Context, s store.Storer, am schemaAltManager, dmm dalModelManager) (err error) {
	// Get all available namespaces
	nn, _, err := store.SearchComposeNamespaces(ctx, s, types.NamespaceFilter{})
	if err != nil {
		return
	}

	// Get all available connections
	mm, _, err := store.SearchComposeModules(ctx, s, types.ModuleFilter{})
	if err != nil {
		return
	}

	err = loadModuleFields(ctx, s, mm...)
	if err != nil {
		return err
	}

	// Reload!
	for _, ns := range nn {
		err = DalModelReplace(ctx, s, am, dmm, ns, modulesForNamespace(ns, mm)...)
		if err != nil {
			return
		}
	}

	return
}

// modulesForNamespace returns all of the modules belonging to that namespace
// @todo implement some indexing at an earlier step for faster processing; will do for now
func modulesForNamespace(ns *types.Namespace, mm types.ModuleSet) (out types.ModuleSet) {
	out = make(types.ModuleSet, 0, len(mm))
	for _, m := range mm {
		if m.NamespaceID == ns.ID {
			out = append(out, m)
		}
	}

	return
}

// Replaces all given connections
func DalModelReplace(ctx context.Context, s store.Storer, am schemaAltManager, dmm dalModelManager, ns *types.Namespace, modules ...*types.Module) (err error) {
	var (
		models      dal.ModelSet
		currentAlts []*dal.Alteration
		newAlts     []*dal.Alteration
	)

	models, err = ModulesToModelSet(dmm, ns, modules...)
	if err != nil {
		return
	}

	for _, m := range models {
		if !reflect2.IsNil(am) {
			// @todo this would need to use s from here, not service
			currentAlts, err = am.ModelAlterations(ctx, m)
			if err != nil {
				return
			}
		}

		newAlts, err = dmm.ReplaceModel(ctx, currentAlts, m)
		if err != nil {
			return
		}

		if !reflect2.IsNil(am) {
			err = am.SetAlterations(ctx, s, m, currentAlts, newAlts...)
			if err != nil {
				return
			}
		}
	}

	return
}

// Removes a connection from DAL service
func DalModelRemove(ctx context.Context, dmm dalModelManager, mm ...*types.Module) (err error) {
	for _, m := range mm {
		if err = dmm.RemoveModel(ctx, m.Config.DAL.ConnectionID, m.ID); err != nil {
			return err
		}
	}

	return
}

// ModulesToModelSet takes a modules for a namespace and converts all of them
// into a model set for the DAL
//
// Ident partition placeholders are replaced here as well alongside
// with the revision models where revisions are enabled
func ModulesToModelSet(dmm dalModelManager, ns *types.Namespace, mm ...*types.Module) (out dal.ModelSet, err error) {
	var (
		conn  *dal.ConnectionWrap
		model *dal.Model

		// partition replace pairs
		modPartition []string

		// namespace partition replacement pairs
		// {{namespace}} is replaced with the namespace handle (slug)
		nsPartition = []string{"{{namespace}}", ns.Slug}

		defConnID uint64
		defConn   = dmm.GetConnectionByID(0)
	)

	if defConn != nil {
		defConnID = defConn.ID
	}

	for connectionID, modules := range modulesByConnection(defConnID, mm...) {
		// Get the connection meta
		conn = dmm.GetConnectionByID(connectionID)

		// Convert all modules to models
		for _, mod := range modules {
			if conn == nil {
				// construct a simplified model w/o attributes, connection
				// this will allow us to manage model's issues within
				// the DAL service
				model = &dal.Model{
					Label:      mod.Handle,
					Resource:   mod.RbacResource(),
					ResourceID: mod.ID,
				}

				out = append(out, model)
				continue
			}

			// convert each module to model
			model, err = ModuleToModel(ns, mod, conn.Config.ModelIdent)
			if err != nil {
				return
			}

			// construct partition replacement pairs from namespace & module handles
			// {{module}} is replaced with module handle
			modPartition = append(nsPartition, "{{module}}", mod.Handle)

			// replace all partition replacement pairs
			model.Ident = strings.NewReplacer(modPartition...).Replace(model.Ident)

			// @todo validate ident with connection's ident validator

			model.Constraints = modelBaseConstraints(model, mod)

			model.ConnectionID = connectionID
			out = append(out, model)

			if mod.Config.RecordRevisions.Enabled {
				rModel := revisions.Model()

				// reuse the connection from the module
				rModel.ConnectionID = connectionID
				rModel.Resource = model.Resource
				rModel.ResourceID = nextID()

				if rModel.Ident = mod.Config.RecordRevisions.Ident; rModel.Ident == "" {
					rModel.Ident = "compose_record_revisions"
				}

				rModel.Ident = strings.NewReplacer(modPartition...).Replace(rModel.Ident)

				// @todo validate ident with connection's ident validator

				out = append(out, rModel)
			}
		}
	}

	return
}

func modelBaseConstraints(model *dal.Model, mod *types.Module) (out map[string][]any) {

	// If we're writting to the default table apply additional constraints
	// @todo there should be more logic here, but for now this is what we had
	//       elsewhere.
	if model.Ident == recordTable {
		out = map[string][]any{
			recordFieldModuleID:    {mod.ID},
			recordFieldNamespaceID: {mod.NamespaceID},
		}
	}

	return
}

// ModuleToModel converts a module with fields to DAL model and attributes
//
// note: this function does not do any partition placeholder replacements
func ModuleToModel(ns *types.Namespace, mod *types.Module, inhIdent string) (model *dal.Model, err error) {
	var (
		attrAux dal.AttributeSet
	)

	model = &dal.Model{
		Label:              mod.Handle,
		Resource:           mod.RbacResource(),
		ResourceID:         mod.ID,
		ResourceType:       types.ModuleResourceType,
		SensitivityLevelID: mod.Config.Privacy.SensitivityLevelID,
	}

	userDefinedFieldIdents := make(map[string]bool)

	if model.Ident = mod.Config.DAL.Ident; model.Ident == "" {
		// try with explicitly set ident on module's DAL config
		// and fallback connection's default if it is empty
		model.Ident = inhIdent
	}

	// Refs for lookups
	var (
		nsSlug = ""
		nsID   = uint64(0)
	)
	if ns != nil {
		nsSlug = ns.Slug
		nsID = ns.ID
	}
	model.Refs = map[string]any{
		"module":      mod.Handle,
		"moduleID":    mod.ID,
		"namespace":   nsSlug,
		"namespaceID": nsID,
	}

	// Convert user-defined fields to attributes
	attrAux, err = moduleFieldsToAttributes(mod)
	if err != nil {
		return
	}

	for _, attr := range attrAux {
		userDefinedFieldIdents[attr.Ident] = true
	}

	model.Attributes = append(model.Attributes, attrAux...)

	// Convert system fields to attribute
	attrAux, err = moduleSystemFieldsToAttributes(mod)
	if err != nil {
		return
	}
	for _, attr := range attrAux {
		ok, _ := userDefinedFieldIdents[attr.Ident]
		if !ok {
			// make sure we're backward compatible:
			// if, by some weird case, someone managed to get a system field name into
			// the store, we'll turn a blind eye. We need to make sure not to include the field twice in this situation.
			model.Attributes = append(model.Attributes, attr)
		}
	}

	return
}

// moduleFieldsToAttributes converts all user-defined module fields to attributes
func moduleFieldsToAttributes(mod *types.Module) (out dal.AttributeSet, err error) {
	out = make(dal.AttributeSet, 0, len(mod.Fields))
	var (
		attr *dal.Attribute
	)

	for _, f := range mod.Fields {
		attr, err = moduleFieldToAttribute(f)
		if err != nil {
			return
		}

		if attr == nil {
			// when instructed to omit the attribute
			// by field's encoding strategy
			continue
		}

		out = append(out, attr)
	}

	return
}

// moduleSystemFieldsToAttributes converts all system-defined module fields to attributes
func moduleSystemFieldsToAttributes(mod *types.Module) (out dal.AttributeSet, err error) {
	var (
		sysEnc = mod.Config.DAL.SystemFieldEncoding

		// generate dal.Codec for each attribute
		// using encoding strategy for that attribute
		// with failsafe on CodecAlias
		mfc = func(defStoreIdent string, es *types.EncodingStrategy) dal.Codec {
			switch {
			case es != nil && es.EncodingStrategyAlias != nil:
				return &dal.CodecAlias{
					Ident: es.EncodingStrategyAlias.Ident,
				}
			case es != nil && es.EncodingStrategyJSON != nil:
				return &dal.CodecRecordValueSetJSON{
					Ident: es.EncodingStrategyJSON.Ident,
				}
			case es != nil:
				// assuming omit!
				return nil
			default:
				return &dal.CodecAlias{
					Ident: defStoreIdent,
				}
			}
		}

		// takes a slice of attributes removes one with nil store codec
		filterSkippedAttribtues = func(in ...*dal.Attribute) (out dal.AttributeSet) {
			for _, attr := range in {
				if attr.Store != nil {
					out = append(out, attr)
				}
			}
			return
		}
	)

	aa := filterSkippedAttribtues(
		dal.PrimaryAttribute(sysID, mfc(colSysID, sysEnc.ID)),
		dal.FullAttribute(sysModuleID, &dal.TypeID{}, mfc(colSysModuleID, sysEnc.ModuleID)),
		dal.FullAttribute(sysDeletedBy, &dal.TypeRef{RefModel: &dal.ModelRef{ResourceType: "corteza::system:user"}, Nullable: true}, mfc(colSysDeletedBy, sysEnc.DeletedBy)),
		dal.FullAttribute(sysNamespaceID, &dal.TypeID{}, mfc(colSysNamespaceID, sysEnc.NamespaceID)),
		dal.FullAttribute(sysRevision, &dal.TypeNumber{HasDefault: true, DefaultValue: 0, Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"}}, mfc(colSysRevision, sysEnc.Revision)),
		dal.FullAttribute(sysMeta, &dal.TypeJSON{}, mfc(colSysMeta, sysEnc.Meta)),
		dal.FullAttribute(sysOwnedBy, &dal.TypeRef{RefModel: &dal.ModelRef{ResourceType: "corteza::system:user"}}, mfc(colSysOwnedBy, sysEnc.OwnedBy)),
		dal.FullAttribute(sysCreatedAt, &dal.TypeTimestamp{}, mfc(colSysCreatedAt, sysEnc.CreatedAt)),
		dal.FullAttribute(sysCreatedBy, &dal.TypeRef{RefModel: &dal.ModelRef{ResourceType: "corteza::system:user"}}, mfc(colSysCreatedBy, sysEnc.CreatedBy)),
		dal.FullAttribute(sysUpdatedAt, &dal.TypeTimestamp{Nullable: true}, mfc(colSysUpdatedAt, sysEnc.UpdatedAt)),
		dal.FullAttribute(sysUpdatedBy, &dal.TypeRef{RefModel: &dal.ModelRef{ResourceType: "corteza::system:user"}, Nullable: true}, mfc(colSysUpdatedBy, sysEnc.UpdatedBy)),
		dal.FullAttribute(sysDeletedAt, &dal.TypeTimestamp{Nullable: true}, mfc(colSysDeletedAt, sysEnc.DeletedAt)),
	)

	for _, a := range aa {
		a.System = true
	}

	return append(out, aa...), nil
}

// moduleFieldToAttribute converts the given module field to a DAL attribute
func moduleFieldToAttribute(f *types.ModuleField) (out *dal.Attribute, err error) {
	var (
		// generate dal.Codec for each attribute
		// using encoding strategy for that attribute
		// with failsafe on JSON RVS.
		codec = func(f *types.ModuleField) dal.Codec {
			var es = f.Config.DAL.EncodingStrategy

			switch {
			case es != nil && es.EncodingStrategyPlain != nil:
				return &dal.CodecPlain{}
			case es != nil && es.EncodingStrategyAlias != nil:
				return &dal.CodecAlias{
					Ident: es.EncodingStrategyAlias.Ident,
				}
			case es != nil && es.EncodingStrategyJSON != nil:
				return &dal.CodecRecordValueSetJSON{
					Ident: es.EncodingStrategyJSON.Ident,
				}
			// Only omit if explicitly told to; for module fields, default to JSON
			case es != nil && es.Omit:
				return nil
			default:
				// defaulting to RecordValueSetJSON with
				// default attribute ident from connection
				return &dal.CodecRecordValueSetJSON{
					// ensure JSON encoded record values always have
					// "values" as col ident as a failsafe
					Ident: "values",
				}
			}
		}(f)
	)

	if codec == nil {
		return
	}

	switch strings.ToLower(f.Kind) {
	case "bool", "boolean":
		at := &dal.TypeBoolean{
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
	case "datetime":
		switch {
		case f.IsDateOnly():
			at := &dal.TypeDate{
				Nullable: !f.Required,
			}
			out = dal.FullAttribute(f.Name, at, codec)
		case f.IsTimeOnly():
			at := &dal.TypeTime{
				Nullable: !f.Required,
			}
			out = dal.FullAttribute(f.Name, at, codec)
		default:
			at := &dal.TypeTimestamp{
				Nullable: !f.Required,
			}
			out = dal.FullAttribute(f.Name, at, codec)
		}
	case "email":
		at := &dal.TypeText{
			Length:   emailLength,
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
	case "file":
		at := &dal.TypeRef{
			RefModel: &dal.ModelRef{Resource: "corteza::system:attachment"},
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
	case "number":
		at := &dal.TypeNumber{
			// @todo precision/scale need a proper rework. Options.Precision() is
			// the number of decimal places, so it is the SCALE; the column keeps
			// a fixed total precision. Emitting it as Precision (with scale 0)
			// produced DECIMAL(<places>,0) and silently truncated every decimal.
			Precision: maxPrecisionLength,
			Scale:     int(f.Options.Precision()),
			Nullable:  !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
	case "record":
		at := &dal.TypeRef{
			RefModel: &dal.ModelRef{
				ResourceID:   f.Options.UInt64("moduleID"),
				ResourceType: types.ModuleResourceType,
			},
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
	case "select":
		at := &dal.TypeEnum{
			Values:   f.SelectOptions(),
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
	case "url":
		at := &dal.TypeText{
			Length:   urlLength,
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
	case "user":
		at := &dal.TypeRef{
			RefModel: &dal.ModelRef{
				ResourceType: systemTypes.UserResourceType,
			},
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
	case "geometry":
		at := &dal.TypeGeometry{
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)
		out.Filterable = false
		out.Sortable = false

	default:
		at := &dal.TypeText{
			Nullable: !f.Required,
		}
		out = dal.FullAttribute(f.Name, at, codec)

	}

	out.SensitivityLevelID = f.Config.Privacy.SensitivityLevelID
	out.Label = f.Label
	out.MultiValue = f.Multi
	return
}

// modulesByConnection groups given modules by the common connectionID
func modulesByConnection(defConnID uint64, modules ...*types.Module) map[uint64]types.ModuleSet {
	var (
		id  uint64
		out = make(map[uint64]types.ModuleSet)
	)
	for _, mod := range modules {
		if id = mod.Config.DAL.ConnectionID; id == 0 {
			// connection not explicitly set on module
			// use default
			id = defConnID
		}

		out[id] = append(out[id], mod)
	}

	return out
}

// handleDalSysFieldEncodingUpdate prevents the mapping from being disabled for certain system field
//
//	IE. `recordID` -> `Module.Config.DAL.SystemFieldEncoding.ID`
func handleDalSysFieldEncodingUpdate(mod *types.Module) error {
	if mod.Config.DAL.SystemFieldEncoding.ID != nil && mod.Config.DAL.SystemFieldEncoding.ID.Omit {
		mod.Config.DAL.SystemFieldEncoding.ID.Omit = false
	}
	return nil
}
