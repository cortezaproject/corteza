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
	NamespaceImport struct {
		set types.NamespaceSet

		modules map[string]*ModuleImport
		charts  map[string]*ChartImport
		pages   map[string]*PageImport

		permissions importer.PermissionImporter

		namespaceFinder namespaceFinder
		moduleFinder    moduleFinder
		chartFinder     chartFinder
		pageFinder      pageFinder
	}

	namespaceFinder interface {
		FindByHandle(string) (*types.Namespace, error)
	}

	namespaceKeeper interface {
		Update(*types.Namespace) (*types.Namespace, error)
		Create(*types.Namespace) (*types.Namespace, error)
	}
)

func NewNamespaceImporter(nsf namespaceFinder, mf moduleFinder, cf chartFinder, pf pageFinder, p importer.PermissionImporter) *NamespaceImport {
	return &NamespaceImport{
		set: types.NamespaceSet{},

		modules: map[string]*ModuleImport{},
		charts:  map[string]*ChartImport{},
		pages:   map[string]*PageImport{},

		namespaceFinder: nsf,
		moduleFinder:    mf,
		chartFinder:     cf,
		pageFinder:      pf,

		permissions: p,
	}
}

// CastSet resolves permission rules:
// { <namespace-handle>: { namespace } } or [ { namespace }, ... ]
func (imp *NamespaceImport) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Namespaces defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return imp.Cast(handle, def)
	})
}

// Cast resolves permission rules:
// { <namespace-handle>: { namespace } } or [ { namespace }, ... ]
func (imp *NamespaceImport) Cast(handle string, def interface{}) (err error) {
	if !deinterfacer.IsMap(def) {
		return errors.New("expecting map of values for namespace")
	}

	var namespace *types.Namespace

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid namespace handle")
	}

	handle = importer.NormalizeHandle(handle)
	if namespace, err = imp.Get(handle); err != nil {
		return err
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "handle":
			// already handled
		case "name":
			namespace.Name = deinterfacer.ToString(val)
		case "enabled":
			namespace.Enabled = deinterfacer.ToBool(val)
		case "meta":
			// @todo Namespace.Meta

		case "modules":
			return imp.castModules(handle, val)

		case "charts":
			return imp.castCharts(handle, val)

		case "pages":
			return imp.castPages(handle, val)

		case "allow", "deny":
			return imp.permissions.CastSet(types.NamespacePermissionResource.String()+namespace.Slug, key, val)

		default:
			return fmt.Errorf("unexpected key %q for namespace %q", key, namespace.Slug)
		}

		return err
	})
}

func (imp *NamespaceImport) castModules(namespace string, def interface{}) error {
	return imp.modules[namespace].CastSet(def)
}

func (imp *NamespaceImport) castCharts(namespace string, def interface{}) error {
	return imp.charts[namespace].CastSet(def)
}

func (imp *NamespaceImport) castPages(namespace string, def interface{}) error {
	return imp.pages[namespace].CastSet(def)
}

func (imp *NamespaceImport) Exists(handle string) bool {
	handle = importer.NormalizeHandle(handle)

	var (
		namespace *types.Namespace
		err       error
	)

	namespace = imp.set.FindByHandle(handle)
	if namespace != nil {
		return true
	}

	if imp.namespaceFinder != nil {
		namespace, err = imp.namespaceFinder.FindByHandle(handle)
		if err == nil && namespace != nil {
			imp.set = append(imp.set, namespace)
			return true
		}
	}

	return false
}

// Get finds or creates a new namespace
func (imp *NamespaceImport) Get(handle string) (*types.Namespace, error) {
	handle = importer.NormalizeHandle(handle)

	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid namespace handle")
	}

	if !imp.Exists(handle) {
		imp.set = append(imp.set, &types.Namespace{
			Slug: handle,
			Name: handle,
		})
	}

	ns := imp.set.FindByHandle(handle)

	imp.pages[handle] = NewPageImporter(ns, imp.pageFinder, imp.permissions)
	imp.modules[handle] = NewModuleImporter(ns, imp.moduleFinder, imp.pages[handle], imp.permissions)
	imp.charts[handle] = NewChartImporter(ns, imp.chartFinder, imp.permissions)

	return ns, nil
}

func (imp *NamespaceImport) Store(ctx context.Context, nsk namespaceKeeper, mk moduleKeeper, ck chartKeeper, pk pageKeeper) error {
	return imp.set.Walk(func(namespace *types.Namespace) (err error) {
		var handle = namespace.Slug

		if namespace.ID == 0 {
			namespace, err = nsk.Create(namespace)
		} else {
			namespace, err = nsk.Update(namespace)
		}

		if err != nil {
			return
		}

		imp.permissions.UpdateResources(types.NamespacePermissionResource.String(), handle, namespace.ID)

		if _, ok := imp.modules[handle]; ok {
			imp.modules[handle].namespace = namespace
			if err = imp.modules[handle].Store(ctx, mk); err != nil {
				return errors.Wrap(err, "could not import modules")
			}

		}

		if err = imp.charts[handle].Store(ctx, ck); err != nil {
			return errors.Wrap(err, "could not import charts")
		}

		if err = imp.pages[handle].Store(ctx, pk); err != nil {
			return errors.Wrap(err, "could not import pages")
		}

		return
	})
}
