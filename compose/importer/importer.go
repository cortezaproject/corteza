package importer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	Importer struct {
		namespaces *Namespace

		namespaceFinder  namespaceFinder
		moduleFinder     moduleFinder
		chartFinder      chartFinder
		pageFinder       pageFinder
		automationFinder automationFinder

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

	automationScriptKeeper interface {
		UpdateScript(context.Context, *automation.Script) error
		CreateScript(context.Context, *automation.Script) error
	}

	roleFinder interface {
		Find(context.Context) (sysTypes.RoleSet, error)
	}
)

func NewImporter(nsf namespaceFinder, mf moduleFinder, cf chartFinder, pf pageFinder, af automationFinder, p importer.PermissionImporter) *Importer {
	imp := &Importer{
		namespaceFinder:  nsf,
		moduleFinder:     mf,
		chartFinder:      cf,
		pageFinder:       pf,
		automationFinder: af,

		permissions: p,
	}

	imp.namespaces = NewNamespaceImporter(imp)
	return imp
}

func (imp *Importer) GetNamespaceImporter() *Namespace {
	return imp.namespaces
}

func (imp *Importer) GetModuleImporter(handle string) *Module {
	return imp.namespaces.modules[handle]
}

func (imp *Importer) GetPageImporter(handle string) *Page {
	return imp.namespaces.pages[handle]
}

func (imp *Importer) GetChartImporter(handle string) *Chart {
	return imp.namespaces.charts[handle]
}

func (imp *Importer) Cast(in interface{}) (err error) {
	return deinterfacer.Each(in, func(index int, key string, val interface{}) (err error) {
		switch key {
		case "namespaces":
			return imp.namespaces.CastSet(val)

		case "namespace":
			return imp.namespaces.CastSet([]interface{}{val})

		case "allow", "deny":
			return imp.permissions.CastResourcesSet(key, val)

		default:
			err = fmt.Errorf("unexpected key %q", key)
		}

		return err
	})
}

func (imp *Importer) Store(
	ctx context.Context,
	nsStore namespaceKeeper,
	mStore moduleKeeper,
	cStore chartKeeper,
	pStore pageKeeper,
	asStore automationScriptKeeper,
	pk permissions.ImportKeeper,
	roles sysTypes.RoleSet,
) (err error) {
	err = imp.namespaces.Store(ctx, nsStore, mStore, cStore, pStore, asStore)
	if err != nil {
		return errors.Wrap(err, "could not import namespaces")
	}

	// Make sure we properly replace role handles with IDs
	_ = roles.Walk(func(role *sysTypes.Role) error {
		imp.permissions.UpdateRoles(role.Handle, role.ID)
		return nil
	})

	err = imp.permissions.Store(ctx, pk)
	if err != nil {
		return errors.Wrap(err, "could not import permissions")
	}

	return nil
}
