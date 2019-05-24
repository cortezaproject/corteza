package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/internal/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	module struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac moduleAccessController

		moduleRepo repository.ModuleRepository
		pageRepo   repository.PageRepository
		nsRepo     repository.NamespaceRepository
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
		pageRepo:   repository.Page(ctx, db),
		nsRepo:     repository.Namespace(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc module) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc module) FindByID(namespaceID, moduleID uint64) (m *types.Module, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	if m, err = svc.moduleRepo.FindByID(namespaceID, moduleID); err != nil {
		return
	} else if !svc.ac.CanReadModule(svc.ctx, m) {
		return nil, ErrNoReadPermissions.withStack()
	}

	var ff types.ModuleFieldSet
	if ff, err = svc.moduleRepo.FindFields(m.ID); err != nil {
		return
	} else {
		_ = ff.Walk(func(f *types.ModuleField) error {
			m.Fields = append(m.Fields, f)
			return nil
		})
	}

	return
}

func (svc module) Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error) {
	set, f, err = svc.moduleRepo.Find(filter)
	if err != nil {
		return
	}

	set, _ = set.Filter(func(m *types.Module) (bool, error) {
		return svc.ac.CanReadModule(svc.ctx, m), nil
	})

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

func (svc module) Create(mod *types.Module) (*types.Module, error) {

	if mod.NamespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	if ns, err := svc.loadNamespace(mod.NamespaceID); err != nil {
		return nil, err
	} else if !svc.ac.CanCreateModule(svc.ctx, ns) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	return svc.moduleRepo.Create(mod)
}

func (svc module) Update(mod *types.Module) (m *types.Module, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if m, err = svc.moduleRepo.FindByID(mod.NamespaceID, mod.ID); err != nil {
		return
	}

	if isStale(mod.UpdatedAt, m.UpdatedAt, m.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if !svc.ac.CanUpdateModule(svc.ctx, m) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	m.Name = mod.Name
	m.Meta = mod.Meta
	m.Fields = mod.Fields

	return svc.moduleRepo.Update(m)
}

func (svc module) DeleteByID(namespaceID, moduleID uint64) error {
	if moduleID == 0 {
		return ErrInvalidID.withStack()
	}

	if namespaceID == 0 {
		return ErrNamespaceRequired.withStack()
	}

	if c, err := svc.moduleRepo.FindByID(namespaceID, moduleID); err != nil {
		return err
	} else if !svc.ac.CanDeleteModule(svc.ctx, c) {
		return ErrNoDeletePermissions.withStack()
	}

	return svc.moduleRepo.DeleteByID(namespaceID, moduleID)
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
