package service

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/revisions"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"

	"github.com/cortezaproject/corteza-server/compose/dalutils"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	module struct {
		actionlog actionlog.Recorder
		ac        moduleAccessController
		eventbus  eventDispatcher
		store     store.Storer
		locale    ResourceTranslationsManagerService

		dal dalModelManager
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

		ReplaceModel(context.Context, *dal.Model) error
		RemoveModel(ctx context.Context, connectionID, ID uint64) error
		ReplaceModelAttribute(ctx context.Context, model *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) (err error)
		SearchModelIssues(ID uint64) []error
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

func Module() *module {
	return &module{
		ac:        DefaultAccessControl,
		eventbus:  eventbus.Service(),
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		locale:    DefaultResourceTranslation,
		dal:       dal.Service(),
	}
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
				filter.ModuleID...,
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

// FindByID tries to find module by ID
func (svc module) FindByID(ctx context.Context, namespaceID, moduleID uint64) (m *types.Module, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		if moduleID == 0 {
			return nil, ModuleErrInvalidID()
		}

		aProps.module.ID = moduleID
		return store.LookupComposeModuleByID(ctx, svc.store, moduleID)
	})
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
	if svc.locale == nil || svc.locale.Locale() == nil {
		return
	}

	tag := locale.GetAcceptLanguageFromContext(ctx)
	m.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, m.ResourceTranslation()))

	m.Fields.Walk(func(mf *types.ModuleField) error {
		mf.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, mf.ResourceTranslation()))
		return nil
	})
}

func (svc module) procDal(m *types.Module) {
	if svc.dal == nil {
		return
	}

	ii := svc.dal.SearchModelIssues(m.ID)
	if len(ii) == 0 {
		m.Issues = nil
		return
	}

	m.Issues = make([]string, len(ii))
	for i, err := range ii {
		m.Issues[i] = err.Error()
	}
}

