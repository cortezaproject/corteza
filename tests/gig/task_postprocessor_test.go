package gig

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
	"github.com/cortezaproject/corteza-server/store"
)

func test_postprocessor_tasks_worker(t *testing.T, h helper, s store.Storer) gig.Worker {
	return gig.WorkerNoop()
}

func test_postprocessor_tasks_noop(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, tc string) {
	g, err := svc.SetPostprocessors(ctx, g, gig.PostprocessorNoop())
	h.a.NoError(err)

	ss, err := gig.PrepareSourceFromDirectory(ctx, testSource(t, "_base"))
	h.a.NoError(err)
	g, err = svc.AddSources(ctx, g, gig.ToSourceWrap(ss...))
	h.a.NoError(err)

	g, err = svc.Exec(ctx, g)
	h.a.NoError(err)

	out, err := svc.Output(ctx, g)
	h.a.NoError(err)

	h.a.Len(out, 2)
	h.a.True(gig.ToSourceWrap(out...).HasByName("secret-key.txt"))
	h.a.True(gig.ToSourceWrap(out...).HasByName("secret-message.txt"))
}

func test_postprocessor_tasks_discard(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, tc string) {
	g, err := svc.SetPostprocessors(ctx, g, gig.PostprocessorDiscard())
	h.a.NoError(err)

	ss, err := gig.PrepareSourceFromDirectory(ctx, testSource(t, "_base"))
	h.a.NoError(err)
	g, err = svc.AddSources(ctx, g, gig.ToSourceWrap(ss...))
	h.a.NoError(err)

	g, err = svc.Exec(ctx, g)
	h.a.NoError(err)

	out, err := svc.Output(ctx, g)
	h.a.NoError(err)

	h.a.Len(out, 0)
}

// @todo pending implementation
func test_postprocessor_tasks_save(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, tc string) {
}

func test_postprocessor_tasks_archive(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, tc string) {
	g, err := svc.SetPostprocessors(ctx, g, gig.PostprocessorArchive(gig.ArchiveTar))
	h.a.NoError(err)

	ss, err := gig.PrepareSourceFromDirectory(ctx, testSource(t, "_base"))
	h.a.NoError(err)
	g, err = svc.AddSources(ctx, g, gig.ToSourceWrap(ss...))
	h.a.NoError(err)

	g, err = svc.Exec(ctx, g)
	h.a.NoError(err)

	out, err := svc.Output(ctx, g)
	h.a.NoError(err)

	h.a.Len(out, 1)
	h.a.True(gig.ToSourceWrap(out...).HasByName("archive.tar.gz"))
}
