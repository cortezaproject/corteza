package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"sort"
	"strconv"
)

type (
	module struct {
		ctx       context.Context
		actionlog actionlog.Recorder
		ac        moduleAccessController
		eventbus  eventDispatcher
		store     store.Storer
	}

	moduleAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreateModule(context.Context, *types.Namespace) bool
		CanReadModule(context.Context, *types.Module) bool
		CanUpdateModule(context.Context, *types.Module) bool
		CanDeleteModule(context.Context, *types.Module) bool
	}

	ModuleService interface {
		With(ctx context.Context) ModuleService

		FindByID(namespaceID, moduleID uint64) (*types.Module, error)
		FindByName(namespaceID uint64, name string) (*types.Module, error)
		FindByHandle(namespaceID uint64, handle string) (*types.Module, error)
		FindByAny(namespaceID uint64, identifier interface{}) (*types.Module, error)
		Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)

		Create(module *types.Module) (*types.Module, error)
		Update(module *types.Module) (*types.Module, error)
		DeleteByID(namespaceID, moduleID uint64) error
	}

	moduleUpdateHandler func(ctx context.Context, ns *types.Namespace, c *types.Module) (bool, bool, error)
)

func Module() ModuleService {
	return (&module{
		ctx:      context.Background(),
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
	}).With(context.Background())
}

func (svc module) With(ctx context.Context) ModuleService {
	return &module{
		ctx:       ctx,
		actionlog: DefaultActionlog,
		ac:        svc.ac,
		eventbus:  svc.eventbus,
		store:     DefaultStore,
	}
}

func (svc module) Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error) {
	var (
		aProps = &moduleActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Module) (bool, error) {
		if !svc.ac.CanReadModule(svc.ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if ns, err := loadNamespace(svc.ctx, svc.store, filter.NamespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if set, f, err = store.SearchComposeModules(svc.ctx, svc.store, filter); err != nil {
			return err
		}

		return loadModuleFields(svc.ctx, svc.store, set...)
	}()

	return set, f, svc.recordAction(svc.ctx, aProps, ModuleActionSearch, err)
}

// FindByID tries to find module by ID
func (svc module) FindByID(namespaceID, moduleID uint64) (m *types.Module, err error) {
	return svc.lookup(namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		if moduleID == 0 {
			return nil, ModuleErrInvalidID()
		}

		aProps.module.ID = moduleID
		return store.LookupComposeModuleByID(svc.ctx, svc.store, moduleID)
	})
}

// FindByName tries to find module by name
func (svc module) FindByName(namespaceID uint64, name string) (m *types.Module, err error) {
	return svc.lookup(namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		aProps.module.Name = name
		return store.LookupComposeModuleByNamespaceIDName(svc.ctx, svc.store, namespaceID, name)
	})
}

// FindByHandle tries to find module by handle
func (svc module) FindByHandle(namespaceID uint64, h string) (m *types.Module, err error) {
	return svc.lookup(namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		if !handle.IsValid(h) {
			return nil, ModuleErrInvalidHandle()
		}

		aProps.module.Handle = h
		return store.LookupComposeModuleByNamespaceIDHandle(svc.ctx, svc.store, namespaceID, h)
	})
}

// FindByAny tries to find module in a particular namespace by id, handle or name
func (svc module) FindByAny(namespaceID uint64, identifier interface{}) (m *types.Module, err error) {
	if ID, ok := identifier.(uint64); ok {
		m, err = svc.FindByID(namespaceID, ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			m, err = svc.FindByID(namespaceID, ID)
		} else {
			m, err = svc.FindByHandle(namespaceID, strIdentifier)
			if err == nil && m.ID == 0 {
				m, err = svc.FindByName(namespaceID, strIdentifier)
			}
		}
	} else {
		// force invalid ID error
		// we do that to wrap error with lookup action context
		_, err = svc.FindByID(namespaceID, 0)
	}

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (svc module) Create(new *types.Module) (*types.Module, error) {
	var (
		ns     *types.Namespace
		aProps = &moduleActionProps{changed: new}
	)

	err := store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Handle) {
			return ModuleErrInvalidHandle()
		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		if !svc.ac.CanCreateModule(ctx, ns) {
			return ModuleErrNotAllowedToCreate()
		}

		aProps.setNamespace(ns)

		// Calling before-create scripts
		if err = svc.eventbus.WaitFor(ctx, event.ModuleBeforeCreate(new, nil, ns)); err != nil {
			return err
		}

		if err = svc.uniqueCheck(new); err != nil {
			return err
		}

		new.ID = id.Next()
		new.CreatedAt = *nowPtr()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		if new.Fields != nil {
			_ = new.Fields.Walk(func(f *types.ModuleField) error {
				f.ModuleID = new.ID
				f.CreatedAt = *nowPtr()
				f.UpdatedAt = nil
				f.DeletedAt = nil
				return nil
			})
		}

		aProps.setModule(new)

		if err = store.CreateComposeModule(ctx, s, new); err != nil {
			return err
		}

		if err = store.CreateComposeModuleField(ctx, s, new.Fields...); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(ctx, event.ModuleAfterCreate(new, nil, ns))
		return nil
	})

	return new, svc.recordAction(svc.ctx, aProps, ModuleActionCreate, err)
}

