package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"reflect"
	"strconv"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/text/language"
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
		Clone(ctx context.Context, namespaceID uint64, dup *types.Namespace, decoder func() (resource.InterfaceSet, error), encoder func(resource.InterfaceSet) error) (ns *types.Namespace, err error)
		Export(ctx context.Context, namespaceID uint64, archive string, decoder func() (resource.InterfaceSet, error), encoder func(resource.InterfaceSet) (envoy.Streamer, error)) (r io.ReadSeeker, err error)
		Import(ctx context.Context, f multipart.File, size int64, encoder func(resource.InterfaceSet) error) (ns *types.Namespace, err error)
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
		ac:      DefaultAccessControl,
		modAc:   DefaultAccessControl,
		pageAc:  DefaultAccessControl,
		chartAc: DefaultAccessControl,

		eventbus:  eventbus.Service(),
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		locale:    DefaultResourceTranslation,
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

		// i18n
		tag := locale.GetAcceptLanguageFromContext(ctx)
		set.Walk(func(n *types.Namespace) error {
			n.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, n.ResourceTranslation()))
			return nil
		})

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

		if contentLang := locale.GetContentLanguageFromContext(ctx); contentLang != language.Und {
			tt := new.EncodeTranslations()
			tt.SetLanguage(contentLang)
			err = DefaultResourceTranslation.Upsert(ctx, tt)
			if err != nil {
				return err
			}
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

func (svc namespace) Clone(ctx context.Context, namespaceID uint64, dup *types.Namespace, decoder func() (resource.InterfaceSet, error), encoder func(resource.InterfaceSet) error) (ns *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{namespace: dup}
	)

	err = func() error {
		// Preparation
		// - target namespace
		targetNs, err := loadNamespace(ctx, svc.store, namespaceID)
		if errors.IsNotFound(err) {
			return NamespaceErrNotFound()
		} else if err != nil {
			return err
		}
		aProps.setNamespace(targetNs)

		// - destination namespace
		dstNs, err := store.LookupComposeNamespaceBySlug(ctx, svc.store, dup.Slug)
		if err != nil && err != store.ErrNotFound {
			return err
		}
		if dstNs != nil {
			return NamespaceErrHandleNotUnique()
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

		// some meta bits
		sNsID := strconv.FormatUint(namespaceID, 10)
		oldNsRef := resource.MakeRef(types.NamespaceResourceType, resource.MakeIdentifiers(sNsID))
		newNsRef := resource.MakeRef(types.NamespaceResourceType, resource.MakeIdentifiers(dup.Slug, dup.Name))
		prune := resource.RefSet{resource.MakeWildRef(automationTypes.WorkflowResourceType)}

		// rename the namespace
		//
		// For now we will find the namespace in set and change it's name, handle.
		// The rest of the resources can stay as are.
		//
		// @todo add a more flexible system for such modifications
		auxNs := resource.FindComposeNamespace(nn, oldNsRef.Identifiers)
		auxNs.ID = 0
		auxNs.Name = dup.Name
		auxNs.Slug = dup.Slug
		dup = auxNs
		aProps.setNamespace(dup)

		// Correct internal references
		// - namespace identifiers
		nn.SearchForIdentifiers(oldNsRef.Identifiers).Walk(func(r resource.Interface) error {
			r.ReID(newNsRef.Identifiers)
			return nil
		})

		// - relations
		nn.SearchForReferences(oldNsRef).Walk(func(r resource.Interface) error {
			r.ReRef(resource.RefSet{oldNsRef}, resource.RefSet{newNsRef})

			// - additional pruning
			pp, ok := r.(resource.PrunableInterface)
			if !ok {
				return nil
			}

			for _, p := range prune {
				pp.Prune(p)
			}
			return nil
		})

		// encode
		return encoder(nn)
	}()

	return dup, svc.recordAction(ctx, aProps, NamespaceActionClone, err)
}

func (svc namespace) Export(ctx context.Context, namespaceID uint64, archive string, decoder func() (resource.InterfaceSet, error), encoder func(resource.InterfaceSet) (envoy.Streamer, error)) (r io.ReadSeeker, err error) {
	var (
		aProps = &namespaceActionProps{archiveFormat: archive}
	)

	// make archive
	buf := bytes.NewBuffer(nil)
	w := zip.NewWriter(buf)

	err = func() error {
		if archive != "zip" {
			return NamespaceErrUnsupportedExportFormat()
		}

		// initial validation
		// - target namespace
		targetNs, err := store.LookupComposeNamespaceByID(ctx, svc.store, namespaceID)
		if err != nil && err != store.ErrNotFound {
			return err
		}
		aProps.setNamespace(targetNs)

		// - ac
		if err = svc.canExport(ctx, targetNs); err != nil {
			return err
		}

		// get namespace resources
		nn, err := decoder()
		if err != nil {
			return err
		}

		// some meta bits
		sNsID := strconv.FormatUint(namespaceID, 10)
		oldNsRef := resource.MakeRef(types.NamespaceResourceType, resource.MakeIdentifiers(sNsID))
		prune := resource.RefSet{resource.MakeWildRef(automationTypes.WorkflowResourceType)}

		// - prune resources we won't preserve
		nn.SearchForReferences(oldNsRef).Walk(func(r resource.Interface) error {
			pp, ok := r.(resource.PrunableInterface)
			if !ok {
				return nil
			}

			for _, p := range prune {
				pp.Prune(p)
			}
			return nil
		})

		// encode
		ss, err := encoder(nn)
		if err != nil {
			return err
		}

		// create archive
		for _, s := range ss.Stream() {
			// @todo generalize when needed
			f, err := w.Create(fmt.Sprintf("%s.yaml", s.Resource))
			if err != nil {
				return err
			}

			bb, err := ioutil.ReadAll(s.Source)
			if err != nil {
				return err
			}

			_, err = f.Write(bb)
			if err != nil {
				return err
			}
		}

		return w.Close()
	}()

	return bytes.NewReader(buf.Bytes()), svc.recordAction(ctx, aProps, NamespaceActionExport, err)
}

func (svc namespace) Import(ctx context.Context, f multipart.File, size int64, encoder func(resource.InterfaceSet) error) (ns *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{}
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

		// decode with Envoy
		yd := yaml.Decoder()
		nn := make([]resource.Interface, 0, 10)

		for _, f := range archive.File {
			a, err := f.Open()
			if err != nil {
				return err
			}
			defer a.Close()

			mm, err := yd.Decode(ctx, a, nil)
			if err != nil {
				return err
			}
			nn = append(nn, mm...)
		}

		// encode
		err = encoder(nn)
		if err != nil {
			return err
		}

		// find the ns node
		for _, n := range nn {
			if nsn, ok := n.(*resource.ComposeNamespace); ok {
				ns = nsn.Res
				break
			}
		}

		aProps.setNamespace(ns)

		return nil
	}()

	return ns, svc.recordAction(ctx, aProps, NamespaceActionImport, err)
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

		if contentLang := locale.GetContentLanguageFromContext(ctx); contentLang != language.Und {
			tt := ns.EncodeTranslations()
			tt.SetLanguage(contentLang)
			err = DefaultResourceTranslation.Upsert(ctx, tt)
			if err != nil {
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
