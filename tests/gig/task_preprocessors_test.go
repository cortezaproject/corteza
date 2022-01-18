package gig

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
	"github.com/cortezaproject/corteza-server/store"
)

func test_preprocessor_tasks_worker_noop(t *testing.T, h helper, s store.Storer) gig.Worker {
	return gig.WorkerNoop()
}

func test_preprocessor_tasks_worker_attachment(t *testing.T, h helper, s store.Storer) gig.Worker {
	return gig.WorkerAttachment()
}

func test_preprocessor_tasks_worker_export(t *testing.T, h helper, s store.Storer) gig.Worker {
	return gig.WorkerExport(s)
}

func test_preprocessor_tasks_worker_import(t *testing.T, h helper, s store.Storer) gig.Worker {
	return gig.WorkerImport(s)
}

// Worker noop

func test_preprocessor_tasks_noop_noop(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, worker, task string) {
	g, err := svc.SetPreprocessors(ctx, g, gig.PreprocessorNoop())
	h.a.NoError(err)

	g, err = svc.Prepare(ctx, g)
	h.a.NoError(err)
}

// Worker attachment
// @todo

func test_preprocessor_tasks_attachment_attachmentRemove(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, worker, task string) {

}

func test_preprocessor_tasks_attachment_attachmentTransform(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, worker, task string) {

}

func test_preprocessor_tasks_attachment_noop(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, worker, task string) {

}

// Worker envoy

func test_preprocessor_tasks_export_experimentalExport(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, worker, task string) {
	// id := id.Next()
	// ns := &types.Namespace{
	// 	ID:   id,
	// 	Slug: "ns1",
	// }
	// err := store.CreateComposeNamespace(ctx, s, ns)
	// h.a.NoError(err)
	// defer store.DeleteComposeNamespace(ctx, s, ns)

	// g, err = svc.SetPreprocessors(ctx, g,
	// 	gig.PreprocessorNamespaceLoadHandle("ns1"),
	// 	gig.PreprocessorResourceRemove(gig.ResourceTranslationResourceType))
	// h.a.NoError(err)

	// g, err = svc.Prepare(ctx, g)
	// h.a.NoError(err)

	// ss, err := svc.State(ctx, g)
	// h.a.NoError(err)

	// cs := ss.(gig.WorkerStateEnvoy)
	// h.a.Len(cs.Resources, 1)
	// h.a.Equal(gig.ComposeNamespaceResourceType, cs.Resources[0].ResourceType)
}

func test_preprocessor_tasks_export_noop(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, worker, task string) {
	g, err := svc.AddSources(ctx, g, gig.SourceWrapSet{{
		Uri: testSource(t, "_base", "namespace.yaml"),
	}})
	h.a.NoError(err)

	g, err = svc.SetPreprocessors(ctx, g, gig.PreprocessorNoop())
	h.a.NoError(err)

	g, err = svc.Prepare(ctx, g)
	h.a.NoError(err)

	ss, err := svc.State(ctx, g)
	h.a.NoError(err)

	cs := ss.(gig.WorkerStateEnvoy)
	h.a.Len(cs.Resources, 2)
}

func test_preprocessor_tasks_import_noop(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, worker, task string) {
	g, err := svc.AddSources(ctx, g, gig.SourceWrapSet{{
		Uri: testSource(t, "_base", "namespace.yaml"),
	}})
	h.a.NoError(err)

	g, err = svc.SetPreprocessors(ctx, g, gig.PreprocessorNoop())
	h.a.NoError(err)

	g, err = svc.Prepare(ctx, g)
	h.a.NoError(err)

	ss, err := svc.State(ctx, g)
	h.a.NoError(err)

	cs := ss.(gig.WorkerStateEnvoy)
	h.a.Len(cs.Resources, 2)
}
