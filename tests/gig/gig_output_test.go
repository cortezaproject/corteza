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
		out, ng, err = svc.OutputAll(ctx, g)
		h.a.NoError(err)
		h.a.Nil(out)
	})

	t.Run("output list", func(_ *testing.T) {
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

		var list gig.SourceWrapSet
		list, ng, err = svc.Output(ctx, ng)
		h.a.NoError(err)
		h.a.Len(list, 1)
		h.a.NotEqual(0, list[0].ID)

		t.Run("specific output", func(t *testing.T) {
			var out gig.Source
			out, ng, err = svc.OutputSpecific(ctx, ng, list[0].ID)
			h.a.NoError(err)
			h.a.NotNil(out)
			h.a.Equal(list[0].ID, out.ID())
		})

		t.Run("specific output; not found", func(t *testing.T) {
			_, ng, err = svc.OutputSpecific(ctx, ng, 1)
			h.a.Error(err)
		})
	})

	t.Run("all outputs", func(_ *testing.T) {
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

		out, ng, err = svc.OutputAll(ctx, ng)
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

		_, ng, err = svc.OutputAll(ctx, ng)
		h.a.NoError(err)
		h.a.NotNil(ng.CompletedAt)
	})
}