func (svc module) Create(ctx context.Context, new *types.Module) (*types.Module, error) {
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

		if err != nil {

		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if !svc.ac.CanCreateModuleOnNamespace(ctx, ns) {
			return ModuleErrNotAllowedToCreate()
		}

		// Calling before-create scripts
		if err = svc.eventbus.WaitFor(ctx, event.ModuleBeforeCreate(new, nil, ns)); err != nil {
			return err
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

		if err = updateTranslations(ctx, svc.ac, svc.locale, tt...); err != nil {
			return
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		if err = DalModelReplace(ctx, svc.dal, ns, new); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(ctx, event.ModuleAfterCreate(new, nil, ns))

		svc.procDal(new)
		return nil
	})

	return new, svc.recordAction(ctx, aProps, ModuleActionCreate, err)
}

func (svc module) Update(ctx context.Context, upd *types.Module) (c *types.Module, err error) {
	return svc.updater(ctx, upd.NamespaceID, upd.ID, ModuleActionUpdate, svc.handleUpdate(ctx, upd))
}

func (svc module) DeleteByID(ctx context.Context, namespaceID, moduleID uint64) error {
	return trim1st(svc.updater(ctx, namespaceID, moduleID, ModuleActionDelete, svc.handleDelete))
}

func (svc module) UndeleteByID(ctx context.Context, namespaceID, moduleID uint64) error {
	return trim1st(svc.updater(ctx, namespaceID, moduleID, ModuleActionUndelete, svc.handleUndelete))
}

// ReloadDALModels reconstructs the DAL's data model based on the store.Storer
//
// Directly using store so we don't spam the action log
func (svc *module) ReloadDALModels(ctx context.Context) (err error) {
	return DalModelReload(ctx, svc.store, svc.dal)
}

// SearchSensitive will list all module with at least one private module field
func (svc module) SearchSensitive(ctx context.Context, filter types.PrivacyModuleFilter) (set []types.PrivacyModule, f types.PrivacyModuleFilter, err error) {
	var (
		mm types.ModuleSet

		reqConnes    = make(map[uint64]bool)
		hasReqConnes = len(filter.ConnectionID) > 0
	)

	for _, connectionID := range filter.ConnectionID {
		reqConnes[connectionID] = true
	}

	err = func() error {
		mm, _, err = svc.Find(ctx, types.ModuleFilter{NamespaceID: filter.NamespaceID})
		if err != nil {
			return err
		}

		for _, m := range mm {
			conn := svc.dal.GetConnectionByID(m.Config.DAL.ConnectionID)
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
			m.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, m.ResourceTranslation()))

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

func (svc module) updater(ctx context.Context, namespaceID, moduleID uint64, action func(...*moduleActionProps) *moduleAction, fn moduleUpdateHandler) (*types.Module, error) {
	var (
		changes moduleChanges

		ns     *types.Namespace
		m, old *types.Module
		aProps = &moduleActionProps{module: &types.Module{ID: moduleID, NamespaceID: namespaceID}}
		err    error

		defConn *dal.ConnectionWrap
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, m, err = loadModuleCombo(ctx, s, namespaceID, moduleID)
		if err != nil {
			return
		}

		if err = loadModuleLabels(ctx, svc.store, m); err != nil {
			return err
		}

		old = m.Clone()
		// so we can get issues
		svc.procDal(old)

		aProps.setNamespace(ns)
		aProps.setChanged(m)

		if m.DeletedAt == nil {
			err = svc.eventbus.WaitFor(ctx, event.ModuleBeforeUpdate(m, old, ns))
		} else {
			err = svc.eventbus.WaitFor(ctx, event.ModuleBeforeDelete(m, old, ns))
		}

		if err != nil {
			return
		}

		if changes, err = fn(ctx, ns, m); err != nil {
			return err
		}

		if changes&moduleChanged > 0 {
			{
				// properly resolve connection ID 0 to the actual ID of the default connection
				if defConn = svc.dal.GetConnectionByID(0); defConn == nil {
					return fmt.Errorf("could not find default DAL connection")
				}

				if old.Config.DAL.ConnectionID == 0 {
					old.Config.DAL.ConnectionID = defConn.ID
				}
				if m.Config.DAL.ConnectionID == 0 {
					m.Config.DAL.ConnectionID = defConn.ID
				}
			}

			if old.Config.DAL.ConnectionID != m.Config.DAL.ConnectionID {
				return fmt.Errorf("unable to switch connection for existing models: run data migration")
			}

			if err = store.UpdateComposeModule(ctx, svc.store, m); err != nil {
				return err
			}
		}

		if changes&moduleFieldsChanged > 0 {
			var (
				hasRecords bool
				set        types.RecordSet
			)

			// @todo rethink how model issues and attempted module update with records should interact.
			// 			 this is a temporary solution but should be re-thinked.
			modelIssues := svc.dal.SearchModelIssues(m.ID)
			if len(modelIssues) == 0 {
				if set, _, err = dalutils.ComposeRecordsList(ctx, svc.dal, m, types.RecordFilter{Paging: filter.Paging{Limit: 1}, Check: func(r *types.Record) (bool, error) { return true, nil }}); err != nil {
					return err
				}
				hasRecords = len(set) > 0
			} else {
				hasRecords = false
			}

			if err = updateModuleFields(ctx, s, m, old, hasRecords); err != nil {
				return err
			}
		}

		// i18n
		tt := m.EncodeTranslations()
		for _, f := range m.Fields {
			tt = append(tt, f.EncodeTranslations()...)
		}

		if err = updateTranslations(ctx, svc.ac, svc.locale, tt...); err != nil {
			return
		}

		if changes&moduleLabelsChanged > 0 {
			if err = label.Update(ctx, s, m); err != nil {
				return
			}
		}

		if m.DeletedAt == nil {
			if err = svc.eventbus.WaitFor(ctx, event.ModuleAfterUpdate(m, old, ns)); err != nil {
				return err
			}
			if err = DalModelReplace(ctx, svc.dal, ns, old, m); err != nil {
				return err
			}
			if err = dalAttributeReplace(ctx, svc.dal, ns, old, m); err != nil {
				return err
			}
		} else {
			if err = svc.eventbus.WaitFor(ctx, event.ModuleAfterDelete(nil, old, ns)); err != nil {
				return
			}
			if err = DalModelRemove(ctx, svc.dal, m); err != nil {
				return err
			}
		}

		svc.procDal(m)
		return err
	})

	return m, svc.recordAction(ctx, aProps, action, err)
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

func (svc module) handleUpdate(ctx context.Context, upd *types.Module) moduleUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace, res *types.Module) (changes moduleChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return moduleUnchanged, ModuleErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return moduleUnchanged, ModuleErrInvalidHandle()
		}

		if err = svc.uniqueCheck(ctx, upd); err != nil {
			return moduleUnchanged, err
		}

		if !svc.ac.CanUpdateModule(ctx, res) {
			return moduleUnchanged, ModuleErrNotAllowedToUpdate()
		}

		// Get max validatorID for later use
		vvID := make(map[uint64]uint64)
		for _, f := range res.Fields {
			for _, v := range f.Expressions.Validators {
				if vvID[f.ID] < v.ValidatorID {
					vvID[f.ID] = v.ValidatorID
				}
			}
		}

		if res.Name != upd.Name {
			changes |= moduleChanged
			res.Name = upd.Name
		}

		if res.Handle != upd.Handle {
			changes |= moduleChanged
			res.Handle = upd.Handle
		}

		{
			oldMeta := res.Meta.String()
			if oldMeta == "{}" {
				oldMeta = ""
			}

			newMeta := upd.Meta.String()
			if newMeta == "{}" {
				newMeta = ""
			}

			if oldMeta != newMeta {
				changes |= moduleChanged
				res.Meta = upd.Meta
			}

		}

		if !reflect.DeepEqual(res.Config, upd.Config) {
			changes |= moduleChanged
			res.Config = upd.Config
		}

		// @todo make field-change detection more optimal
		if !reflect.DeepEqual(res.Fields, upd.Fields) {
			changes |= moduleFieldsChanged
			res.Fields = upd.Fields
		}

		// Assure validatorIDs
		for _, f := range res.Fields {
			for j, v := range f.Expressions.Validators {
				if v.ValidatorID == 0 {
					vvID[f.ID] += 1
					v.ValidatorID = vvID[f.ID]
					f.Expressions.Validators[j] = v

					changes |= moduleFieldsChanged
				}
			}
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= moduleLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if changes&moduleChanged > 0 {
			res.UpdatedAt = now()
		}

		// for now, we assume that
		return
	}
}

