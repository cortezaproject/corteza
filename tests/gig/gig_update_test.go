package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_gig_update(t *testing.T) {
	var (
		ctx, svc, h, s, g = setupWithNoopGig(t)
		err               error
	)
	_ = s
	_ = svc
	_ = g

	fork := func(g gig.Gig) gig.Gig {
		return g
	}

	t.Run("invalid gig update attempts", func(_ *testing.T) {
		// worker can't be changed
		_, err = svc.Update(ctx, g, gig.UpdatePayload{
			Worker: &workerFail{},
		})
		h.a.Error(err)

		// prepared gigs can't be updated
		ng := fork(g)
		ng.PreparedAt = now()
		_, err = svc.Update(ctx, ng, gig.UpdatePayload{})
		h.a.Error(err)

		// completed gigs can't be updated
		ng = fork(g)
		ng.CompletedAt = now()
		_, err = svc.Update(ctx, ng, gig.UpdatePayload{})
		h.a.Error(err)
	})

	t.Run("update decoders", func(_ *testing.T) {
		t.Run("implicit with source", func(_ *testing.T) {
			var ng gig.Gig
			_ = ng

			ng, err = svc.Update(ctx, g, gig.UpdatePayload{
				Decode: gig.DecoderSet{
					decoderSafe(gig.DecoderNoopSource(0)),
				},
				Sources: []gig.SourceWrap{
					{
						Uri:  testSource(t, "secret-key.txt"),
						Name: "secret-key.txt",
					},
				},
			})
			h.a.NoError(err)

			h.a.Len(ng.Sources, 1)
			h.a.Len(ng.Sources[0].Decoders(), 1)
			h.a.NotNil(ng.UpdatedAt)
		})

		t.Run("explicit after source", func(_ *testing.T) {
			g, err := svc.Create(ctx, gig.UpdatePayload{
				Worker: gig.WorkerNoop(),
				Sources: []gig.SourceWrap{
					{
						Uri:  testSource(t, "secret-key.txt"),
						Name: "secret-key.txt",
					},
				},
			})
			h.a.NoError(err)

			sourceID := g.Sources[0].ID()
			var ng gig.Gig
			_ = ng

			ng, err = svc.Update(ctx, g, gig.UpdatePayload{
				Sources: gig.ToSourceWrap(g.Sources...),
				Decode: gig.DecoderSet{
					decoderSafe(gig.DecoderNoopSource(sourceID)),
				},
			})
			h.a.NoError(err)

			h.a.Len(ng.Sources, 1)
			h.a.Len(ng.Sources[0].Decoders(), 1)
			h.a.NotNil(ng.UpdatedAt)
		})
	})

	t.Run("update preprocessors", func(_ *testing.T) {
		var ng gig.Gig
		_ = ng

		ng, err = svc.Update(ctx, g, gig.UpdatePayload{
			Preprocess: gig.PreprocessorSet{
				gig.PreprocessorNoop(),
			},
		})
		h.a.NoError(err)
		h.a.Len(ng.Preprocess, 1)
		h.a.NotNil(ng.UpdatedAt)
	})

	t.Run("update postprocessors", func(_ *testing.T) {
		var ng gig.Gig
		_ = ng

		ng, err = svc.Update(ctx, g, gig.UpdatePayload{
			Postprocess: gig.PostprocessorSet{
				gig.PostprocessorNoop(),
			},
		})
		h.a.NoError(err)
		h.a.Len(ng.Postprocess, 1)
		h.a.NotNil(ng.UpdatedAt)
	})
}
