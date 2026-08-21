package agentic

import (
	"context"
	"strconv"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/weburl"
)

// Links a single-item result carries, so a caller that has just changed
// something can say where to go and look at it.
//
// A link is a convenience and never a reason to fail: every builder here
// returns "" rather than an error when it cannot resolve what it needs, and
// toolkit.JSONResultWith drops an empty value instead of emitting a blank key.
// The alternative is a write that succeeded reported as a failure because a
// slug lookup did not.

// nsURLPart is what goes in the :slug segment of every compose route.
//
// The webapp resolves that segment by slug OR by namespaceID
// (lib/vue/src/stores/useNamespaceStore.js getByUrlPart), and a slug is
// optional on a namespace — plenty of real ones have none. Falling back to the
// ID is what keeps those linkable at all; keying on the slug alone silently
// produced no link for them.
func nsURLPart(ctx context.Context, nsID uint64) string {
	ns, err := cmpService.DefaultNamespace.FindByID(ctx, nsID)
	if err != nil || ns == nil {
		return ""
	}
	if ns.Slug != "" {
		return ns.Slug
	}
	return strconv.FormatUint(ns.ID, 10)
}

func namespaceLinks(ns *cmpTypes.Namespace) map[string]string {
	if ns == nil {
		return nil
	}
	part := ns.Slug
	if part == "" {
		part = strconv.FormatUint(ns.ID, 10)
	}
	return map[string]string{
		"url":     weburl.ComposeNamespace(part),
		"editUrl": weburl.ComposeNamespaceEdit(part),
	}
}

// pageLinks points at the page as a user sees it, and at the builder, which is
// where its blocks are arranged.
func pageLinks(ctx context.Context, pg *cmpTypes.Page) map[string]string {
	if pg == nil {
		return nil
	}
	slug := nsURLPart(ctx, pg.NamespaceID)
	return map[string]string{
		"url":     weburl.ComposePage(slug, pg.ID),
		"editUrl": weburl.ComposePageBuilder(slug, pg.ID),
	}
}

// layoutLinks point at the builder too: a layout has no screen of its own, it
// is what the builder edits.
func layoutLinks(ctx context.Context, layout *cmpTypes.PageLayout) map[string]string {
	if layout == nil {
		return nil
	}
	return map[string]string{"url": weburl.ComposePageBuilder(nsURLPart(ctx, layout.NamespaceID), layout.PageID)}
}

func moduleLinks(ctx context.Context, mod *cmpTypes.Module) map[string]string {
	if mod == nil {
		return nil
	}
	return map[string]string{"url": weburl.ComposeModuleEdit(nsURLPart(ctx, mod.NamespaceID), mod.ID)}
}

func chartLinks(ctx context.Context, chart *cmpTypes.Chart) map[string]string {
	if chart == nil {
		return nil
	}
	return map[string]string{"url": weburl.ComposeChartEdit(nsURLPart(ctx, chart.NamespaceID), chart.ID)}
}

// recordLinks use the module admin view, the one screen every record has. A
// record also opens on any page bound to its module, but which page that should
// be is the caller's choice and not something a link can guess.
func recordLinks(ctx context.Context, rec *cmpTypes.Record) map[string]string {
	if rec == nil {
		return nil
	}
	return map[string]string{"url": weburl.ComposeRecord(nsURLPart(ctx, rec.NamespaceID), rec.ModuleID, rec.ID)}
}
