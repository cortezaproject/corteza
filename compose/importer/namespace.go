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

		set types.NamespaceSet

		// modules per namespace
		modules map[string]*Module

		// charts per namespace
		charts map[string]*Chart

		// pages per namespace
		pages map[string]*Page
	}

	namespaceFinder interface {
		FindByHandle(string) (*types.Namespace, error)
	}

	namespaceKeeper interface {
		Update(*types.Namespace) (*types.Namespace, error)
		Create(*types.Namespace) (*types.Namespace, error)
	}
)

func NewNamespaceImporter(imp *Importer) *Namespace {
	return &Namespace{
		imp: imp,

		set: types.NamespaceSet{},

		modules: map[string]*Module{},
		charts:  map[string]*Chart{},
		pages:   map[string]*Page{},
	}
}

// CastSet resolves permission rules:
// { <namespace-handle>: { namespace } } or [ { namespace }, ... ]
func (nsImp *Namespace) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Namespaces defined as collection
			deinterfacer.KVsetString(&handle, "slug", def)
			deinterfacer.KVsetString(&handle, "handle", handle)
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
	}

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

		case "allow", "deny":
			return nsImp.imp.permissions.CastSet(types.NamespacePermissionResource.String()+namespace.Slug, key, val)

		default:
			return fmt.Errorf("unexpected key %q for namespace %q", key, namespace.Slug)
		}

		return err
	})
}

func (cImp *Namespace) castMeta(ns *types.Namespace, def interface{}) (types.NamespaceMeta, error) {
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

func (nsImp *Namespace) Exists(handle string) bool {
	handle = importer.NormalizeHandle(handle)

	var (
		namespace *types.Namespace
		err       error
	)

	namespace = nsImp.set.FindByHandle(handle)
	if namespace != nil {
		return true
	}

	if nsImp.imp.namespaceFinder != nil {
		namespace, err = nsImp.imp.namespaceFinder.FindByHandle(handle)
		if err == nil && namespace != nil {
			nsImp.set = append(nsImp.set, namespace)
			return true
		}
	}

	return false
}

// Get finds or creates a new namespace
func (nsImp *Namespace) Get(handle string) (*types.Namespace, error) {
	handle = importer.NormalizeHandle(handle)

	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid namespace handle")
	}

	if !nsImp.Exists(handle) {
		nsImp.set = append(nsImp.set, &types.Namespace{
			Slug:    handle,
			Name:    handle,
			Enabled: true,
		})
	}

	namespace := nsImp.set.FindByHandle(handle)

	nsImp.pages[handle] = NewPageImporter(nsImp.imp, namespace)
	nsImp.modules[handle] = NewModuleImporter(nsImp.imp, namespace)
	nsImp.charts[handle] = NewChartImporter(nsImp.imp, namespace)

	return namespace, nil
}

func (nsImp *Namespace) Store(ctx context.Context, nsk namespaceKeeper, mk moduleKeeper, ck chartKeeper, pk pageKeeper) error {
	return nsImp.set.Walk(func(namespace *types.Namespace) (err error) {
		var handle = namespace.Slug

		if namespace.ID == 0 {
			namespace, err = nsk.Create(namespace)
		} else {
			namespace, err = nsk.Update(namespace)
		}

		if err != nil {
			return
		}

		nsImp.imp.permissions.UpdateResources(types.NamespacePermissionResource.String(), handle, namespace.ID)

		if _, ok := nsImp.modules[handle]; ok {
			nsImp.modules[handle].namespace = namespace
			if err = nsImp.modules[handle].Store(ctx, mk); err != nil {
				return errors.Wrap(err, "could not import modules")
			}

		}

		if err = nsImp.charts[handle].Store(ctx, ck); err != nil {
			return errors.Wrap(err, "could not import charts")
		}

		if err = nsImp.pages[handle].Store(ctx, pk); err != nil {
			return errors.Wrap(err, "could not import pages")
		}

		return
	})
}
