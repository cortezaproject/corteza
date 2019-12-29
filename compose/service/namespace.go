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
	namespace struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac       namespaceAccessController
		eventbus eventDispatcher

		namespaceRepo repository.NamespaceRepository
	}

	namespaceAccessController interface {
		CanCreateNamespace(context.Context) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanUpdateNamespace(context.Context, *types.Namespace) bool
		CanDeleteNamespace(context.Context, *types.Namespace) bool

		Grant(ctx context.Context, rr ...*permissions.Rule) error

		FilterReadableNamespaces(ctx context.Context) *permissions.ResourceFilter
	}

	NamespaceService interface {
		With(ctx context.Context) NamespaceService

		FindByID(namespaceID uint64) (*types.Namespace, error)
		FindByHandle(handle string) (*types.Namespace, error)
		Find(types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error)

		Create(namespace *types.Namespace) (*types.Namespace, error)
		Update(namespace *types.Namespace) (*types.Namespace, error)
		DeleteByID(namespaceID uint64) error
	}
)

func Namespace() NamespaceService {
	return (&namespace{
		logger:   DefaultLogger.Named("namespace"),
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
	}).With(context.Background())
}

func (svc namespace) With(ctx context.Context) NamespaceService {
	db := repository.DB(ctx)
	return &namespace{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac:       svc.ac,
		eventbus: svc.eventbus,

		namespaceRepo: repository.Namespace(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc namespace) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc namespace) FindByID(ID uint64) (ns *types.Namespace, err error) {
	return svc.checkPermissions(svc.namespaceRepo.FindByID(ID))
}

func (svc namespace) FindByHandle(handle string) (ns *types.Namespace, err error) {
	return svc.checkPermissions(svc.namespaceRepo.FindBySlug(handle))
}

func (svc namespace) FindBySlug(slug string) (ns *types.Namespace, err error) {
	return svc.checkPermissions(svc.namespaceRepo.FindBySlug(slug))
}

func (svc namespace) checkPermissions(p *types.Namespace, err error) (*types.Namespace, error) {
	if err != nil {
		return nil, err
	} else if !svc.ac.CanReadNamespace(svc.ctx, p) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return p, err
}

func (svc namespace) Find(filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	filter.IsReadable = svc.ac.FilterReadableNamespaces(svc.ctx)

	set, f, err = svc.namespaceRepo.Find(filter)
	if err != nil {
		return
	}

	return
}

// Create adds namespace and presets access rules for role everyone
func (svc namespace) Create(new *types.Namespace) (ns *types.Namespace, err error) {
	if !handle.IsValid(new.Slug) {
		return nil, ErrInvalidHandle
	}

	if !svc.ac.CanCreateNamespace(svc.ctx) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	if err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceBeforeCreate(new, nil)); err != nil {
		return
	}

	if err = svc.UniqueCheck(new); err != nil {
		return
	}

	if ns, err = svc.namespaceRepo.Create(new); err != nil {
		return nil, err
	}

	defer svc.eventbus.Dispatch(svc.ctx, event.NamespaceAfterCreate(ns, nil))
	return
}

func (svc namespace) Update(upd *types.Namespace) (ns *types.Namespace, err error) {
	if upd.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if !handle.IsValid(upd.Slug) {
		return nil, ErrInvalidHandle
	}

	if ns, err = svc.FindByID(upd.ID); err != nil {
		return nil, err
	}

	if isStale(upd.UpdatedAt, ns.UpdatedAt, ns.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if !svc.ac.CanUpdateNamespace(svc.ctx, ns) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	if err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceBeforeUpdate(upd, ns)); err != nil {
		return
	}

	if err = svc.UniqueCheck(upd); err != nil {
		return
	}

	// Copy changes
	ns.Name = upd.Name
	ns.Slug = upd.Slug
	ns.Meta = upd.Meta
	ns.Enabled = upd.Enabled

	if ns, err = svc.namespaceRepo.Update(ns); err != nil {
		return nil, err
	}

	defer svc.eventbus.Dispatch(svc.ctx, event.NamespaceAfterUpdate(upd, ns))
	return
}

func (svc namespace) DeleteByID(namespaceID uint64) (err error) {
	var (
		del *types.Namespace
	)

	if namespaceID == 0 {
		return ErrInvalidID.withStack()
	}

	if del, err = svc.namespaceRepo.FindByID(namespaceID); err != nil {
		return
	} else if !svc.ac.CanDeleteNamespace(svc.ctx, del) {
		return ErrNoDeletePermissions.withStack()
	}

	if err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceBeforeDelete(nil, del)); err != nil {
		return
	}

	if err = svc.namespaceRepo.DeleteByID(namespaceID); err != nil {
		return
	}

	defer svc.eventbus.Dispatch(svc.ctx, event.NamespaceAfterDelete(nil, del))
	return
}

func (svc namespace) UniqueCheck(ns *types.Namespace) (err error) {
	if ns.Slug != "" {
		if e, _ := svc.namespaceRepo.FindBySlug(ns.Slug); e != nil && e.ID != ns.ID {
			return repository.ErrNamespaceSlugNotUnique
		}
	}

	return nil
}
