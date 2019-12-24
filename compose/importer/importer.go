package importer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	Importer struct {
		namespaces *Namespace

		namespaceFinder namespaceFinder
		moduleFinder    moduleFinder
		chartFinder     chartFinder
		pageFinder      pageFinder

		permissions importer.PermissionImporter
		settings    importer.SettingImporter
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

	recordKeeper interface {
		Update(*types.Record) (*types.Record, error)
		Create(*types.Record) (*types.Record, error)
	}

	automationScriptKeeper interface {
		UpdateScript(context.Context, *automation.Script) error
		CreateScript(context.Context, *automation.Script) error
	}
)

func NewImporter(nsf namespaceFinder, mf moduleFinder, cf chartFinder, pf pageFinder, p importer.PermissionImporter, s importer.SettingImporter) *Importer {
	imp := &Importer{
		namespaceFinder: nsf,
		moduleFinder:    mf,
		chartFinder:     cf,
		pageFinder:      pf,

		permissions: p,
		settings:    s,
	}

	if nsf != nil {
		imp.namespaces = NewNamespaceImporter(imp)
	}
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

func (imp *Importer) GetRecordImporter(handle string) *Record {
	return imp.namespaces.records[handle]
}

func (imp *Importer) GetChartImporter(handle string) *Chart {
	return imp.namespaces.charts[handle]
}

func (imp *Importer) Cast(def interface{}) (err error) {
	var nsHandle string
	// Solving a special case where namespace is defined as string
	// and we're treating value as namespace's handle
	deinterfacer.KVsetString(&nsHandle, "namespace", def)
	if nsHandle != "" {
		delete(def.(map[interface{}]interface{}), "namespace")
		if imp.namespaces != nil {
			return imp.namespaces.Cast(nsHandle, def)
		} else {
			return nil
		}
	}

	return deinterfacer.Each(def, func(index int, key string, val interface{}) (err error) {
		switch key {
		case "namespaces":
			if imp.namespaces != nil {
				return imp.namespaces.CastSet(val)
			}

		case "namespace":
			if imp.namespaces != nil {
				return imp.namespaces.CastSet([]interface{}{val})
			}

		case "settings":
			if imp.settings != nil {
				return imp.settings.CastSet(val)
			}

		case "allow", "deny":
			if imp.permissions != nil {
				return imp.permissions.CastResourcesSet(key, val)
			}

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
	rStore recordKeeper,
	pk permissions.ImportKeeper,
	sk settings.ImportKeeper,
	roles sysTypes.RoleSet,
) (err error) {
	if imp.namespaces != nil {
		err = imp.namespaces.Store(ctx, nsStore, mStore, cStore, pStore, rStore)
		if err != nil {
			return errors.Wrap(err, "could not import namespaces")
		}

	}

	// Make sure we properly replace role handles with IDs
	if imp.permissions != nil {
		_ = roles.Walk(func(role *sysTypes.Role) error {
			imp.permissions.UpdateRoles(role.Handle, role.ID)
			return nil
		})

		err = imp.permissions.Store(ctx, pk)
		if err != nil {
			return errors.Wrap(err, "could not import permissions")
		}
	}

	if imp.settings != nil {
		err = imp.settings.Store(ctx, sk)
		if err != nil {
			return errors.Wrap(err, "could not import settings")
		}
	}

	return nil
}
