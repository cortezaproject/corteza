package gig

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
	"github.com/cortezaproject/corteza-server/store"
)

func test_decoder_tasks_worker(t *testing.T, h helper, s store.Storer) gig.Worker {
	return gig.WorkerNoop()
}

func test_decoder_tasks_noop(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, tc string) {
	g, err := svc.AddSources(ctx, g, gig.SourceWrapSet{{
		Uri: testSource(t, tc, "secret-key.txt"),
	}}, decoderSafe(gig.DecoderNoopSource(0)))
	h.a.NoError(err)

	g, err = svc.Prepare(ctx, g)
	h.a.NoError(err)

	state, err := svc.State(ctx, g)
	h.a.NoError(err)

	cs := state.(gig.WorkerNoopState)
	h.a.Len(cs.Sources, 1)

	h.a.Equal("secret-key.txt", cs.Sources[0].Name)
}

func test_decoder_tasks_archive(ctx context.Context, t *testing.T, h helper, svc gig.Service, s store.Storer, g gig.Gig, tc string) {
	g, err := svc.AddSources(ctx, g, gig.SourceWrapSet{{
		Uri: testSource(t, tc, "archive.tar.gz"),
	}}, decoderSafe(gig.DecoderArchiveSource(0)))
	h.a.NoError(err)

	g, err = svc.Prepare(ctx, g)
	h.a.NoError(err)

	state, err := svc.State(ctx, g)
	h.a.NoError(err)
	cs := state.(gig.WorkerNoopState)
	h.a.Len(cs.Sources, 2)

	h.a.True(cs.Sources.HasByName("namespaces.yaml"))
	h.a.True(cs.Sources.HasByName("modules.yaml"))
}
