package service

import (
	"archive/zip"
	"context"
	"mime/multipart"
	"reflect"
	"strconv"
	"time"

	automationService "github.com/crusttech/human/server/automation/service"
	"github.com/crusttech/human/server/compose/service/event"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
	"github.com/gabriel-vasile/mimetype"
)

type (
	namespace struct {
		actionlog actionlog.Recorder
		ac        namespaceAccessController
		modAc     moduleAccessController
		pageAc    pageAccessController
		chartAc   chartAccessController

		eventbus eventDispatcher
		store    store.Storer
		locale   ResourceTranslationsManagerService
		envoy    *envoyx.Service
	}

	namespaceImportSession struct {
		Name        string `json:"name"`
		Slug        string `json:"handle"`
		NamespaceID uint64 `json:"namespaceID,string"`
		SessionID   uint64 `json:"sessionID,string"`
		UserID      uint64 `json:"userID,string"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`

		Nodes envoyx.NodeSet `json:"-"`
	}

	namespaceAccessController interface {
		CanManageResourceTranslations(ctx context.Context) bool
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
		Search(context.Context, types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error)
		FindByAny(context.Context, interface{}) (*types.Namespace, error)

		Create(ctx context.Context, namespace *types.Namespace) (*types.Namespace, error)
		Update(ctx context.Context, namespace *types.Namespace) (*types.Namespace, error)
		Clone(ctx context.Context, namespaceID uint64, dup *types.Namespace, decoder func() (envoyx.NodeSet, error)) (ns *types.Namespace, err error)
		ImportInit(ctx context.Context, f multipart.File, size int64) (namespaceImportSession, error)
		ImportRun(ctx context.Context, sessionID uint64, dup *types.Namespace) (ns *types.Namespace, err error)
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

var (
	// @todo this is a temporary implementation; we will rework resource import/export
	//       in the following versions
	namespaceSessionStore = make(map[uint64]namespaceImportSession)
)

func Namespace() *namespace {
	return &namespace{
		ac:      DefaultAccessControl,
		modAc:   DefaultAccessControl,
		pageAc:  DefaultAccessControl,
		chartAc: DefaultAccessControl,

		eventbus:  eventbus.Service(),
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		locale:    DefaultResourceTranslation,
		envoy:     envoyx.Global(),
	}
}

// onSearch is the generated Search body handler.
//
// The recordAction wrapper, aProps (filter) and the standard CanSearchNamespaces
// access check are owned by the generated namespace.gen.go; this handler runs the
// store search, i18n decode and label load. It also installs filter.Check
// (the generated wrapper does not emit it for custom-body search).
func (svc namespace) onSearch(ctx context.Context, filter types.NamespaceFilter, aProps *namespaceActionProps) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Namespace) (bool, error) {
		if !svc.ac.CanReadNamespace(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	if len(filter.Labels) > 0 {
		filter.LabeledIDs, err = label.Search(
			ctx,
			svc.store,
			types.Namespace{}.LabelResourceKind(),
			filter.Labels,
		)

		if err != nil {
			return
		}

		// labels specified but no labeled resources found
		if len(filter.LabeledIDs) == 0 {
			return
		}
	}

	if set, f, err = store.SearchComposeNamespaces(ctx, svc.store, filter); err != nil {
		return
	}

	// i18n
	tag := locale.GetAcceptLanguageFromContext(ctx)
	set.Walk(func(n *types.Namespace) error {
		n.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, n.ResourceTranslation()))
		return nil
	})

	if err = label.Load(ctx, svc.store, toLabeledNamespaces(set)...); err != nil {
		return
	}

	return
}