func (svc module) Update(upd *types.Module) (c *types.Module, err error) {
	return svc.updater(upd.NamespaceID, upd.ID, ModuleActionUpdate, svc.handleUpdate(upd))
}

func (svc module) DeleteByID(namespaceID, moduleID uint64) error {
	return trim1st(svc.updater(namespaceID, moduleID, ModuleActionDelete, svc.handleDelete))
}

func (svc module) UndeleteByID(namespaceID, moduleID uint64) error {
	return trim1st(svc.updater(namespaceID, moduleID, ModuleActionUndelete, svc.handleUndelete))
}

func (svc module) updater(namespaceID, moduleID uint64, action func(...*moduleActionProps) *moduleAction, fn moduleUpdateHandler) (*types.Module, error) {
	var (
		moduleChanged, fieldsChanged bool

		ns     *types.Namespace
		m, old *types.Module
		aProps = &moduleActionProps{module: &types.Module{ID: moduleID, NamespaceID: namespaceID}}
		err    error
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, m, err = loadModuleWithNamespace(svc.ctx, s, namespaceID, moduleID)
		if err != nil {
			return
		}

		old = m.Clone()

		aProps.setNamespace(ns)
		aProps.setChanged(m)

		if m.DeletedAt == nil {
			err = svc.eventbus.WaitFor(svc.ctx, event.ModuleBeforeUpdate(m, old, ns))
		} else {
			err = svc.eventbus.WaitFor(svc.ctx, event.ModuleBeforeDelete(m, old, ns))
		}

		if err != nil {
			return
		}

		if moduleChanged, fieldsChanged, err = fn(svc.ctx, ns, m); err != nil {
			return err
		}

		if moduleChanged {
			if err = store.UpdateComposeModule(svc.ctx, svc.store, m); err != nil {
				return err
			}
		}

		if fieldsChanged {
			var (
				hasRecords bool
				set        types.RecordSet
			)

			if set, _, err = store.SearchComposeRecords(ctx, s, m, types.RecordFilter{Paging: filter.Paging{Limit: 1}}); err != nil {
				return err
			}

			hasRecords = len(set) > 0

			if err = updateModuleFields(ctx, s, m, m.Fields, hasRecords); err != nil {

			}
		}

		if m.DeletedAt == nil {
			err = svc.eventbus.WaitFor(svc.ctx, event.ModuleAfterUpdate(m, old, ns))
		} else {
			err = svc.eventbus.WaitFor(svc.ctx, event.ModuleAfterDelete(nil, old, ns))
		}

		return err
	})

	return m, svc.recordAction(svc.ctx, aProps, action, err)
}

// lookup fn() orchestrates module lookup, namespace preload and check, module reading...
func (svc module) lookup(namespaceID uint64, lookup func(*moduleActionProps) (*types.Module, error)) (m *types.Module, err error) {
	var aProps = &moduleActionProps{module: &types.Module{NamespaceID: namespaceID}}

	err = func() error {
		if ns, err := loadNamespace(svc.ctx, svc.store, namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if m, err = lookup(aProps); errors.Is(err, store.ErrNotFound) {
			return ModuleErrNotFound()
		} else if err != nil {
			return err
		}

		aProps.setModule(m)

		if !svc.ac.CanReadModule(svc.ctx, m) {
			return ModuleErrNotAllowedToRead()
		}

		return loadModuleFields(svc.ctx, svc.store, m)

	}()

	return m, svc.recordAction(svc.ctx, aProps, ModuleActionLookup, err)
}

func (svc module) uniqueCheck(m *types.Module) (err error) {
	if m.Handle != "" {
		if e, _ := store.LookupComposeModuleByNamespaceIDHandle(svc.ctx, svc.store, m.NamespaceID, m.Handle); e != nil && e.ID > 0 && e.ID != m.ID {
			return ModuleErrHandleNotUnique()
		}
	}

	if m.Name != "" {
		if e, _ := store.LookupComposeModuleByNamespaceIDName(svc.ctx, svc.store, m.NamespaceID, m.Name); e != nil && e.ID > 0 && e.ID != m.ID {
			return ModuleErrNameNotUnique()
		}
	}

	return nil
}

func (svc module) handleUpdate(upd *types.Module) moduleUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace, m *types.Module) (mch bool, fch bool, err error) {
		if isStale(upd.UpdatedAt, m.UpdatedAt, m.CreatedAt) {
			return false, false, ModuleErrStaleData()
		}

		if upd.Handle != m.Handle && !handle.IsValid(upd.Handle) {
			return false, false, ModuleErrInvalidHandle()
		}

		if err = svc.uniqueCheck(upd); err != nil {
			return false, false, err
		}

		if !svc.ac.CanUpdateModule(svc.ctx, m) {
			return false, false, ModuleErrNotAllowedToUpdate()
		}

		if m.Name != upd.Name {
			mch = true
			m.Name = upd.Name
		}

		if m.Handle != upd.Handle {
			mch = true
			m.Handle = upd.Handle
		}

		if m.Meta.String() != upd.Meta.String() {
			mch = true
			m.Meta = upd.Meta
		}

		// @todo make field-change detection more optimal
		if len(upd.Fields) > 0 {
			fch = true
			m.Fields = upd.Fields
		}

		if mch {
			m.UpdatedAt = nowPtr()
		}

		// for now, we assume that
		return mch, fch, nil
	}
}

