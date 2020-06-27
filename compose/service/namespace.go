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
	namespace struct {
		db  *factory.DB
		ctx context.Context

		actionlog actionlog.Recorder

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
		FindByAny(interface{}) (*types.Namespace, error)

		Create(namespace *types.Namespace) (*types.Namespace, error)
		Update(namespace *types.Namespace) (*types.Namespace, error)
		DeleteByID(namespaceID uint64) error
	}
)

func Namespace() NamespaceService {
	return (&namespace{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
	}).With(context.Background())
}

func (svc namespace) With(ctx context.Context) NamespaceService {
	db := repository.DB(ctx)
	return &namespace{
		db:  db,
		ctx: ctx,

		actionlog: DefaultActionlog,

		ac:       svc.ac,
		eventbus: svc.eventbus,

		namespaceRepo: repository.Namespace(ctx, db),
	}
}

// lookup fn() orchestrates namespace lookup, and check
func (svc namespace) lookup(lookup func(*namespaceActionProps) (*types.Namespace, error)) (m *types.Namespace, err error) {
	var aProps = &namespaceActionProps{namespace: &types.Namespace{}}

	err = svc.db.Transaction(func() error {
		if m, err = lookup(aProps); err != nil {
			if repository.ErrNamespaceNotFound.Eq(err) {
				return NamespaceErrNotFound()
			}

			return err
		}

		aProps.setNamespace(m)

		if !svc.ac.CanReadNamespace(svc.ctx, m) {
			return NamespaceErrNotAllowedToRead()
		}

		return nil
	})

	return m, svc.recordAction(svc.ctx, aProps, NamespaceActionLookup, err)
}

func (svc namespace) FindByID(ID uint64) (ns *types.Namespace, err error) {
	return svc.lookup(func(aProps *namespaceActionProps) (*types.Namespace, error) {
		if ID == 0 {
			return nil, NamespaceErrInvalidID()
		}

		aProps.namespace.ID = ID
		return svc.namespaceRepo.FindByID(ID)
	})
}

// FindByHandle is an alias for FindBySlug
func (svc namespace) FindByHandle(handle string) (ns *types.Namespace, err error) {
	return svc.FindBySlug(handle)
}

func (svc namespace) FindBySlug(slug string) (ns *types.Namespace, err error) {
	return svc.lookup(func(aProps *namespaceActionProps) (*types.Namespace, error) {
		if !handle.IsValid(slug) {
			return nil, NamespaceErrInvalidHandle()
		}

		aProps.namespace.Slug = slug
		return svc.namespaceRepo.FindBySlug(slug)
	})
}

// FindByAny tries to find namespace by id, handle or slug
func (svc namespace) FindByAny(identifier interface{}) (r *types.Namespace, err error) {
	if ID, ok := identifier.(uint64); ok {
		r, err = svc.FindByID(ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			r, err = svc.FindByID(ID)
		} else {
			r, err = svc.FindByHandle(strIdentifier)
			if err == nil && r.ID == 0 {
				r, err = svc.FindBySlug(strIdentifier)
			}
		}
	} else {
		err = NamespaceErrInvalidID()
	}

	if err != nil {
		return
	}

	return
}

// search fn() orchestrates pages search, namespace preload and check
func (svc namespace) search(filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	var (
		aProps = &namespaceActionProps{filter: &filter}
	)

	f = filter
	f.IsReadable = svc.ac.FilterReadableNamespaces(svc.ctx)

	err = svc.db.Transaction(func() error {
		if set, f, err = svc.namespaceRepo.Find(f); err != nil {
			return err
		}

		return nil
	})

	return set, f, svc.recordAction(svc.ctx, aProps, NamespaceActionSearch, err)
}

func (svc namespace) Find(filter types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error) {
	return svc.search(filter)
}

// Create adds namespace and presets access rules for role everyone
func (svc namespace) Create(new *types.Namespace) (ns *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{changed: new}
	)

	err = svc.db.Transaction(func() (err error) {
		if !handle.IsValid(new.Slug) {
			return NamespaceErrInvalidHandle()
		}

		if !svc.ac.CanCreateNamespace(svc.ctx) {
			return NamespaceErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceBeforeCreate(new, nil)); err != nil {
			return
		}

		if err = svc.UniqueCheck(new); err != nil {
			return
		}

		if ns, err = svc.namespaceRepo.Create(new); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.NamespaceAfterCreate(ns, nil))
		return
	})

	return ns, svc.recordAction(svc.ctx, aProps, NamespaceActionCreate, err)
}

func (svc namespace) Update(upd *types.Namespace) (ns *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{changed: upd}
	)

	err = svc.db.Transaction(func() (err error) {
		if upd.ID == 0 {
			return NamespaceErrInvalidID()
		}

		if !handle.IsValid(upd.Slug) {
			return NamespaceErrInvalidHandle()
		}

		if ns, err = svc.FindByID(upd.ID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if isStale(upd.UpdatedAt, ns.UpdatedAt, ns.CreatedAt) {
			return NamespaceErrStaleData()
		}

		if !svc.ac.CanUpdateNamespace(svc.ctx, ns) {
			return NamespaceErrNotAllowedToUpdate()
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
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.NamespaceAfterUpdate(upd, ns))
		return
	})

	return ns, svc.recordAction(svc.ctx, aProps, NamespaceActionUpdate, err)
}

func (svc namespace) DeleteByID(namespaceID uint64) (err error) {
	var (
		del    *types.Namespace
		aProps = &namespaceActionProps{namespace: &types.Namespace{ID: namespaceID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if namespaceID == 0 {
			return NamespaceErrInvalidID()
		}

		if del, err = svc.namespaceRepo.FindByID(namespaceID); err != nil {
			if repository.ErrNamespaceNotFound.Eq(err) {
				return NamespaceErrNotFound()
			}

			return
		}

		aProps.setChanged(del)

		if !svc.ac.CanDeleteNamespace(svc.ctx, del) {
			return NamespaceErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceBeforeDelete(nil, del)); err != nil {
			return
		}

		if err = svc.namespaceRepo.DeleteByID(namespaceID); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.NamespaceAfterDelete(nil, del))
		return
	})

	return svc.recordAction(svc.ctx, aProps, NamespaceActionDelete, err)
}

func (svc namespace) UniqueCheck(ns *types.Namespace) (err error) {
	if ns.Slug != "" {
		if e, _ := svc.namespaceRepo.FindBySlug(ns.Slug); e != nil && e.ID != ns.ID {
			return NamespaceErrHandleNotUnique()
		}
	}

	return nil
}
