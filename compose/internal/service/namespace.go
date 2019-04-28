package service

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/compose/types"
)

type (
	namespace struct {
		db  *factory.DB
		ctx context.Context

		prmSvc PermissionsService

		namespaceRepo repository.NamespaceRepository
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
		prmSvc: DefaultPermissions,
	}).With(context.Background())
}

func (svc *namespace) With(ctx context.Context) NamespaceService {
	db := repository.DB(ctx)
	return &namespace{
		db:  db,
		ctx: ctx,

		prmSvc: svc.prmSvc.With(ctx),

		namespaceRepo: repository.Namespace(ctx, db),
	}
}

func (svc *namespace) FindByID(ID uint64) (n *types.Namespace, err error) {
	if ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if n, err = svc.namespaceRepo.FindByID(ID); err != nil {
		return
	} else if !svc.prmSvc.CanReadNamespace(n) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}

func (svc *namespace) Find(filter types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error) {
	nn, f, err := svc.namespaceRepo.Find(filter)
	if err != nil {
		return nil, f, err
	}

	nn, _ = nn.Filter(func(m *types.Namespace) (bool, error) {
		return svc.prmSvc.CanReadNamespace(m), nil
	})

	return nn, f, nil
}

func (svc *namespace) Create(mod *types.Namespace) (*types.Namespace, error) {
	if !svc.prmSvc.CanCreateNamespace() {
		return nil, ErrNoCreatePermissions.withStack()
	}

	return svc.namespaceRepo.Create(mod)
}

func (svc *namespace) Update(updated *types.Namespace) (m *types.Namespace, err error) {
	m, err = svc.FindByID(updated.ID)
	if err != nil {
		return nil, err
	}

	if isStale(updated.UpdatedAt, m.UpdatedAt, m.CreatedAt) {
		return nil, ErrStaleData
	}

	if !svc.prmSvc.CanUpdateNamespace(m) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	m.Name = updated.Name
	m.Slug = updated.Slug
	m.Meta = updated.Meta
	m.Enabled = updated.Enabled

	return svc.namespaceRepo.Update(m)
}

func (svc *namespace) DeleteByID(ID uint64) error {
	if m, err := svc.namespaceRepo.FindByID(ID); err != nil {
		return err
	} else if !svc.prmSvc.CanDeleteNamespace(m) {
		return ErrNoDeletePermissions.withStack()
	}

	return svc.namespaceRepo.DeleteByID(ID)
}

// Data is stale when new date does not match updatedAt or createdAt (before first update)
func isStale(new *time.Time, updatedAt *time.Time, createdAt time.Time) bool {

	if new == nil {
		return false
	}

	if updatedAt != nil {
		return !new.Equal(*updatedAt)
	}

	return new.Equal(createdAt)
}
