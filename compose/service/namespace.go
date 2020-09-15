package service

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/store"
	"strconv"
)

type (
	namespace struct {
		ctx       context.Context
		actionlog actionlog.Recorder
		ac        namespaceAccessController
		eventbus  eventDispatcher
		store     store.Storer
	}

	namespaceAccessController interface {
		CanCreateNamespace(context.Context) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanUpdateNamespace(context.Context, *types.Namespace) bool
		CanDeleteNamespace(context.Context, *types.Namespace) bool

		Grant(ctx context.Context, rr ...*permissions.Rule) error
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

	namespaceUpdateHandler func(ctx context.Context, ns *types.Namespace) (bool, error)
)

func Namespace() NamespaceService {
	return (&namespace{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
	}).With(context.Background())
}

func (svc namespace) With(ctx context.Context) NamespaceService {
	return &namespace{
		ctx:       ctx,
		actionlog: DefaultActionlog,
		ac:        svc.ac,
		eventbus:  svc.eventbus,
		store:     DefaultStore,
	}
}

// search fn() orchestrates pages search, namespace preload and check
func (svc namespace) Find(filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	var (
		aProps = &namespaceActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Namespace) (bool, error) {
		if !svc.ac.CanReadNamespace(svc.ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if set, f, err = store.SearchComposeNamespaces(svc.ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(svc.ctx, aProps, NamespaceActionSearch, err)
}

func (svc namespace) FindByID(ID uint64) (ns *types.Namespace, err error) {
	return svc.lookup(func(aProps *namespaceActionProps) (*types.Namespace, error) {
		if ID == 0 {
			return nil, NamespaceErrInvalidID()
		}

		aProps.namespace.ID = ID
		return store.LookupComposeNamespaceByID(svc.ctx, svc.store, ID)
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
		return store.LookupComposeNamespaceBySlug(svc.ctx, svc.store, slug)
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

// Create adds namespace and presets access rules for role everyone
func (svc namespace) Create(new *types.Namespace) (ns *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{changed: new}
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if !handle.IsValid(new.Slug) {
			return NamespaceErrInvalidHandle()
		}

		if !svc.ac.CanCreateNamespace(svc.ctx) {
			return NamespaceErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceBeforeCreate(new, nil)); err != nil {
			return err
		}

		if err = svc.uniqueCheck(new); err != nil {
			return err
		}

		if err = store.CreateComposeNamespace(svc.ctx, svc.store, new); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.NamespaceAfterCreate(ns, nil))
		return nil
	})

	return ns, svc.recordAction(svc.ctx, aProps, NamespaceActionCreate, err)
}

func (svc namespace) Update(upd *types.Namespace) (c *types.Namespace, err error) {
	return svc.updater(upd.ID, NamespaceActionUpdate, svc.handleUpdate(upd))
}

func (svc namespace) DeleteByID(namespaceID uint64) error {
	return trim1st(svc.updater(namespaceID, NamespaceActionDelete, svc.handleDelete))
}

func (svc namespace) UndeleteByID(namespaceID uint64) error {
	return trim1st(svc.updater(namespaceID, NamespaceActionUndelete, svc.handleUndelete))
}

func (svc namespace) updater(namespaceID uint64, action func(...*namespaceActionProps) *namespaceAction, fn namespaceUpdateHandler) (*types.Namespace, error) {
	var (
		changed bool
		ns, old *types.Namespace
		aProps  = &namespaceActionProps{namespace: &types.Namespace{ID: namespaceID}}
		err     error
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, err = loadNamespace(svc.ctx, s, namespaceID)
		if err != nil {
			return
		}

		old = ns.Clone()

		aProps.setNamespace(ns)
		aProps.setChanged(ns)

		if ns.DeletedAt == nil {
			err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceBeforeUpdate(ns, old))
		} else {
			err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceBeforeDelete(ns, old))
		}

		if err != nil {
			return
		}

		if changed, err = fn(svc.ctx, ns); err != nil {
			return err
		}

		if changed {
			if err = store.UpdateComposeNamespace(svc.ctx, svc.store, ns); err != nil {
				return err
			}
		}

		if ns.DeletedAt == nil {
			err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceAfterUpdate(ns, old))
		} else {
			err = svc.eventbus.WaitFor(svc.ctx, event.NamespaceAfterDelete(nil, old))
		}

		return err
	})

	return ns, svc.recordAction(svc.ctx, aProps, action, err)
}

// lookup fn() orchestrates namespace lookup, and check
func (svc namespace) lookup(lookup func(*namespaceActionProps) (*types.Namespace, error)) (ns *types.Namespace, err error) {
	var aProps = &namespaceActionProps{namespace: &types.Namespace{}}

	err = func() error {
		if ns, err = lookup(aProps); errors.Is(err, store.ErrNotFound) {
			return NamespaceErrNotFound()
		} else if err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if !svc.ac.CanReadNamespace(svc.ctx, ns) {
			return NamespaceErrNotAllowedToRead()
		}

		return nil
	}()

	return ns, svc.recordAction(svc.ctx, aProps, NamespaceActionLookup, err)
}

func (svc namespace) uniqueCheck(ns *types.Namespace) (err error) {
	if ns.Slug != "" {
		if e, _ := store.LookupComposeNamespaceBySlug(svc.ctx, svc.store, ns.Slug); e != nil && e.ID != ns.ID {
			return NamespaceErrHandleNotUnique()
		}
	}

	return nil
}

func (svc namespace) handleUpdate(upd *types.Namespace) namespaceUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace) (bool, error) {
		if isStale(upd.UpdatedAt, ns.UpdatedAt, ns.CreatedAt) {
			return false, NamespaceErrStaleData()
		}

		if upd.Slug != ns.Slug && !handle.IsValid(upd.Slug) {
			return false, NamespaceErrInvalidHandle()
		}

		if err := svc.uniqueCheck(upd); err != nil {
			return false, err
		}

		if !svc.ac.CanUpdateNamespace(svc.ctx, ns) {
			return false, NamespaceErrNotAllowedToUpdate()
		}

		ns.Name = upd.Name
		ns.Slug = upd.Slug
		ns.Meta = upd.Meta
		ns.Enabled = upd.Enabled
		ns.UpdatedAt = nowPtr()

		return true, nil
	}
}

func (svc namespace) handleDelete(ctx context.Context, ns *types.Namespace) (bool, error) {
	if !svc.ac.CanDeleteNamespace(ctx, ns) {
		return false, NamespaceErrNotAllowedToDelete()
	}

	if ns.DeletedAt != nil {
		// namespace already deleted
		return false, nil
	}

	ns.DeletedAt = nowPtr()
	return true, nil
}

func (svc namespace) handleUndelete(ctx context.Context, ns *types.Namespace) (bool, error) {
	if !svc.ac.CanDeleteNamespace(ctx, ns) {
		return false, NamespaceErrNotAllowedToUndelete()
	}

	if ns.DeletedAt == nil {
		// namespace not deleted
		return false, nil
	}

	ns.DeletedAt = nil
	return true, nil
}

func loadNamespace(ctx context.Context, s store.Storer, namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, ChartErrInvalidNamespaceID()
	}

	if ns, err = store.LookupComposeNamespaceByID(ctx, s, namespaceID); errors.Is(err, store.ErrNotFound) {
		return nil, NamespaceErrNotFound()
	}

	return
}