func (svc module) handleDelete(ctx context.Context, ns *types.Namespace, m *types.Module) (moduleChanges, error) {
	if !svc.ac.CanDeleteModule(ctx, m) {
		return moduleUnchanged, ModuleErrNotAllowedToDelete()
	}

	if m.DeletedAt != nil {
		// module already deleted
		return moduleUnchanged, nil
	}

	m.DeletedAt = now()
	return moduleChanged, nil
}

func (svc module) handleUndelete(ctx context.Context, ns *types.Namespace, m *types.Module) (moduleChanges, error) {
	if !svc.ac.CanDeleteModule(ctx, m) {
		return moduleUnchanged, ModuleErrNotAllowedToUndelete()
	}

	if m.DeletedAt == nil {
		// module not deleted
		return moduleUnchanged, nil
	}

	m.DeletedAt = nil
	return moduleChanged, nil
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

	m, err = loadModule(ctx, s, namespaceID, moduleID)
	return
}

func loadModule(ctx context.Context, s store.Storer, namespaceID, moduleID uint64) (res *types.Module, err error) {
	if moduleID == 0 {
		return nil, ModuleErrInvalidID()
	}

	if res, err = store.LookupComposeModuleByID(ctx, s, moduleID); errors.IsNotFound(err) {
		err = ModuleErrNotFound()
	}

	if err == nil && namespaceID != res.NamespaceID {
		// Make sure chart belongs to the right namespace
		return nil, ModuleErrNotFound()
	}

	if err == nil {
		err = loadModuleFields(ctx, s, res)
	}

	return
}

