package service

import (
	"context"
	"strconv"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	module struct {
		db  *factory.DB
		ctx context.Context

		actionlog actionlog.Recorder

		ac       moduleAccessController
		eventbus eventDispatcher

		moduleRepo repository.ModuleRepository
		recordRepo repository.RecordRepository
		pageRepo   repository.PageRepository
		nsRepo     repository.NamespaceRepository
	}

	moduleAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreateModule(context.Context, *types.Namespace) bool
		CanReadModule(context.Context, *types.Module) bool
		CanUpdateModule(context.Context, *types.Module) bool
		CanDeleteModule(context.Context, *types.Module) bool

		FilterReadableModules(ctx context.Context) *permissions.ResourceFilter
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
)

func Module() ModuleService {
	return (&module{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
	}).With(context.Background())
}

func (svc module) With(ctx context.Context) ModuleService {
	db := repository.DB(ctx)
	return &module{
		db:  db,
		ctx: ctx,

		actionlog: DefaultActionlog,

		ac:       svc.ac,
		eventbus: svc.eventbus,

		moduleRepo: repository.Module(ctx, db),
		recordRepo: repository.Record(ctx, db),
		pageRepo:   repository.Page(ctx, db),
		nsRepo:     repository.Namespace(ctx, db),
	}
}

// lookup fn() orchestrates module lookup, namespace preload and check, module reading...
func (svc module) lookup(namespaceID uint64, lookup func(*moduleActionProps) (*types.Module, error)) (m *types.Module, err error) {
	var aProps = &moduleActionProps{module: &types.Module{NamespaceID: namespaceID}}

	err = svc.db.Transaction(func() error {
		if ns, err := svc.loadNamespace(namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if m, err = lookup(aProps); err != nil {
			if repository.ErrModuleNotFound.Eq(err) {
				return ModuleErrNotFound()
			}

			return err
		}

		aProps.setModule(m)

		if !svc.ac.CanReadModule(svc.ctx, m) {
			return ModuleErrNotAllowedToRead()
		}

		if m.Fields, err = svc.moduleRepo.FindFields(m.ID); err != nil {
			return err
		}

		return nil
	})

	return m, svc.recordAction(svc.ctx, aProps, ModuleActionLookup, err)
}

// FindByID tries to find module by ID
func (svc module) FindByID(namespaceID, moduleID uint64) (m *types.Module, err error) {
	return svc.lookup(namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		if moduleID == 0 {
			return nil, ModuleErrInvalidID()
		}

		aProps.module.ID = moduleID
		return svc.moduleRepo.FindByID(namespaceID, moduleID)
	})
}

// FindByName tries to find module by name
func (svc module) FindByName(namespaceID uint64, name string) (m *types.Module, err error) {
	return svc.lookup(namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		aProps.module.Name = name
		return svc.moduleRepo.FindByName(namespaceID, name)
	})
}

