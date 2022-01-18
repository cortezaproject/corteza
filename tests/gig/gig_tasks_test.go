package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_gig_tasks(t *testing.T) {
	var (
		ctx, svc, h, _ = setup(t)
		g              gig.Gig
		err            error
	)

	t.Run("decoders", func(_ *testing.T) {
		t.Run("ok", func(_ *testing.T) {
			g, err = svc.Create(ctx, gig.UpdatePayload{
				Worker: gig.WorkerNoop(),
				Sources: []gig.SourceWrap{{
					Uri: testSource(t, "secret-key.txt"),
				}},
			})
			h.a.NoError(err)

			g, err = svc.SetDecoders(ctx, g, decoderSafe(gig.DecoderNoopSource(g.Sources[0].ID())))
			h.a.NoError(err)

			h.a.Len(g.Sources[0].Decoders(), 1)
		})

		t.Run("ok", func(_ *testing.T) {
			g, err = svc.Create(ctx, gig.UpdatePayload{
				Worker: gig.WorkerNoop(),
				Sources: []gig.SourceWrap{{
					Uri: testSource(t, "secret-key.txt"),
				}},
			})
			h.a.NoError(err)

			g, err = svc.Prepare(ctx, g)
			h.a.NoError(err)

			g, err = svc.SetDecoders(ctx, g, decoderSafe(gig.DecoderNoopSource(g.Sources[0].ID())))
			h.a.Error(err)
		})
	})

	t.Run("preprocessors", func(_ *testing.T) {
		t.Run("ok", func(_ *testing.T) {
			g, err = svc.Create(ctx, gig.UpdatePayload{
				Worker: gig.WorkerNoop(),
				Sources: []gig.SourceWrap{{
					Uri: testSource(t, "secret-key.txt"),
				}},
			})
			h.a.NoError(err)

			g, err = svc.SetPreprocessors(ctx, g, preprocessorSafe(gig.PreprocessorNoopParams(nil)))
			h.a.NoError(err)
		})

		t.Run("ok", func(_ *testing.T) {
			g, err = svc.Create(ctx, gig.UpdatePayload{
				Worker: gig.WorkerNoop(),
				Sources: []gig.SourceWrap{{
					Uri: testSource(t, "secret-key.txt"),
				}},
			})
			h.a.NoError(err)

			g, err = svc.Prepare(ctx, g)
			h.a.NoError(err)

			g, err = svc.SetPreprocessors(ctx, g, preprocessorSafe(gig.PreprocessorNoopParams(nil)))
			h.a.Error(err)
		})
	})

	t.Run("postprocessors", func(_ *testing.T) {
		t.Run("ok", func(_ *testing.T) {
			g, err = svc.Create(ctx, gig.UpdatePayload{
				Worker: gig.WorkerNoop(),
				Sources: []gig.SourceWrap{{
					Uri: testSource(t, "secret-key.txt"),
				}},
			})
			h.a.NoError(err)

			g, err = svc.SetPostprocessors(ctx, g, postprocessorSafe(gig.PostprocessorNoopParams(nil)))
			h.a.NoError(err)
		})

		t.Run("ok", func(_ *testing.T) {
			g, err = svc.Create(ctx, gig.UpdatePayload{
				Worker: gig.WorkerNoop(),
				Sources: []gig.SourceWrap{{
					Uri: testSource(t, "secret-key.txt"),
				}},
			})
			h.a.NoError(err)

			g, err = svc.Prepare(ctx, g)
			h.a.NoError(err)

			g, err = svc.SetPostprocessors(ctx, g, postprocessorSafe(gig.PostprocessorNoopParams(nil)))
			h.a.Error(err)
		})
	})
}