func loadModuleField(ctx context.Context, s store.Storer, namespaceID, moduleID, fieldID uint64) (res *types.ModuleField, err error) {
	if moduleID == 0 {
		return nil, ModuleErrInvalidID()
	}

	if res, err = store.LookupComposeModuleFieldByID(ctx, s, fieldID); errors.IsNotFound(err) {
		err = ModuleErrNotFound()
	}

	if err == nil && (namespaceID != res.NamespaceID || moduleID != res.ModuleID) {
		// Make sure chart belongs to the right namespace
		return nil, ModuleErrNotFound()
	}

	return
}

// loadLabeledModules loads labels on one or more modules and their fields
//
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

// DalModelReload reloads all defined compose modules into the DAL
func DalModelReload(ctx context.Context, s store.Storer, dmm dalModelManager) (err error) {
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
		err = DalModelReplace(ctx, dmm, ns, modulesForNamespace(ns, mm)...)
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
func DalModelReplace(ctx context.Context, dmm dalModelManager, ns *types.Namespace, modules ...*types.Module) (err error) {
	var (
		models dal.ModelSet
	)

	models, err = modulesToModelSet(dmm, ns, modules...)
	if err != nil {
		return
	}

	for _, m := range models {
		err = dmm.ReplaceModel(ctx, m)
		if err != nil {
			return
		}
	}

	return
}

func dalAttributeReplace(ctx context.Context, dmm dalModelManager, ns *types.Namespace, old, new *types.Module) (err error) {
	oldModel, err := modulesToModelSet(dmm, ns, old)
	if err != nil {
		return
	}
	newModel, err := modulesToModelSet(dmm, ns, new)
	if err != nil {
		return
	}

	diff := oldModel[0].Diff(newModel[0])
	for _, d := range diff {
		if err = dmm.ReplaceModelAttribute(ctx, oldModel[0], d.Original, d.Asserted); err != nil {
			return
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

// modulesToModelSet takes a modules for a namespace and converts all of them
// into a model set for the DAL
//
// Ident partition placeholders are replaced here as well alongside
// with the revision models where revisions are enabled
func modulesToModelSet(dmm dalModelManager, ns *types.Namespace, mm ...*types.Module) (out dal.ModelSet, err error) {
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
		Operations:         mod.Config.DAL.Operations,
	}

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
	model.Attributes = append(model.Attributes, attrAux...)

	// Convert system fields to attribute
	attrAux, err = moduleSystemFieldsToAttributes(mod)
	if err != nil {
		return
	}
	model.Attributes = append(model.Attributes, attrAux...)

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
			default:
				return &dal.CodecAlias{
					Ident: defStoreIdent,
				}
			}
		}
	)

	return append(out,
		dal.PrimaryAttribute(sysID, mfc(colSysID, sysEnc.ID)),
		dal.FullAttribute(sysModuleID, &dal.TypeID{}, mfc(colSysModuleID, sysEnc.ModuleID)),
		dal.FullAttribute(sysDeletedBy, &dal.TypeRef{RefModel: &dal.ModelRef{ResourceType: "corteza::system:user"}, Nullable: true}, mfc(colSysDeletedBy, sysEnc.DeletedBy)),
		dal.FullAttribute(sysNamespaceID, &dal.TypeID{}, mfc(colSysNamespaceID, sysEnc.NamespaceID)),
		dal.FullAttribute(sysRevision, &dal.TypeID{}, mfc(colSysRevision, sysEnc.Revision)),
		dal.FullAttribute(sysMeta, &dal.TypeJSON{}, mfc(colSysMeta, sysEnc.Meta)),
		dal.FullAttribute(sysOwnedBy, &dal.TypeRef{RefModel: &dal.ModelRef{ResourceType: "corteza::system:user"}}, mfc(colSysOwnedBy, sysEnc.OwnedBy)),
		dal.FullAttribute(sysCreatedAt, &dal.TypeTimestamp{}, mfc(colSysCreatedAt, sysEnc.CreatedAt)),
		dal.FullAttribute(sysCreatedBy, &dal.TypeRef{RefModel: &dal.ModelRef{ResourceType: "corteza::system:user"}}, mfc(colSysCreatedBy, sysEnc.CreatedBy)),
		dal.FullAttribute(sysUpdatedAt, &dal.TypeTimestamp{Nullable: true}, mfc(colSysUpdatedAt, sysEnc.UpdatedAt)),
		dal.FullAttribute(sysUpdatedBy, &dal.TypeRef{RefModel: &dal.ModelRef{ResourceType: "corteza::system:user"}, Nullable: true}, mfc(colSysUpdatedBy, sysEnc.UpdatedBy)),
		dal.FullAttribute(sysDeletedAt, &dal.TypeTimestamp{Nullable: true}, mfc(colSysDeletedAt, sysEnc.DeletedAt)),
	), nil
}

