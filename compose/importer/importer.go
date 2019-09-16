package importer

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	Importer struct {
		baseNamespace string
		namespaces    *NamespaceImport

		namespaceFinder namespaceFinder
		moduleFinder    moduleFinder
		chartFinder     chartFinder
		pageFinder      pageFinder

		permissions importer.PermissionImporter
	}

	moduleKeeper interface {
		Update(*types.Module) (*types.Module, error)
		Create(*types.Module) (*types.Module, error)
	}

	chartKeeper interface {
		Update(*types.Chart) (*types.Chart, error)
		Create(*types.Chart) (*types.Chart, error)
	}

	pageKeeper interface {
		Update(*types.Page) (*types.Page, error)
		Create(*types.Page) (*types.Page, error)
	}
)

func NewImporter(nsf namespaceFinder, mf moduleFinder, cf chartFinder, pf pageFinder, p importer.PermissionImporter) *Importer {
	return &Importer{
		namespaceFinder: nsf,
		moduleFinder:    mf,
		chartFinder:     cf,
		pageFinder:      pf,

		permissions: p,
	}
}

func (imp *Importer) YAML(r io.Reader) (err error) {
	var aux interface{}

	if err = yaml.NewDecoder(r).Decode(&aux); err != nil {
		return
	}

	return imp.Cast(aux)
}

func (imp *Importer) Cast(in interface{}) (err error) {
	if imp.namespaces == nil {
		imp.namespaces = NewNamespaceImporter(
			imp.namespaceFinder,
			imp.moduleFinder,
			imp.chartFinder,
			imp.pageFinder,
			imp.permissions,
		)
	}

	return deinterfacer.Each(in, func(index int, key string, val interface{}) (err error) {
		switch key {
		case "namespaces":
			return imp.namespaces.CastSet(val)
		case "namespace":
			return imp.namespaces.CastSet([]interface{}{val})

		case "modules":
			return imp.namespaces.castModules(imp.baseNamespace, val)
		case "module":
			return imp.namespaces.castModules(imp.baseNamespace, []interface{}{val})

		case "charts":
			return imp.namespaces.castCharts(imp.baseNamespace, val)
		case "chart":
			return imp.namespaces.castCharts(imp.baseNamespace, []interface{}{val})

		case "pages":
			return imp.namespaces.castPages(imp.baseNamespace, val)
		case "page":
			return imp.namespaces.castPages(imp.baseNamespace, []interface{}{val})

		case "allow", "deny":
			return imp.permissions.CastResourcesSet(key, val)

		default:
			err = fmt.Errorf("unexpected key %q", key)
		}

		return err
	})
}

func (imp *Importer) Store(ctx context.Context, nsStore namespaceKeeper, mStore moduleKeeper, cStore chartKeeper, pStore pageKeeper, pk permissions.ImportKeeper) (err error) {
	err = imp.namespaces.Store(ctx, nsStore, mStore, cStore, pStore)
	if err != nil {
		return errors.Wrap(err, "could not import namespaces")
	}

	err = imp.permissions.Store(ctx, pk)
	if err != nil {
		return errors.Wrap(err, "could not import permissions")
	}

	return nil
}