// FindByHandle tries to find module by handle
func (svc module) FindByHandle(namespaceID uint64, h string) (m *types.Module, err error) {
	return svc.lookup(namespaceID, func(aProps *moduleActionProps) (*types.Module, error) {
		if !handle.IsValid(h) {
			return nil, ModuleErrInvalidHandle()
		}

		aProps.module.Handle = h
		return svc.moduleRepo.FindByHandle(namespaceID, h)
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

func (svc module) Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error) {
	var (
		aProps = &moduleActionProps{filter: &filter}
	)

	err = svc.db.Transaction(func() error {
		filter.IsReadable = svc.ac.FilterReadableModules(svc.ctx)

		set, f, err = svc.moduleRepo.Find(filter)
		if err != nil {
			return err
		}

		// Preload all fields and update all modules
		var ff types.ModuleFieldSet
		if ff, err = svc.moduleRepo.FindFields(set.IDs()...); err != nil {
			return err
		}

		return set.Walk(func(m *types.Module) error {
			m.Fields = ff.FilterByModule(m.ID)
			return nil
		})
	})

	return set, f, svc.recordAction(svc.ctx, aProps, ModuleActionSearch, err)
}

func (svc module) Create(new *types.Module) (m *types.Module, err error) {
	var (
		ns     *types.Namespace
		aProps = &moduleActionProps{changed: new}
	)

	err = svc.db.Transaction(func() error {
		if !handle.IsValid(new.Handle) {
			return ModuleErrInvalidHandle()
		}

		if ns, err = svc.loadNamespace(new.NamespaceID); err != nil {
			return err
		}

		if !svc.ac.CanCreateModule(svc.ctx, ns) {
			return ModuleErrNotAllowedToCreate()
		}

		aProps.setNamespace(ns)

		// Calling before-create scripts
		if err = svc.eventbus.WaitFor(svc.ctx, event.ModuleBeforeCreate(new, nil, ns)); err != nil {
			return err
		}

		if err = svc.UniqueCheck(new); err != nil {
			return err
		}

		if m, err = svc.moduleRepo.Create(new); err != nil {
			return err
		}

		aProps.setModule(m)

		if err = svc.moduleRepo.UpdateFields(m.ID, m.Fields, false); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.ModuleAfterCreate(m, nil, ns))
		return nil
	})

	return m, svc.recordAction(svc.ctx, aProps, ModuleActionCreate, err)
}

func (svc module) Update(upd *types.Module) (m *types.Module, err error) {
	var (
		ns     *types.Namespace
		aProps = &moduleActionProps{changed: upd}
	)

	err = svc.db.Transaction(func() error {

		if upd.ID == 0 {
			return ModuleErrInvalidID()
		}

		if !handle.IsValid(upd.Handle) {
			return ModuleErrInvalidHandle()
		}

		if ns, err = svc.loadNamespace(upd.NamespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if m, err = svc.moduleRepo.FindByID(upd.NamespaceID, upd.ID); err != nil {
			return err
		}

		aProps.setModule(m)

		if isStale(upd.UpdatedAt, m.UpdatedAt, m.CreatedAt) {
			return ModuleErrStaleData()
		}

		if !svc.ac.CanUpdateModule(svc.ctx, m) {
			return ModuleErrNotAllowedToUpdate()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.ModuleBeforeUpdate(upd, m, ns)); err != nil {
			return err
		}

		if err = svc.UniqueCheck(upd); err != nil {
			return err
		}

		m.Name = upd.Name
		m.Handle = upd.Handle
		m.Meta = upd.Meta
		m.Fields = upd.Fields

		if m, err = svc.moduleRepo.Update(m); err != nil {
			return err
		}

		// select 1 record to see how fields can be updated
		var rf = types.RecordFilter{}
		rf.Limit = 1
		if _, rf, err = svc.recordRepo.Find(m, rf); err != nil {
			return err
		}

		if err = svc.moduleRepo.UpdateFields(m.ID, m.Fields, rf.Count > 0); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.ModuleAfterUpdate(upd, m, ns))
		return nil
	})

	return m, svc.recordAction(svc.ctx, aProps, ModuleActionUpdate, err)
}

func (svc module) DeleteByID(namespaceID, moduleID uint64) (err error) {
	var (
		m      *types.Module
		ns     *types.Namespace
		aProps = &moduleActionProps{module: &types.Module{ID: moduleID, NamespaceID: namespaceID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if moduleID == 0 {
			return ModuleErrInvalidID()
		}

		if ns, err = svc.loadNamespace(namespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if m, err = svc.moduleRepo.FindByID(namespaceID, moduleID); err != nil {
			if repository.ErrModuleNotFound.Eq(err) {
				return ModuleErrNotFound()
			}

			return err
		} else if !svc.ac.CanDeleteModule(svc.ctx, m) {
			return ModuleErrNotAllowedToDelete()
		}

		aProps.setChanged(m)

		if err = svc.eventbus.WaitFor(svc.ctx, event.ModuleBeforeDelete(nil, m, ns)); err != nil {
			return err
		}

		if err = svc.moduleRepo.DeleteByID(namespaceID, moduleID); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.ModuleAfterDelete(nil, m, ns))
		return err
	})

	return svc.recordAction(svc.ctx, aProps, ModuleActionDelete, err)

}

func (svc module) UniqueCheck(m *types.Module) (err error) {
	if m.Handle != "" {
		if e, _ := svc.moduleRepo.FindByHandle(m.NamespaceID, m.Handle); e != nil && e.ID > 0 && e.ID != m.ID {
			return ModuleErrHandleNotUnique()
		}
	}

	if m.Name != "" {
		if e, _ := svc.moduleRepo.FindByName(m.NamespaceID, m.Name); e != nil && e.ID > 0 && e.ID != m.ID {
			return ModuleErrNameNotUnique()
		}
	}

	return nil
}

// Namespace loader
//
func (svc module) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
	var (
		moProps = &moduleActionProps{namespace: &types.Namespace{ID: namespaceID}}
	)

	if namespaceID == 0 {
		return nil, ModuleErrInvalidID(moProps)
	}

	if ns, err = svc.nsRepo.FindByID(namespaceID); err != nil {
		if repository.ErrNamespaceNotFound.Eq(err) {
			return nil, ModuleErrNamespaceNotFound()
		}

		return nil, err
	}

	moProps.setNamespace(ns)

	if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, ModuleErrNotAllowedToReadNamespace(moProps)
	}

	return
}
