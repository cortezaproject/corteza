package importer

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	Page struct {
		imp *Importer

		namespace *types.Namespace
		set       types.PageSet
		dirty     map[uint64]bool

		// page => module handle
		modules map[string]string

		// child => parent handle
		tree map[string][]string
	}

	// @todo remove finder strategy, directly provide set of items
	pageFinder interface {
		Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error)
	}
)

func NewPageImporter(imp *Importer, ns *types.Namespace) *Page {
	out := &Page{
		imp: imp,

		namespace: ns,

		set:   types.PageSet{},
		dirty: make(map[uint64]bool),

		modules: map[string]string{},
		tree:    map[string][]string{},
	}

	if imp.pageFinder != nil && ns.ID > 0 {
		out.set, _, _ = imp.pageFinder.Find(types.PageFilter{NamespaceID: ns.ID})
	}

	return out
}

func (pImp *Page) getModule(handle string) (*types.Module, error) {
	if g, ok := pImp.imp.namespaces.modules[pImp.namespace.Slug]; !ok {
		return nil, errors.Errorf("could not get modules %q from non existing namespace %q", handle, pImp.namespace.Slug)
	} else {
		return g.Get(handle)
	}
}

func (pImp *Page) getScript(name string) (*automation.Script, error) {
	if g, ok := pImp.imp.namespaces.scripts[pImp.namespace.Slug]; !ok {
		return nil, errors.Errorf("could not get scripts %q from non existing namespace %q", name, pImp.namespace.Slug)
	} else {
		return g.Get(name)
	}
}

func (pImp *Page) getChart(handle string) (*types.Chart, error) {
	if g, ok := pImp.imp.namespaces.charts[pImp.namespace.Slug]; !ok {
		return nil, errors.Errorf("could not get chart %q from non existing namespace %q", handle, pImp.namespace.Slug)
	} else {
		return g.Get(handle)
	}
}

// CastSet Resolves permission rules:
// { <page-handle>: { page } } or [ { page }, ... ]
func (pImp *Page) CastSet(set interface{}) error {
	return pImp.castSet("", set)
}

// CastSet Resolves permission rules:
// { <page-handle>: { page } } or [ { page }, ... ]
func (pImp *Page) castSet(parent string, set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Pages defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return pImp.cast(parent, handle, def)
	})
}

func (pImp *Page) Cast(handle string, def interface{}) (err error) {
	return pImp.cast("", handle, def)
}

// Cast Resolves permission rules:
// { <page-handle>: { page } } or [ { page }, ... ]
func (pImp *Page) cast(parent, handle string, def interface{}) (err error) {
	var page *types.Page

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid page handle")
	}

	handle = importer.NormalizeHandle(handle)

	if page, err = pImp.Get(handle); err != nil {
		return err
	} else if page == nil {
		page = &types.Page{
			Handle:  handle,
			Title:   handle,
			Visible: true,
		}

		pImp.set = append(pImp.set, page)
	} else if page.ID == 0 {
		return errors.Errorf("page handle %q already defined in this import session", page.Handle)
	} else {
		pImp.dirty[page.ID] = true
	}

	pImp.tree[parent] = append(pImp.tree[parent], handle)

	// Make pages are always sorted
	sort.Strings(pImp.tree[parent])

	if title, ok := def.(string); ok && title != "" {
		page.Title = title
		return nil
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "namespace":
			// namespace value sanity check
			if deinterfacer.ToString(val, pImp.namespace.Slug) != pImp.namespace.Slug {
				return fmt.Errorf("explicitly set namespace on page %q shadows inherited namespace", pImp.namespace.Slug)
			}

		case "handle":
			// handle value sanity check
			if deinterfacer.ToString(val, handle) != handle {
				return fmt.Errorf("explicitly set handle on page %q shadows inherited handle", handle)
			}

		case "module":
			pImp.modules[handle] = deinterfacer.ToString(val)

		case "visible":
			page.Visible = deinterfacer.ToBool(val)

		case "title", "name", "label":
			page.Title = deinterfacer.ToString(val)

		case "description":
			page.Description = deinterfacer.ToString(val)

		case "blocks":
			return pImp.castBlocks(page, val)

		case "pages":
			return pImp.castSet(handle, val)

		case "allow", "deny":
			return pImp.imp.permissions.CastSet(types.PagePermissionResource.String()+handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for page %q", key, handle)
		}

		return err
	})
}

func (pImp *Page) castBlocks(page *types.Page, def interface{}) error {
	page.Blocks = types.PageBlocks{}

	return deinterfacer.Each(def, func(b int, _ string, blockDef interface{}) (err error) {
		block := types.PageBlock{}

		err = deinterfacer.Each(blockDef, func(_ int, key string, val interface{}) (err error) {
			switch key {
			case "title", "name", "label":
				block.Title = deinterfacer.ToString(val)

			case "description":
				block.Description = deinterfacer.ToString(val)

			case "kind":
				block.Kind = deinterfacer.ToString(val)

			case "options":
				block.Options, err = pImp.castBlockOptions(val)
				return err

			case "style":
				block.Style, err = pImp.castBlockStyle(page, b, val)
				return

			case "XYWH", "xywh", "dim", "dimension":
				xywh := deinterfacer.ToInts(val)
				if len(xywh) != 4 {
					return errors.New("invalid dimension (xywh) value, expecting slice with 4 integers")
				}

				block.XYWH = [4]int{xywh[0], xywh[1], xywh[2], xywh[3]}

			default:
				return fmt.Errorf("unexpected key %q for block on page %q", key, page.Handle)

			}

			return nil
		})

		if err != nil {
			return err
		}

		page.Blocks = append(page.Blocks, block)
		return
	})
}

