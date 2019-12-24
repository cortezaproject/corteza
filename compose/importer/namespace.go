package importer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	Namespace struct {
		imp *Importer

		set   types.NamespaceSet
		dirty map[uint64]bool

		// modules per namespace
		modules map[string]*Module

		// charts per namespace
		charts map[string]*Chart

		// pages per namespace
		pages map[string]*Page

		// records per namespace
		records map[string]*Record
	}

	// @todo remove finder strategy, directly provide set of items
	namespaceFinder interface {
		Find(filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error)
	}

	namespaceKeeper interface {
		Update(*types.Namespace) (*types.Namespace, error)
		Create(*types.Namespace) (*types.Namespace, error)
	}
)

func NewNamespaceImporter(imp *Importer) *Namespace {
	out := &Namespace{
		imp: imp,

		set:   types.NamespaceSet{},
		dirty: make(map[uint64]bool),

		modules: map[string]*Module{},
		charts:  map[string]*Chart{},
		pages:   map[string]*Page{},
		records: map[string]*Record{},
	}

	if imp.namespaceFinder != nil {
		out.set, _, _ = imp.namespaceFinder.Find(types.NamespaceFilter{})
	}

	return out
}

// CastSet resolves permission rules:
// { <namespace-handle>: { namespace } } or [ { namespace }, ... ]
func (nsImp *Namespace) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Namespaces defined as collection
			deinterfacer.KVsetString(&handle, "slug", def)
			deinterfacer.KVsetString(&handle, "handle", def, handle)
		}

		return nsImp.Cast(handle, def)
	})
}

// Cast resolves permission rules:
// { <namespace-handle>: { namespace } } or [ { namespace }, ... ]
func (nsImp *Namespace) Cast(handle string, def interface{}) (err error) {
	if !deinterfacer.IsMap(def) {
		return errors.New("expecting map of values for namespace")
	}

	var namespace *types.Namespace

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid namespace handle")
	}

	handle = importer.NormalizeHandle(handle)

	if namespace, err = nsImp.Get(handle); err != nil {
		return err
	} else if namespace == nil {
		namespace = &types.Namespace{
			Slug:    handle,
			Name:    handle,
			Enabled: true,
		}
	} else if namespace.ID == 0 {
		// We will ignore that namespace has already been defined because
		// we want to support multiple calls to Cast() fn, ie when namespace
		// config is split into multiple files
	} else {
		nsImp.dirty[namespace.ID] = true
	}

	nsImp.Setup(namespace)

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "handle", "slug":
		// already handled

		case "name", "title", "label":
			namespace.Name = deinterfacer.ToString(val)

		case "enabled":
			namespace.Enabled = deinterfacer.ToBool(val)

		case "meta":
			namespace.Meta, err = nsImp.castMeta(namespace, val)

			return

		case "modules":
			return nsImp.castModules(handle, val)

		case "charts":
			return nsImp.castCharts(handle, val)

		case "pages":
			return nsImp.castPages(handle, val)

		case "records":
			return nsImp.castRecords(handle, val)

		case "allow", "deny":
			return nsImp.imp.permissions.CastSet(types.NamespacePermissionResource.String()+namespace.Slug, key, val)

		default:
			return fmt.Errorf("unexpected key %q for namespace %q", key, namespace.Slug)
		}

		return err
	})
}

func (nsImp *Namespace) castMeta(ns *types.Namespace, def interface{}) (types.NamespaceMeta, error) {
	var meta = types.NamespaceMeta{}

	return meta, deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "subtitle":
			meta.Subtitle = deinterfacer.ToString(val)

		case "description":
			meta.Description = deinterfacer.ToString(val)

		default:
			return fmt.Errorf("unexpected key %q for namespace %q meta", key, ns.Slug)

		}
		return
	})
}

func (nsImp *Namespace) castModules(handle string, def interface{}) error {
	if nsImp.modules[handle] == nil {
		return fmt.Errorf("unknown namespace %q", handle)

	}

	return nsImp.modules[handle].CastSet(def)
}

func (nsImp *Namespace) castCharts(handle string, def interface{}) error {
	if nsImp.charts[handle] == nil {
		return fmt.Errorf("unknown namespace %q", handle)

	}

	return nsImp.charts[handle].CastSet(def)
}

func (nsImp *Namespace) castPages(handle string, def interface{}) error {
	if nsImp.pages[handle] == nil {
		return fmt.Errorf("unknown namespace %q", handle)

	}

	return nsImp.pages[handle].CastSet(def)
}

func (nsImp *Namespace) castRecords(handle string, def interface{}) error {
	if nsImp.records[handle] == nil {
		return fmt.Errorf("unknown namespace %q", handle)

	}

	return nsImp.records[handle].CastSet(def)
}

// Get finds or creates a new namespace
func (nsImp *Namespace) Get(handle string) (*types.Namespace, error) {
	handle = importer.NormalizeHandle(handle)

	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid namespace handle")
	}

	return nsImp.set.FindByHandle(handle), nil
}

func (nsImp *Namespace) Setup(namespace *types.Namespace) {
	if nsImp.set.FindByHandle(namespace.Slug) == nil {
		nsImp.set = append(nsImp.set, namespace)
	}

	if _, has := nsImp.modules[namespace.Slug]; !has {
		nsImp.modules[namespace.Slug] = NewModuleImporter(nsImp.imp, namespace)
		nsImp.pages[namespace.Slug] = NewPageImporter(nsImp.imp, namespace)
		nsImp.charts[namespace.Slug] = NewChartImporter(nsImp.imp, namespace)
		nsImp.records[namespace.Slug] = NewRecordImporter(nsImp.imp, namespace)
	}
}

func (nsImp *Namespace) Store(ctx context.Context, nsk namespaceKeeper, mk moduleKeeper, ck chartKeeper, pk pageKeeper, rk recordKeeper) error {
	return nsImp.set.Walk(func(namespace *types.Namespace) (err error) {
		var handle = namespace.Slug

		if namespace.ID == 0 {
			namespace, err = nsk.Create(namespace)
		} else if nsImp.dirty[namespace.ID] {
			namespace, err = nsk.Update(namespace)
		}

		if err != nil {
			return
		}

		nsImp.dirty[namespace.ID] = false
		nsImp.imp.permissions.UpdateResources(types.NamespacePermissionResource.String(), handle, namespace.ID)

		if _, ok := nsImp.modules[handle]; ok {
			nsImp.modules[handle].namespace = namespace
			if err = nsImp.modules[handle].Store(ctx, mk); err != nil {
				return errors.Wrap(err, "could not import modules")
			}

			nsImp.charts[handle].namespace = namespace
			if err = nsImp.charts[handle].Store(ctx, ck); err != nil {
				return errors.Wrap(err, "could not import charts")
			}

			nsImp.pages[handle].namespace = namespace
			if err = nsImp.pages[handle].Store(ctx, pk); err != nil {
				return errors.Wrap(err, "could not import pages")
			}

			nsImp.records[handle].namespace = namespace
			if err = nsImp.records[handle].Store(ctx, rk); err != nil {
				return errors.Wrap(err, "could not import records")
			}
		}

		return
	})
}
