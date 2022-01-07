package gig

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/store"
)

func (w *workerEnvoy) filterComposeNamespace(base *store.DecodeFilter, defs preprocessorResourceLoad) (out *store.DecodeFilter) {
	out = base

	if defs.id != 0 {
		out = out.ComposeNamespace(&types.NamespaceFilter{
			NamespaceID: []uint64{defs.id},
		})
	}

	if defs.handle != "" {
		out = out.ComposeNamespace(&types.NamespaceFilter{
			Slug: defs.handle,
		})
	}

	return
}

func (w *workerEnvoy) filterComposeModule(base *store.DecodeFilter, defs preprocessorResourceLoad) (out *store.DecodeFilter) {
	out = base

	if defs.id != 0 {
		out = out.ComposeModule(&types.ModuleFilter{
			ModuleID: []uint64{defs.id},
		})
	}

	if defs.handle != "" {
		out = out.ComposeModule(&types.ModuleFilter{
			Handle: defs.handle,
		})
	}

	return
}

func (w *workerEnvoy) filterComposeChart(base *store.DecodeFilter, defs preprocessorResourceLoad) (out *store.DecodeFilter) {
	out = base

	if defs.id != 0 {
		out = out.ComposeChart(&types.ChartFilter{
			ChartID: []uint64{defs.id},
		})
	}

	if defs.handle != "" {
		out = out.ComposeChart(&types.ChartFilter{
			Handle: defs.handle,
		})
	}

	return
}

func (w *workerEnvoy) filterComposePage(base *store.DecodeFilter, defs preprocessorResourceLoad) (out *store.DecodeFilter) {
	out = base

	// if defs.id != 0 {
	// 	out = out.ComposePage(&types.PageFilter{
	// 		PageID: defs.id,
	// 	})
	// }

	if defs.handle != "" {
		out = out.ComposePage(&types.PageFilter{
			Handle: defs.handle,
		})
	}

	return
}

func (w *workerEnvoy) noop(_ context.Context, _ preprocessorNoop) error {
	return nil
}

func (w *workerEnvoy) resourceRemove(ctx context.Context, t preprocessorResourceRemove) error {
	var ref *resource.Ref
	if t.resource != "" {
		if t.identifier != "" {
			ref = resource.MakeRef(t.resource, resource.MakeIdentifiers(t.identifier))
		} else {
			ref = resource.MakeWildRef(t.resource)
		}
	}

	if t.identifier == "" && t.resource == "" {
		return fmt.Errorf("invalid parameters: at least resType must be provided")
	}

	out := make([]resource.Interface, 0, len(w.resources))
	for _, r := range w.resources {
		if r.ResourceType() != ref.ResourceType && !r.Identifiers().HasAny(ref.Identifiers) {
			out = append(out, r)
		}
	}

	w.resources = out
	return nil
}

func (w *workerEnvoy) resourceLoad(ctx context.Context, params preprocessorResourceLoad) error {
	df := store.NewDecodeFilter()

	switch params.resource {
	case ComposeNamespaceResourceType:
		df = w.filterComposeNamespace(df, params)
	case ComposeModuleResourceType:
		df = w.filterComposeModule(df, params)
	case ComposeChartResourceType:
		df = w.filterComposeChart(df, params)
	case ComposePageResourceType:
		df = w.filterComposePage(df, params)
	}

	res, err := w.getStoreDecoders().Decode(ctx, w.store, df)
	if err != nil {
		return err
	}

	w.resources = append(w.resources, res...)
	return nil
}

func (w *workerEnvoy) namespaceLoad(ctx context.Context, params preprocessorNamespaceLoad) error {
	df := store.NewDecodeFilter()

	if params.id != 0 {
		df = df.ComposeNamespace(&types.NamespaceFilter{
			NamespaceID: []uint64{params.id},
		})
	} else {
		df = df.ComposeNamespace(&types.NamespaceFilter{
			Slug: params.handle,
		})
	}

	df = df.
		ComposeModule(&types.ModuleFilter{}).
		ComposePage(&types.PageFilter{}).
		ComposeChart(&types.ChartFilter{})

	res, err := w.getStoreDecoders().Decode(ctx, w.store, df)
	if err != nil {
		return err
	}

	w.resources = append(w.resources, res...)
	return nil
}

func preprocessorResourceRemoveTransformer(base preprocessorResourceRemove) preprocessorResourceRemove {
	if base.identifier == "" {
		base.identifier = "*"
	}

	return base
}