func (pImp *Page) castBlockOptions(def interface{}) (opt map[string]interface{}, err error) {
	opt = make(map[string]interface{})

	return opt, deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		opt[key] = deinterfacer.Simplify(val)
		return nil
	})
}

func (pImp *Page) castBlockStyle(page *types.Page, n int, def interface{}) (s types.PageBlockStyle, err error) {
	s = types.PageBlockStyle{}

	return s, deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "variants":
			s.Variants = map[string]string{}
			return deinterfacer.Each(val, func(_ int, key string, val interface{}) (err error) {
				s.Variants[key] = deinterfacer.ToString(val)
				return
			})
		default:
			return fmt.Errorf("unexpected key %q for block #%d on page %q", key, n+1, page.Handle)

		}
	})
}

// Get existing pages
func (pImp *Page) Get(handle string) (*types.Page, error) {
	handle = importer.NormalizeHandle(handle)
	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid page handle")
	}

	return pImp.set.FindByHandle(handle), nil
}

func (pImp *Page) Store(ctx context.Context, k pageKeeper) error {
	return pImp.storeChildren(ctx, "", k)
}

func (pImp *Page) storeChildren(ctx context.Context, parent string, k pageKeeper) (err error) {
	children, ok := pImp.tree[parent]
	if !ok {
		// No children...
		return nil
	}

	var page *types.Page

	for _, child := range children {
		if page, err = pImp.Get(child); err != nil {
			return
		}

		if err = pImp.resolveRefs(page); err != nil {
			return
		}

		if page.ID == 0 {
			page.NamespaceID = pImp.namespace.ID
			page, err = k.Create(page)
		} else if pImp.dirty[page.ID] {
			page, err = k.Update(page)
		}

		if err != nil {
			return
		}

		pImp.dirty[page.ID] = false
		if page.Handle == "" {
			continue
		}

		pImp.imp.permissions.UpdateResources(types.PagePermissionResource.String(), page.Handle, page.ID)

		if err = pImp.storeChildren(ctx, page.Handle, k); err != nil {
			return err
		}
	}

	return
}

// Resolve all refs for this page (page module, inside block)
func (pImp *Page) resolveRefs(page *types.Page) error {
	if moduleHandle, ok := pImp.modules[page.Handle]; ok {
		if module, err := pImp.getModule(moduleHandle); err != nil {
			return err
		} else if module == nil {
			return errors.Wrapf(err, "could not load module %q for page %q",
				moduleHandle, page.Handle)
		} else {
			page.ModuleID = module.ID
		}
	}

	for i, b := range page.Blocks {
		if b.Options == nil {
			continue
		}

		if h, ok := b.Options["module"]; ok {
			if refm, err := pImp.getModule(deinterfacer.ToString(h)); err != nil || refm == nil {
				return errors.Errorf("could not load module %q for page %q block #%d (%v)",
					h, page.Handle, i+1, err)
			} else {
				b.Options["moduleID"] = strconv.FormatUint(refm.ID, 10)
				delete(b.Options, "module")
			}
		}

		if h, ok := b.Options["page"]; ok {
			if refp, err := pImp.Get(deinterfacer.ToString(h)); err != nil || refp == nil {
				return errors.Errorf("could not load page %q for page %q block #%d (%v)",
					h, page.Handle, i+1, err)
			} else {
				b.Options["pageID"] = strconv.FormatUint(refp.ID, 10)
				delete(b.Options, "page")
			}
		}

		if h, ok := b.Options["chart"]; ok {
			if refc, err := pImp.getChart(deinterfacer.ToString(h)); err != nil || refc == nil {
				return errors.Errorf("could not load chart %q for page %q block #%d (%v)",
					h, page.Handle, i+1, err)
			} else {
				b.Options["chartID"] = strconv.FormatUint(refc.ID, 10)
				delete(b.Options, "chart")
			}
		}

		if b.Kind == "Automation" {
			bb := make([]interface{}, 0)
			err := deinterfacer.Each(b.Options["buttons"], func(_ int, _ string, btn interface{}) (err error) {
				button := map[string]interface{}{}

				err = deinterfacer.Each(btn, func(_ int, k string, v interface{}) error {
					switch k {
					case "script":
						if s, err := pImp.getScript(deinterfacer.ToString(v)); err != nil || s == nil {
							return errors.Errorf("could not load script %q for page %q block #%d (%v)",
								v, page.Handle, i+1, err)
						} else {
							button["scriptID"] = s.ID
						}
					default:
						button[k] = v
					}

					return nil
				})

				if err != nil {
					return err
				}

				bb = append(bb, button)
				return nil
			})

			b.Options["buttons"] = bb

			if err != nil {
				return err
			}
		}
	}

	return nil
}
