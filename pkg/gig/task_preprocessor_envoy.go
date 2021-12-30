package gig

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/store"
)

func (w *workerEnvoy) resourceRemove(ctx context.Context, t preprocessorResourceRemove) error {
	var ref *resource.Ref
	if t.Resource != "" {
		if t.Identifier != "" {
			ref = resource.MakeRef(t.Resource, resource.MakeIdentifiers(t.Identifier))
		} else {
			ref = resource.MakeWildRef(t.Resource)
		}
	}

	if t.Identifier == "" && t.Resource == "" {
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

	switch params.Resource {
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

	if params.ID != 0 {
		df = df.ComposeNamespace(&types.NamespaceFilter{
			NamespaceID: []uint64{params.ID},
		})
	} else {
		df = df.ComposeNamespace(&types.NamespaceFilter{
			Slug: params.Handle,
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
