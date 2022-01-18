package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_gig_output(t *testing.T) {
	var (
		ctx, svc, h, s, g = setupWithNoopGig(t)
		err               error
		out               gig.SourceSet
	)
	_ = s
	_ = svc
	_ = g

	var ng gig.Gig
	t.Run("new gig", func(_ *testing.T) {
		out, ng, err = svc.Output(ctx, g)
		h.a.NoError(err)
		h.a.Nil(out)
	})

	t.Run("post exec", func(_ *testing.T) {
		ng, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: gig.WorkerNoop(),
			Sources: []gig.SourceWrap{{
				Uri: testSource(t, "secret-key.txt"),
			}},
			Postprocess: gig.PostprocessorSet{
				postprocessorSafe(gig.PostprocessorNoopParams(nil)),
			},
		})
		h.a.NoError(err)

		ng, err = svc.Exec(ctx, ng)
		h.a.NoError(err)

		out, ng, err = svc.Output(ctx, ng)
		h.a.NoError(err)
		h.a.NotNil(out)
	})

	t.Run("complete on output", func(_ *testing.T) {
		ng, err = svc.Create(ctx, gig.UpdatePayload{
			Worker:     gig.WorkerNoop(),
			CompleteOn: gig.OnOutput,
		})

		ng, err = svc.Exec(ctx, ng)
		h.a.NoError(err)
		h.a.Nil(ng.CompletedAt)

		_, ng, err = svc.Output(ctx, ng)
		h.a.NoError(err)
		h.a.NotNil(ng.CompletedAt)
	})
}