func (svc module) handleDelete(ctx context.Context, ns *types.Namespace, m *types.Module) (bool, bool, error) {
	if !svc.ac.CanDeleteModule(ctx, m) {
		return false, false, ModuleErrNotAllowedToDelete()
	}

	if m.DeletedAt != nil {
		// module already deleted
		return false, false, nil
	}

	m.DeletedAt = nowPtr()
	return true, false, nil
}

func (svc module) handleUndelete(ctx context.Context, ns *types.Namespace, m *types.Module) (bool, bool, error) {
	if !svc.ac.CanDeleteModule(ctx, m) {
		return false, false, ModuleErrNotAllowedToUndelete()
	}

	if m.DeletedAt == nil {
		// module not deleted
		return false, false, nil
	}

	m.DeletedAt = nil
	return true, false, nil
}

// updates module fields
// expecting to receive all module fields, as it deletes the rest
// also, sort order of the fields is also important as this fn stores and updates field's place as send
func updateModuleFields(ctx context.Context, s store.Storer, m *types.Module, newFields types.ModuleFieldSet, hasRecords bool) (err error) {
	for _, f := range newFields {
		// Set module ID to all new fields
		if f.ModuleID == 0 {
			f.ModuleID = m.ID
		}

		// Make sure all updating fields belong here
		if f.ModuleID != m.ID {
			return fmt.Errorf("module id of field %q does not match the module", f.Name)
		}
	}

	if err = loadModuleFields(ctx, s, m); err != nil {
		return
	}

	for _, ef := range m.Fields {
		f := newFields.FindByID(ef.ID)
		if f != nil || f.DeletedAt == nil {
			continue
		}

		ef.DeletedAt = nowPtr()
		err = store.UpdateComposeModuleField(ctx, s, ef)
		if err != nil {
			return err
		}
	}

	for idx, f := range newFields {
		f.Place = idx
		f.DeletedAt = nil

		if e := m.Fields.FindByID(f.ID); e != nil {
			f.CreatedAt = e.CreatedAt

			// We do not have any other code in place that would handle changes of field name and kind, so we need
			// to reset any changes made to the field.
			// @todo remove when we are able to handle field rename & type change
			if hasRecords {
				f.Name = e.Name
				f.Kind = e.Kind
			}

			f.UpdatedAt = nowPtr()

			err = store.UpdateComposeModuleField(ctx, s, f)
		} else {
			f.ID = id.Next()
			f.CreatedAt = *nowPtr()
			err = store.CreateComposeModuleField(ctx, s, f)
		}

		if err != nil {
			return err
		}
	}

	return nil
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
		sort.Sort(m.Fields)
	}

	return
}

// loads record module with fields and namespace
func loadModuleWithNamespace(ctx context.Context, s store.Storer, namespaceID, moduleID uint64) (ns *types.Namespace, m *types.Module, err error) {
	if moduleID == 0 {
		return nil, nil, ModuleErrInvalidID()
	}

	if ns, err = loadNamespace(ctx, s, namespaceID); err == nil {
		m, err = loadModule(ctx, s, moduleID)
	}

	if err != nil {
		return nil, nil, err
	}

	if namespaceID != m.NamespaceID {
		// Make sure chart belongs to the right namespace
		return nil, nil, ModuleErrNotFound()
	}

	return
}

func loadModule(ctx context.Context, s store.Storer, moduleID uint64) (m *types.Module, err error) {
	if moduleID == 0 {
		return nil, ModuleErrInvalidID()
	}

	if m, err = store.LookupComposeModuleByID(ctx, s, moduleID); errors.Is(err, store.ErrNotFound) {
		err = ModuleErrNotFound()
	}

	if err == nil {
		err = loadModuleFields(ctx, s, m)
	}

	if err != nil {
		return nil, err
	}

	return
}
