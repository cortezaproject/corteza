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
	PageImport struct {
		namespace *types.Namespace
		set       types.PageSet

		// page => module maps (module/record-pages)
		modules map[string]string

		// child => parent handle
		parents map[string]string

		pages importer.Interface

		permissions importer.PermissionImporter

		finder pageFinder
	}

	pageFinder interface {
		FindByHandle(uint64, string) (*types.Page, error)
	}
)

func NewPageImporter(ns *types.Namespace, finder pageFinder, permissions importer.PermissionImporter) *PageImport {
	return &PageImport{
		namespace: ns,

		set: types.PageSet{},

		modules: map[string]string{},
		parents: map[string]string{},

		permissions: permissions,

		finder: finder,
	}
}

// CastSet Resolves permission rules:
// { <page-handle>: { page } } or [ { page }, ... ]
func (imp *PageImport) CastSet(set interface{}) error {
	return imp.castSet("", set)
}

// CastSet Resolves permission rules:
// { <page-handle>: { page } } or [ { page }, ... ]
func (imp *PageImport) castSet(parent string, set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Pages defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return imp.cast(parent, handle, def)
	})
}

func (imp *PageImport) Cast(handle string, def interface{}) (err error) {
	return imp.cast("", handle, def)
}

// Cast Resolves permission rules:
// { <page-handle>: { page } } or [ { page }, ... ]
func (imp *PageImport) cast(parent, handle string, def interface{}) (err error) {
	var page *types.Page

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid page handle")
	}

	handle = importer.NormalizeHandle(handle)
	if page, err = imp.Get(handle); err != nil {
		return err
	}

	if parent != "" {
		imp.parents[handle] = parent
	}

	if title, ok := def.(string); ok && title != "" {
		page.Title = title
		return nil
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "namespace":
			// namespace value sanity check
			if deinterfacer.ToString(val, imp.namespace.Slug) != imp.namespace.Slug {
				return fmt.Errorf("explicitly set namespace on page %q shadows inherited namespace", imp.namespace.Slug)
			}

		case "handle":
			// handle value sanity check
			if deinterfacer.ToString(val, handle) != handle {
				return fmt.Errorf("explicitly set handle on page %q shadows inherited handle", handle)
			}

		case "module":
			imp.modules[handle] = deinterfacer.ToString(val)

		case "visible":
			page.Visible = deinterfacer.ToBool(val)

		case "title":
			page.Title = deinterfacer.ToString(val)

		case "description":
			page.Description = deinterfacer.ToString(val)

		case "blocks":
			// @todo Page.Blocks

		case "pages":
			return imp.castSet(handle, val)

		case "allow", "deny":
			return imp.permissions.CastSet(types.PagePermissionResource.String()+handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for page %q", key, handle)
		}

		return err
	})
}

func (imp *PageImport) Exists(handle string) bool {
	handle = importer.NormalizeHandle(handle)

	var (
		page *types.Page
		err  error
	)

	page = imp.set.FindByHandle(handle)
	if page != nil {
		return true
	}

	if imp.namespace.ID == 0 {
		// Assuming new namespace, nothing exists yet..
		return false
	}

	if imp.finder != nil {
		page, err = imp.finder.FindByHandle(imp.namespace.ID, handle)
		if err == nil && page != nil {
			imp.set = append(imp.set, page)
			return true
		}
	}

	return false
}

// Get finds or makes a new page
func (imp *PageImport) Get(handle string) (*types.Page, error) {
	handle = importer.NormalizeHandle(handle)

	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid page handle")
	}

	if !imp.Exists(handle) {
		imp.set = append(imp.set, &types.Page{
			Handle: handle,
			Title:  handle,
		})
	}

	return imp.set.FindByHandle(handle), nil
}

func (imp *PageImport) Store(ctx context.Context, k pageKeeper) error {
	return imp.set.Walk(func(page *types.Page) (err error) {
		var handle = page.Handle

		if page.ID == 0 {
			page.NamespaceID = imp.namespace.ID
			page, err = k.Create(page)
		} else {
			page, err = k.Update(page)
		}

		// @todo where do we check if page with module ref already exists?
		// @todo store pages & resolve page's parent ref!

		if err != nil {
			return
		}

		imp.permissions.UpdateResources(types.PagePermissionResource.String(), handle, page.ID)

		return
	})
}