// onLookup is the generated FindByID body handler.
//
// The recordAction wrapper and aProps (namespace) are owned by the generated
// namespace.gen.go.
func (svc namespace) onLookup(ctx context.Context, ID uint64, aProps *namespaceActionProps) (ns *types.Namespace, err error) {
	return svc.lookup(ctx, aProps, func(aProps *namespaceActionProps) (*types.Namespace, error) {
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
	var aProps = &namespaceActionProps{namespace: &types.Namespace{}}

	ns, err = svc.lookup(ctx, aProps, func(aProps *namespaceActionProps) (*types.Namespace, error) {
		if !handle.IsValid(slug) {
			return nil, NamespaceErrInvalidHandle()
		}

		aProps.namespace.Slug = slug
		return store.LookupComposeNamespaceBySlug(ctx, svc.store, slug)
	})

	return ns, svc.recordAction(ctx, aProps, NamespaceActionLookup, err)
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

// onCreate is the generated Create body handler. It adds namespace and presets
// access rules for role everyone.
//
// The recordAction wrapper, aProps (namespace) and res=new assignment are owned
// by the generated namespace.gen.go. The CanCreateNamespace access check lives
// here (create is in customAccessOps) because it must run inside the tx, after
// the handle validation.
func (svc namespace) onCreate(ctx context.Context, new *types.Namespace) error {
	var (
		aProps = &namespaceActionProps{namespace: new}
	)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		aProps.setChanged(new)

		if err = store.CreateComposeNamespace(ctx, svc.store, new); err != nil {
			return err
		}

		if err = updateTranslations(ctx, svc.ac, svc.locale, new.EncodeTranslations()...); err != nil {
			return
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.NamespaceAfterCreate(new, nil))
		return nil
	})
}

// onUpdate is the generated Update body handler.
func (svc namespace) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Namespace, aProps *namespaceActionProps, _ func() error, _ func() error) error {
	old := res.Clone()

	if upd.Slug != res.Slug && !handle.IsValid(upd.Slug) {
		return NamespaceErrInvalidHandle()
	}

	if err := svc.uniqueCheck(ctx, upd); err != nil {
		return err
	}

	if !svc.ac.CanUpdateNamespace(ctx, res) {
		return NamespaceErrNotAllowedToUpdate()
	}

	if err := svc.eventbus.WaitFor(ctx, event.NamespaceBeforeUpdate(res, old)); err != nil {
		return err
	}

	if res.Name != upd.Name {
		res.Name = upd.Name
	}

	if res.Slug != upd.Slug {
		res.Slug = upd.Slug
	}

	if res.Enabled != upd.Enabled {
		res.Enabled = upd.Enabled
	}

	if !reflect.DeepEqual(upd.Meta, res.Meta) {
		res.Meta = upd.Meta
	}

	if upd.Labels != nil {
		if label.Changed(res.Labels, upd.Labels) {
			res.Labels = upd.Labels
			if err := label.Update(ctx, s, res); err != nil {
				return err
			}
		}
	}

	res.UpdatedAt = now()
	aProps.setChanged(res)

	if err := store.UpdateComposeNamespace(ctx, s, res); err != nil {
		return err
	}

	if err := updateTranslations(ctx, svc.ac, svc.locale, res.EncodeTranslations()...); err != nil {
		return err
	}

	_ = svc.eventbus.WaitFor(ctx, event.NamespaceAfterUpdate(res, old))
	return nil
}

func (svc namespace) Clone(ctx context.Context, namespaceID uint64, dup *types.Namespace, decoder func() (envoyx.NodeSet, error)) (ns *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{namespace: dup}
	)

	err = func() error {
		err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
			// Preparation
			// - target namespace
			targetNs, err := loadNamespace(ctx, s, namespaceID)
			if errors.IsNotFound(err) {
				return NamespaceErrNotFound()
			} else if err != nil {
				return err
			}
			aProps.setNamespace(targetNs)

			// - destination namespace
			if dup.Slug != "" {
				dstNs, err := store.LookupComposeNamespaceBySlug(ctx, s, dup.Slug)
				if err != nil && err != store.ErrNotFound {
					return err
				}
				if dstNs != nil {
					return NamespaceErrHandleNotUnique()
				}
			}

			// Access control
			if err = svc.canExport(ctx, targetNs); err != nil {
				return err
			}

			// get namespace resources
			nn, err := decoder()
			if err != nil {
				return err
			}

			aProps.setNamespace(dup)
			dup, err = svc.envoyRun(ctx, s, nn, targetNs, dup)
			if err != nil {
				return err
			}

			dup, err = store.LookupComposeNamespaceBySlug(ctx, s, dup.Slug)
			if err != nil {
				return err
			}
			tag := locale.GetAcceptLanguageFromContext(ctx)
			dup.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, dup.ResourceTranslation()))

			aProps.setNamespace(dup)

			return nil
		})
		if err != nil {
			return err
		}

		err = svc.reloadServices(ctx, dup)
		return err
	}()

	return dup, svc.recordAction(ctx, aProps, NamespaceActionClone, err)
}

