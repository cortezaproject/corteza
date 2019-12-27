package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	module struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac moduleAccessController

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
		Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)

		Create(module *types.Module) (*types.Module, error)
		Update(module *types.Module) (*types.Module, error)
		DeleteByID(namespaceID, moduleID uint64) error
	}
)

func Module() ModuleService {
	return (&module{
		logger: DefaultLogger.Named("module"),
		ac:     DefaultAccessControl,
	}).With(context.Background())
}

func (svc module) With(ctx context.Context) ModuleService {
	db := repository.DB(ctx)
	return &module{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac: svc.ac,

		moduleRepo: repository.Module(ctx, db),
		recordRepo: repository.Record(ctx, db),
		pageRepo:   repository.Page(ctx, db),
		nsRepo:     repository.Namespace(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc module) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc module) FindByID(namespaceID, moduleID uint64) (m *types.Module, err error) {
	if _, err = svc.loadNamespace(namespaceID); err != nil {
		return
	} else {
		return svc.loader(svc.moduleRepo.FindByID(namespaceID, moduleID))
	}
}

func (svc module) FindByName(namespaceID uint64, name string) (m *types.Module, err error) {
	if _, err = svc.loadNamespace(namespaceID); err != nil {
		return
	} else {
		return svc.loader(svc.moduleRepo.FindByName(namespaceID, name))
	}
}

func (svc module) FindByHandle(namespaceID uint64, handle string) (m *types.Module, err error) {
	if _, err = svc.loadNamespace(namespaceID); err != nil {
		return
	} else {
		return svc.loader(svc.moduleRepo.FindByHandle(namespaceID, handle))
	}
}

func (svc module) loader(m *types.Module, err error) (*types.Module, error) {
	if err != nil {
		return nil, err
	} else if !svc.ac.CanReadModule(svc.ctx, m) {
		return nil, ErrNoReadPermissions.withStack()
	}

	var ff types.ModuleFieldSet
	if ff, err = svc.moduleRepo.FindFields(m.ID); err != nil {
		return nil, err
	} else {
		_ = ff.Walk(func(f *types.ModuleField) error {
			m.Fields = append(m.Fields, f)
			return nil
		})
	}

	return m, err
}

func (svc module) Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error) {
	filter.IsReadable = svc.ac.FilterReadableModules(svc.ctx)

	set, f, err = svc.moduleRepo.Find(filter)
	if err != nil {
		return
	}

	// Preload all fields and update all modules
	var ff types.ModuleFieldSet
	if ff, err = svc.moduleRepo.FindFields(set.IDs()...); err != nil {
		return
	} else {
		_ = ff.Walk(func(f *types.ModuleField) error {
			set.FindByID(f.ModuleID).Fields = append(set.FindByID(f.ModuleID).Fields, f)
			return nil
		})
	}

	return
}

func (svc module) Create(new *types.Module) (m *types.Module, err error) {
	var (
		ns *types.Namespace
	)

	if !handle.IsValid(new.Handle) {
		return nil, ErrInvalidHandle
	}
	if new.NamespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}
	if ns, err = svc.loadNamespace(new.NamespaceID); err != nil {
		return nil, err
	} else if !svc.ac.CanCreateModule(svc.ctx, ns) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	// Calling before-create scripts
	if err = eventbus.WaitFor(svc.ctx, event.ModuleBeforeCreate(new, nil, ns)); err != nil {
		return
	}

	if err := svc.UniqueCheck(new); err != nil {
		return nil, err
	}

	if m, err = svc.moduleRepo.Create(new); err != nil {
		return nil, err
	}

	err = svc.moduleRepo.UpdateFields(m.ID, m.Fields, false)
	if err != nil {
		return nil, err
	}

	defer eventbus.Dispatch(svc.ctx, event.ModuleAfterCreate(m, nil, ns))
	return
}

func (svc module) Update(upd *types.Module) (m *types.Module, err error) {
	var (
		ns *types.Namespace
	)

	if upd.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if !handle.IsValid(upd.Handle) {
		return nil, ErrInvalidHandle
	}

	if m, err = svc.moduleRepo.FindByID(upd.NamespaceID, upd.ID); err != nil {
		return
	}

	if ns, err = svc.loadNamespace(upd.NamespaceID); err != nil {
		return nil, err
	}

	if isStale(upd.UpdatedAt, m.UpdatedAt, m.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if !svc.ac.CanUpdateModule(svc.ctx, m) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	if err = eventbus.WaitFor(svc.ctx, event.ModuleBeforeUpdate(upd, m, ns)); err != nil {
		return
	}

	if err = svc.UniqueCheck(upd); err != nil {
		return
	}

	m.Name = upd.Name
	m.Handle = upd.Handle
	m.Meta = upd.Meta
	m.Fields = upd.Fields

	m, err = svc.moduleRepo.Update(m)
	if err != nil {
		return nil, err
	}

	_, ff, err := svc.recordRepo.Find(m, types.RecordFilter{})
	if err != nil {
		return nil, err
	}
	err = svc.moduleRepo.UpdateFields(m.ID, m.Fields, ff.Count > 0)
	if err != nil {
		return nil, err
	}

	defer eventbus.Dispatch(svc.ctx, event.ModuleAfterUpdate(upd, m, ns))
	return
}

func (svc module) DeleteByID(namespaceID, moduleID uint64) (err error) {
	var (
		del *types.Module
		ns  *types.Namespace
	)

	if moduleID == 0 {
		return ErrInvalidID.withStack()
	}

	if namespaceID == 0 {
		return ErrNamespaceRequired.withStack()
	}

	if ns, err = svc.loadNamespace(namespaceID); err != nil {
		return err
	}

	if del, err = svc.moduleRepo.FindByID(namespaceID, moduleID); err != nil {
		return err
	} else if !svc.ac.CanDeleteModule(svc.ctx, del) {
		return ErrNoDeletePermissions.withStack()
	}

	if err = eventbus.WaitFor(svc.ctx, event.ModuleBeforeDelete(nil, del, ns)); err != nil {
		return
	}

	if err = svc.moduleRepo.DeleteByID(namespaceID, moduleID); err != nil {
		return
	}

	defer eventbus.Dispatch(svc.ctx, event.ModuleAfterDelete(nil, del, ns))
	return
}

func (svc module) UniqueCheck(m *types.Module) (err error) {
	if m.Handle != "" {
		if e, _ := svc.moduleRepo.FindByHandle(m.NamespaceID, m.Handle); e != nil && e.ID > 0 && e.ID != m.ID {
			return repository.ErrModuleHandleNotUnique
		}
	}

	if m.Name != "" {
		if e, _ := svc.moduleRepo.FindByName(m.NamespaceID, m.Name); e != nil && e.ID > 0 && e.ID != m.ID {
			return repository.ErrModuleNameNotUnique
		}
	}

	return nil
}

func (svc module) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	if ns, err = svc.nsRepo.FindByID(namespaceID); err != nil {
		return
	}

	if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}