// moduleFieldToAttribute converts the given module field to a DAL attribute
func moduleFieldToAttribute(f *types.ModuleField) (out *dal.Attribute, err error) {
	var (
		// generate dal.Codec for each attribute
		// using encoding strategy for that attribute
		// with failsafe on JSON RVS.
		mfc = func(f *types.ModuleField) dal.Codec {
			var es = f.Config.DAL.EncodingStrategy

			switch {
			case es.EncodingStrategyAlias != nil:
				return &dal.CodecAlias{
					Ident: es.EncodingStrategyAlias.Ident,
				}
			case es.EncodingStrategyJSON != nil:
				return &dal.CodecRecordValueSetJSON{
					Ident: es.EncodingStrategyJSON.Ident,
				}
			default:
				// defaulting to RecordValueSetJSON with
				// default attribute ident from connection
				return &dal.CodecRecordValueSetJSON{
					// ensure JSON encoded record values always have
					// "values" as col ident as a failsafe
					Ident: "values",
				}
			}
		}
	)

	switch strings.ToLower(f.Kind) {
	case "bool", "boolean":
		at := &dal.TypeBoolean{}
		out = dal.FullAttribute(f.Name, at, mfc(f))
	case "datetime":
		switch {
		case f.IsDateOnly():
			at := &dal.TypeDate{}
			out = dal.FullAttribute(f.Name, at, mfc(f))
		case f.IsTimeOnly():
			at := &dal.TypeTime{}
			out = dal.FullAttribute(f.Name, at, mfc(f))
		default:
			at := &dal.TypeTimestamp{}
			out = dal.FullAttribute(f.Name, at, mfc(f))
		}
	case "email":
		at := &dal.TypeText{Length: emailLength}
		out = dal.FullAttribute(f.Name, at, mfc(f))
	case "file":
		at := &dal.TypeRef{
			RefModel: &dal.ModelRef{Resource: "corteza::system:attachment"},
		}
		out = dal.FullAttribute(f.Name, at, mfc(f))
	case "number":
		at := &dal.TypeNumber{
			Precision: int(f.Options.Precision()),
		}
		out = dal.FullAttribute(f.Name, at, mfc(f))
	case "record":
		at := &dal.TypeRef{
			RefModel: &dal.ModelRef{
				ResourceID:   f.Options.UInt64("moduleID"),
				ResourceType: types.ModuleResourceType,
			},
		}
		out = dal.FullAttribute(f.Name, at, mfc(f))
	case "select":
		at := &dal.TypeEnum{
			Values: f.SelectOptions(),
		}
		out = dal.FullAttribute(f.Name, at, mfc(f))
	case "url":
		at := &dal.TypeText{
			Length: urlLength,
		}
		out = dal.FullAttribute(f.Name, at, mfc(f))
	case "user":
		at := &dal.TypeRef{
			RefModel: &dal.ModelRef{
				ResourceType: systemTypes.UserResourceType,
			},
		}
		out = dal.FullAttribute(f.Name, at, mfc(f))

	default:
		at := &dal.TypeText{}
		out = dal.FullAttribute(f.Name, at, mfc(f))

	}

	out.SensitivityLevelID = f.Config.Privacy.SensitivityLevelID
	out.Label = f.Name
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
