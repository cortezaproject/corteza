package service

import (
	"context"
	"reflect"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	namespace struct {
		actionlog actionlog.Recorder
		ac        namespaceAccessController
		eventbus  eventDispatcher
		store     store.Storer
	}

	namespaceAccessController interface {
		CanSearchNamespaces(context.Context) bool
		CanCreateNamespace(context.Context) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanUpdateNamespace(context.Context, *types.Namespace) bool
		CanDeleteNamespace(context.Context, *types.Namespace) bool

		Grant(ctx context.Context, rr ...*rbac.Rule) error
	}

	NamespaceService interface {
		FindByID(ctx context.Context, namespaceID uint64) (*types.Namespace, error)
		FindByHandle(ctx context.Context, handle string) (*types.Namespace, error)
		Find(context.Context, types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error)
		FindByAny(context.Context, interface{}) (*types.Namespace, error)

		Create(ctx context.Context, namespace *types.Namespace) (*types.Namespace, error)
		Update(ctx context.Context, namespace *types.Namespace) (*types.Namespace, error)
		DeleteByID(ctx context.Context, namespaceID uint64) error
	}

	namespaceUpdateHandler func(ctx context.Context, ns *types.Namespace) (namespaceChanges, error)
	namespaceChanges       uint8
)

const (
	namespaceUnchanged     namespaceChanges = 0
	namespaceChanged       namespaceChanges = 1
	namespaceLabelsChanged namespaceChanges = 2
)

