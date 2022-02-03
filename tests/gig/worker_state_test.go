package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_worker_state(t *testing.T) {
	var (
		ctx, svc, h, s = setup(t)
		err            error
		g              gig.Gig
	)
	_ = s
	_ = svc
	_ = err

	prep := func(w gig.Worker, src string) gig.Gig {
		g, err := svc.Create(ctx, gig.UpdatePayload{
			Worker: w,
			Sources: []gig.SourceWrap{{
				Uri: testSource(t, src),
			}},
		})
		h.a.NoError(err)

		g, err = svc.Prepare(ctx, g)
		h.a.NoError(err)

		return g
	}

	t.Run("completed gig", func(_ *testing.T) {
		ng, err := svc.Complete(ctx, g)
		h.a.NoError(err)

		_, err = svc.State(ctx, ng)
		h.a.Error(err)
	})

	t.Run("noop", func(_ *testing.T) {
		g = prep(gig.WorkerNoop(), "secret-message.txt")

		ss, err := svc.State(ctx, g)
		h.a.NoError(err)
		h.a.NotNil(ss)
	})

	t.Run("envoy", func(_ *testing.T) {
		g = prep(gig.WorkerImport(s), "namespace.yaml")

		ss, err := svc.State(ctx, g)
		h.a.NoError(err)
		h.a.NotNil(ss)
	})
}
