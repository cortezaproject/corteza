package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/internal/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	namespace struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac namespaceAccessController

		namespaceRepo repository.NamespaceRepository
	}

	namespaceAccessController interface {
		CanCreateNamespace(context.Context) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanUpdateNamespace(context.Context, *types.Namespace) bool
		CanDeleteNamespace(context.Context, *types.Namespace) bool
	}

	NamespaceService interface {
		With(ctx context.Context) NamespaceService

		FindByID(namespaceID uint64) (*types.Namespace, error)
		Find(types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error)

		Create(namespace *types.Namespace) (*types.Namespace, error)
		Update(namespace *types.Namespace) (*types.Namespace, error)
		DeleteByID(namespaceID uint64) error
	}
)

func Namespace() NamespaceService {
	return (&namespace{
		logger: DefaultLogger.Named("namespace"),
		ac:     DefaultAccessControl,
	}).With(context.Background())
}

func (svc namespace) With(ctx context.Context) NamespaceService {
	db := repository.DB(ctx)
	return &namespace{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac: svc.ac,

		namespaceRepo: repository.Namespace(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc namespace) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc namespace) FindByID(ID uint64) (ns *types.Namespace, err error) {
	if ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if ns, err = svc.namespaceRepo.FindByID(ID); err != nil {
		return
	} else if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}

func (svc namespace) Find(filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	set, f, err = svc.namespaceRepo.Find(filter)
	if err != nil {
		return
	}

	set, _ = set.Filter(func(ns *types.Namespace) (bool, error) {
		return svc.ac.CanReadNamespace(svc.ctx, ns), nil
	})

	return
}

func (svc namespace) Create(mod *types.Namespace) (*types.Namespace, error) {
	if !svc.ac.CanCreateNamespace(svc.ctx) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	return svc.namespaceRepo.Create(mod)
}

func (svc namespace) Update(mod *types.Namespace) (ns *types.Namespace, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	ns, err = svc.FindByID(mod.ID)
	if err != nil {
		return nil, err
	}

	if isStale(mod.UpdatedAt, ns.UpdatedAt, ns.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if !svc.ac.CanUpdateNamespace(svc.ctx, ns) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	ns.Name = mod.Name
	ns.Slug = mod.Slug
	ns.Meta = mod.Meta
	ns.Enabled = mod.Enabled

	return svc.namespaceRepo.Update(ns)
}

func (svc namespace) DeleteByID(namespaceID uint64) error {
	if namespaceID == 0 {
		return ErrInvalidID.withStack()
	}

	if ns, err := svc.namespaceRepo.FindByID(namespaceID); err != nil {
		return err
	} else if !svc.ac.CanDeleteNamespace(svc.ctx, ns) {
		return ErrNoDeletePermissions.withStack()
	}

	return svc.namespaceRepo.DeleteByID(namespaceID)
}