func Namespace() *namespace {
	return &namespace{
		ac:        DefaultAccessControl,
		eventbus:  eventbus.Service(),
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// search fn() orchestrates pages search, namespace preload and check
func (svc namespace) Find(ctx context.Context, filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	var (
		aProps = &namespaceActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Namespace) (bool, error) {
		if !svc.ac.CanReadNamespace(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchNamespaces(ctx) {
			return NamespaceErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Namespace{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchComposeNamespaces(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledNamespaces(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, NamespaceActionSearch, err)
}

func (svc namespace) FindByID(ctx context.Context, ID uint64) (ns *types.Namespace, err error) {
	return svc.lookup(ctx, func(aProps *namespaceActionProps) (*types.Namespace, error) {
		if ID == 0 {
			return nil, NamespaceErrInvalidID()
		}

		aProps.namespace.ID = ID
		return store.LookupComposeNamespaceByID(ctx, svc.store, ID)
	})
}

// FindByHandle is an alias for FindBySlug
func (svc namespace) FindByHandle(ctx context.Context, handle string) (ns *types.Namespace, err error) {
	return svc.FindBySlug(ctx, handle)
}

func (svc namespace) FindBySlug(ctx context.Context, slug string) (ns *types.Namespace, err error) {
	return svc.lookup(ctx, func(aProps *namespaceActionProps) (*types.Namespace, error) {
		if !handle.IsValid(slug) {
			return nil, NamespaceErrInvalidHandle()
		}

		aProps.namespace.Slug = slug
		return store.LookupComposeNamespaceBySlug(ctx, svc.store, slug)
	})
}

// FindByAny tries to find namespace by id, handle or slug
func (svc namespace) FindByAny(ctx context.Context, identifier interface{}) (r *types.Namespace, err error) {
	if ID, ok := identifier.(uint64); ok {
		r, err = svc.FindByID(ctx, ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			r, err = svc.FindByID(ctx, ID)
		} else {
			r, err = svc.FindByHandle(ctx, strIdentifier)
			if err == nil && r.ID == 0 {
				r, err = svc.FindBySlug(ctx, strIdentifier)
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
func (svc namespace) Create(ctx context.Context, new *types.Namespace) (*types.Namespace, error) {
	var (
		aProps = &namespaceActionProps{changed: new}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Slug) {
			return NamespaceErrInvalidHandle()
		}

		if !svc.ac.CanCreateNamespace(ctx) {
			return NamespaceErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.NamespaceBeforeCreate(new, nil)); err != nil {
			return err
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		if err = store.CreateComposeNamespace(ctx, svc.store, new); err != nil {
			return err
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.NamespaceAfterCreate(new, nil))
		return nil
	})

	return new, svc.recordAction(ctx, aProps, NamespaceActionCreate, err)
}

func (svc namespace) Update(ctx context.Context, upd *types.Namespace) (c *types.Namespace, err error) {
	return svc.updater(ctx, upd.ID, NamespaceActionUpdate, svc.handleUpdate(ctx, upd))
}

func (svc namespace) DeleteByID(ctx context.Context, namespaceID uint64) error {
	return trim1st(svc.updater(ctx, namespaceID, NamespaceActionDelete, svc.handleDelete))
}

func (svc namespace) UndeleteByID(ctx context.Context, namespaceID uint64) error {
	return trim1st(svc.updater(ctx, namespaceID, NamespaceActionUndelete, svc.handleUndelete))
}

func (svc namespace) updater(ctx context.Context, namespaceID uint64, action func(...*namespaceActionProps) *namespaceAction, fn namespaceUpdateHandler) (*types.Namespace, error) {
	var (
		changes namespaceChanges
		ns, old *types.Namespace
		aProps  = &namespaceActionProps{namespace: &types.Namespace{ID: namespaceID}}
		err     error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, err = loadNamespace(ctx, s, namespaceID)
		if err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, ns); err != nil {
			return err
		}

		old = ns.Clone()

		aProps.setNamespace(ns)
		aProps.setChanged(ns)

		if ns.DeletedAt == nil {
			err = svc.eventbus.WaitFor(ctx, event.NamespaceBeforeUpdate(ns, old))
		} else {
			err = svc.eventbus.WaitFor(ctx, event.NamespaceBeforeDelete(ns, old))
		}

		if err != nil {
			return
		}

		if changes, err = fn(ctx, ns); err != nil {
			return err
		}

		if changes&namespaceChanged > 0 {
			if err = store.UpdateComposeNamespace(ctx, svc.store, ns); err != nil {
				return err
			}
		}

		if changes&namespaceLabelsChanged > 0 {
			if err = label.Update(ctx, s, ns); err != nil {
				return
			}
		}

		if ns.DeletedAt == nil {
			err = svc.eventbus.WaitFor(ctx, event.NamespaceAfterUpdate(ns, old))
		} else {
			err = svc.eventbus.WaitFor(ctx, event.NamespaceAfterDelete(nil, old))
		}

		return err
	})

	return ns, svc.recordAction(ctx, aProps, action, err)
}

// lookup fn() orchestrates namespace lookup, and check
func (svc namespace) lookup(ctx context.Context, lookup func(*namespaceActionProps) (*types.Namespace, error)) (ns *types.Namespace, err error) {
	var aProps = &namespaceActionProps{namespace: &types.Namespace{}}

	err = func() error {
		if ns, err = lookup(aProps); errors.IsNotFound(err) {
			return NamespaceErrNotFound()
		} else if err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if !svc.ac.CanReadNamespace(ctx, ns) {
			return NamespaceErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, ns); err != nil {
			return err
		}

		return nil
	}()

	return ns, svc.recordAction(ctx, aProps, NamespaceActionLookup, err)
}

func (svc namespace) uniqueCheck(ctx context.Context, ns *types.Namespace) (err error) {
	if ns.Slug != "" {
		if e, _ := store.LookupComposeNamespaceBySlug(ctx, svc.store, ns.Slug); e != nil && e.ID != ns.ID {
			return NamespaceErrHandleNotUnique()
		}
	}

	return nil
}

func (svc namespace) handleUpdate(ctx context.Context, upd *types.Namespace) namespaceUpdateHandler {
	return func(ctx context.Context, res *types.Namespace) (changes namespaceChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return namespaceUnchanged, NamespaceErrStaleData()
		}

		if upd.Slug != res.Slug && !handle.IsValid(upd.Slug) {
			return namespaceUnchanged, NamespaceErrInvalidHandle()
		}

		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return namespaceUnchanged, err
		}

		if !svc.ac.CanUpdateNamespace(ctx, res) {
			return namespaceUnchanged, NamespaceErrNotAllowedToUpdate()
		}

		if res.Name != upd.Name {
			changes |= namespaceChanged
			res.Name = upd.Name
		}

		if res.Slug != upd.Slug {
			changes |= namespaceChanged
			res.Slug = upd.Slug
		}

		if res.Enabled != upd.Enabled {
			changes |= namespaceChanged
			res.Enabled = upd.Enabled
		}

		if !reflect.DeepEqual(upd.Meta, res.Meta) {
			changes |= namespaceChanged
			res.Meta = upd.Meta
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= namespaceLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if changes&namespaceChanged > 0 {
			res.UpdatedAt = now()
		}

		return
	}
}

func (svc namespace) handleDelete(ctx context.Context, ns *types.Namespace) (namespaceChanges, error) {
	if !svc.ac.CanDeleteNamespace(ctx, ns) {
		return namespaceUnchanged, NamespaceErrNotAllowedToDelete()
	}

	if ns.DeletedAt != nil {
		// namespace already deleted
		return namespaceUnchanged, nil
	}

	ns.DeletedAt = now()
	return namespaceChanged, nil
}

func (svc namespace) handleUndelete(ctx context.Context, ns *types.Namespace) (namespaceChanges, error) {
	if !svc.ac.CanDeleteNamespace(ctx, ns) {
		return namespaceUnchanged, NamespaceErrNotAllowedToUndelete()
	}

	if ns.DeletedAt == nil {
		// namespace not deleted
		return namespaceUnchanged, nil
	}

	ns.DeletedAt = nil
	return namespaceChanged, nil
}

func loadNamespace(ctx context.Context, s store.Storer, namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, ChartErrInvalidNamespaceID()
	}

	if ns, err = store.LookupComposeNamespaceByID(ctx, s, namespaceID); errors.IsNotFound(err) {
		return nil, NamespaceErrNotFound()
	}

	return
}

// toLabeledNamespaces converts to []label.LabeledResource
//
// This function is auto-generated.
func toLabeledNamespaces(set []*types.Namespace) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