func (svc namespace) ImportInit(ctx context.Context, f multipart.File, size int64) (namespaceImportSession, error) {
	var (
		aProps  = &namespaceActionProps{}
		err     error
		ns      *types.Namespace
		session namespaceImportSession
		nodes   envoyx.NodeSet

		esvc = envoyx.Global()
		nn   envoyx.NodeSet
	)

	err = func() error {
		// access control
		if err := svc.canImport(ctx); err != nil {
			return err
		}

		// archive type check
		mt, err := mimetype.DetectReader(f)
		if err != nil {
			return err
		}
		aProps.setArchiveFormat(mt.Extension())
		if !mt.Is("application/zip") {
			return NamespaceErrUnsupportedImportFormat()
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}

		// un-archive
		archive, err := zip.NewReader(f, size)
		if err != nil {
			return err
		}

		for _, zf := range archive.File {
			if zf.FileInfo().IsDir() {
				continue
			}

			f, err := zf.Open()
			if err != nil {
				return err
			}
			defer f.Close()

			nn, _, err = esvc.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeIO,
				Params: map[string]any{
					"reader": f,
					"mime":   "text/yaml",
				},
			})
			if err != nil {
				return err
			}

			nodes = append(nodes, nn...)
		}

		// store a session for later
		session = namespaceImportSession{
			SessionID: nextID(),
			UserID:    auth.GetIdentityFromContext(ctx).Identity(),

			CreatedAt: *now(),
			Nodes:     nodes,
		}

		// find the ns node
		for _, n := range nodes {
			if n.ResourceType == types.NamespaceResourceType {
				ns = n.Resource.(*types.Namespace)
			}
		}
		if ns == nil {
			return NamespaceErrImportMissingNamespace()
		}

		// session needs to have namespaceID if ns Handle is not provided
		session.NamespaceID = ns.ID
		session.Name = ns.Name
		session.Slug = ns.Slug
		namespaceSessionStore[session.SessionID] = session

		aProps.setNamespace(ns)
		return nil
	}()

	return session, svc.recordAction(ctx, aProps, NamespaceActionImportInit, err)
}

func (svc namespace) ImportRun(ctx context.Context, sessionID uint64, dup *types.Namespace) (ns *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{namespace: dup}
	)

	err = func() error {
		var (
			newNS *types.Namespace
		)

		err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
			// access control
			if err = svc.canImport(ctx); err != nil {
				return err
			}

			if !handle.IsValid(dup.Slug) {
				return NamespaceErrInvalidHandle()
			}

			if dup.Slug != "" {
				// check for duplicate
				dstNs, err := store.LookupComposeNamespaceBySlug(ctx, svc.store, dup.Slug)
				if err != nil && err != store.ErrNotFound {
					return err
				}
				if dstNs != nil {
					return NamespaceErrHandleNotUnique()
				}
			}

			// session
			var (
				session namespaceImportSession
				ok      bool
			)
			if session, ok = namespaceSessionStore[sessionID]; !ok {
				return NamespaceErrImportSessionNotFound()
			}
			defer func() {
				delete(namespaceSessionStore, sessionID)
			}()

			aProps.setNamespace(dup)

			newNS, err = svc.envoyRun(ctx, s, session.Nodes, &types.Namespace{ID: session.NamespaceID, Slug: session.Slug, Name: session.Name}, dup)
			if err != nil {
				return err
			}

			newNS, err = store.LookupComposeNamespaceBySlug(ctx, s, newNS.Slug)
			if err != nil {
				return err
			}

			aProps.setNamespace(newNS)
			return nil
		})
		if err != nil {
			return err
		}

		err = svc.reloadServices(ctx, newNS)
		dup = newNS
		return err
	}()

	return dup, svc.recordAction(ctx, aProps, NamespaceActionImportRun, err)
}

// onDelete is the generated DeleteByID body handler.
func (svc namespace) onDelete(ctx context.Context, s store.Storer, res *types.Namespace, aProps *namespaceActionProps) error {
	if !svc.ac.CanDeleteNamespace(ctx, res) {
		return NamespaceErrNotAllowedToDelete()
	}

	if res.DeletedAt != nil {
		return nil
	}

	old := res.Clone()

	if err := svc.eventbus.WaitFor(ctx, event.NamespaceBeforeDelete(res, old)); err != nil {
		return err
	}

	res.DeletedAt = now()

	if err := store.UpdateComposeNamespace(ctx, s, res); err != nil {
		return err
	}

	if err := updateTranslations(ctx, svc.ac, svc.locale, res.EncodeTranslations()...); err != nil {
		return err
	}

	_ = svc.eventbus.WaitFor(ctx, event.NamespaceAfterDelete(nil, old))
	return nil
}

func (svc namespace) UndeleteByID(ctx context.Context, namespaceID uint64) error {
	var (
		aProps = &namespaceActionProps{}
	)

	_, err := svc.updater(ctx, namespaceID, aProps, NamespaceActionUndelete, svc.handleUndelete)
	return err
}

func (svc namespace) updater(ctx context.Context, namespaceID uint64, aProps *namespaceActionProps, action func(...*namespaceActionProps) *namespaceAction, fn namespaceUpdateHandler) (*types.Namespace, error) {
	var (
		changes namespaceChanges
		ns, old *types.Namespace
		err     error
	)

	if aProps.namespace == nil {
		aProps.namespace = &types.Namespace{ID: namespaceID}
	}

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

		if err = updateTranslations(ctx, svc.ac, svc.locale, ns.EncodeTranslations()...); err != nil {
			return
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

	return ns, svc.recordAction(ctx, aProps, action, err, old, ns)
}

// lookup fn() orchestrates namespace lookup, and check.
//
// The recordAction wrapper (NamespaceActionLookup) is owned by the generated
// FindByID; UndeleteByID/Clone callers that still own their action log pass their
// own aProps. The aProps.namespace is preset to a blank namespace when nil so the
// lookup fn can stamp the ID/slug onto it.
func (svc namespace) lookup(ctx context.Context, aProps *namespaceActionProps, lookup func(*namespaceActionProps) (*types.Namespace, error)) (ns *types.Namespace, err error) {
	if aProps.namespace == nil {
		aProps.namespace = &types.Namespace{}
	}

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

	return ns, err
}

func (svc namespace) uniqueCheck(ctx context.Context, ns *types.Namespace) (err error) {
	if ns.Slug != "" {
		if e, _ := store.LookupComposeNamespaceBySlug(ctx, svc.store, ns.Slug); e != nil && e.ID != ns.ID {
			return NamespaceErrHandleNotUnique()
		}
	}

	return nil
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

func (svc namespace) canExport(ctx context.Context, namespace *types.Namespace) error {
	// Preload all of the relevant stuff for access control
	// - modules
	//   no need to load fields
	mm, _, err := store.SearchComposeModules(ctx, svc.store, types.ModuleFilter{NamespaceID: namespace.ID})
	if err != nil {
		return err
	}
	// - pages
	pp, _, err := store.SearchComposePages(ctx, svc.store, types.PageFilter{NamespaceID: namespace.ID})
	if err != nil {
		return err
	}
	// - charts
	cc, _, err := store.SearchComposeCharts(ctx, svc.store, types.ChartFilter{NamespaceID: namespace.ID})
	if err != nil {
		return err
	}

	// access control
	// - namespace
	if !svc.ac.CanReadNamespace(ctx, namespace) {
		return NamespaceErrNotAllowedToRead()
	}
	// - modules
	for _, m := range mm {
		if !svc.modAc.CanReadModule(ctx, m) {
			return ModuleErrNotAllowedToRead()
		}
	}
	// - pages
	for _, p := range pp {
		if !svc.pageAc.CanReadPage(ctx, p) {
			return PageErrNotAllowedToRead()
		}
	}
	// - charts
	for _, c := range cc {
		if !svc.chartAc.CanReadChart(ctx, c) {
			return ChartErrNotAllowedToRead()
		}
	}
	return nil
}

func (svc namespace) canImport(ctx context.Context) error {

	// If a user is allowed to create a namespace, they are considered to be allowed
	// to create any underlying resource when it comes to importing.
	//
	// This was agreed upon internally and may change in the future.

	if !svc.ac.CanCreateNamespace(ctx) {
		return NamespaceErrNotAllowedToCreate()
	}
	return nil
}

func (svc namespace) envoyRun(ctx context.Context, s store.Storer, nodes envoyx.NodeSet, oldNS, newNS *types.Namespace) (ns *types.Namespace, err error) {
	// Get the NS node
	oldRef := envoyx.Ref{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  envoyx.MakeIdentifiers(oldNS.Slug, oldNS.ID),
		Scope: envoyx.Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(oldNS.Slug, oldNS.ID),
		},
	}
	nsNode := envoyx.NodeForRef(oldRef, nodes...)
	auxNs := nsNode.Resource.(*types.Namespace)

	// Handle renames and references
	auxNs.ID = 0
	auxNs.Name = newNS.Name
	auxNs.Slug = newNS.Slug
	newNS = auxNs
	ns = newNS

	// Change the identifiers and references
	// - identifiers of the NS node
	nsNode.Identifiers = envoyx.MakeIdentifiers(newNS.Slug)
	nsNode.Scope.Identifiers = nsNode.Identifiers

	// - all the child refs
	for _, n := range nodes {
		nr := make(map[string]envoyx.Ref)
		for k, r := range n.References {
			if r.ResourceType == types.NamespaceResourceType {
				r.Identifiers = nsNode.Identifiers
			}
			if r.Scope.ResourceType == types.NamespaceResourceType {
				r.Scope = nsNode.Scope
			}
			nr[k] = r
		}
		n.References = nr
		if n.Scope.ResourceType == nsNode.Scope.ResourceType {
			n.Scope = nsNode.Scope
		}
	}

	// Get expected placeholder refs
	// - roles
	roles, _, err := svc.envoy.Decode(ctx, envoyx.DecodeParams{
		Type: envoyx.DecodeTypeStore,
		Params: map[string]any{
			"storer": s,
			"dal":    dal.Service(),
		},
		Filter: map[string]envoyx.ResourceFilter{
			systemTypes.RoleResourceType: {},
		},
	})
	if err != nil {
		return
	}
	for _, r := range roles {
		r.Placeholder = true
	}
	nodes = append(nodes, roles...)

	// run the import
	gg, err := svc.envoy.Bake(ctx, envoyx.EncodeParams{
		Type: envoyx.EncodeTypeStore,
		Params: map[string]any{
			"storer": s,
			"dal":    dal.Service(),
		},
	}, nil, nodes...)
	if err != nil {
		return
	}

	err = svc.envoy.Encode(ctx, envoyx.EncodeParams{
		Type: envoyx.EncodeTypeStore,
		Params: map[string]any{
			"storer": s,
			"dal":    dal.Service(),
		},
	}, gg)

	return
}

func (svc namespace) reloadServices(ctx context.Context, ns *types.Namespace) (err error) {
	// Adjust name res. tr. since we're changing it
	if err = updateTranslations(ctx, svc.ac, svc.locale, &locale.ResourceTranslation{
		Resource: ns.ResourceTranslation(),
		Key:      types.LocaleKeyNamespaceName.Path,
		Msg:      locale.SanitizeMessage(ns.Name),
	}); err != nil {
		return
	}

	// Reload RBAC rules (in case import brought in something new)
	rbac.Global().Reload(ctx)
	if err = locale.Global().ReloadResourceTranslations(ctx); err != nil {
		return
	}

	{
		// Reload workflow-triggers (in case import brought in something new)
		if err = automationService.DefaultWorkflow.Load(ctx); err != nil {
			// should not be a fatal error
			err = nil
		}

		// Reload workflow-triggers (in case import brought in something new)
	if err = automationService.DefaultNgAutomation.Load(ctx); err != nil {
		// should not be a fatal error
		err = nil
	}
	}

	return
}

