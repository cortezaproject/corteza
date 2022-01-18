package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_decoder_archive(t *testing.T) {
	var (
		ctx, svc, h, s, g = setupWithNoopGig(t)
		err               error
		tmp               gig.SourceSet
	)
	_ = s
	_ = svc
	_ = g

	var ng gig.Gig
	t.Run(".tar.gz", func(_ *testing.T) {
		dc := decoderSafe(gig.DecoderArchiveSource(0))

		ng, err = svc.AddSources(ctx, g, gig.SourceWrapSet{{
			Uri: testSource(t, "archive.tar.gz"),
		}}, dc)
		h.a.NoError(err)

		ng, err = svc.Prepare(ctx, ng)
		h.a.NoError(err)

		ng, err = svc.Exec(ctx, ng)
		h.a.NoError(err)

		tmp = ng.Output

		h.a.Len(tmp, 2)
		h.a.Equal("text/yaml; charset=utf-8", tmp[0].MimeType())
		h.a.Equal("text/yaml; charset=utf-8", tmp[1].MimeType())
	})

	t.Run(".zip", func(_ *testing.T) {
		dc := decoderSafe(gig.DecoderArchiveSource(0))

		ng, err = svc.AddSources(ctx, g, gig.SourceWrapSet{{
			Uri: testSource(t, "archive.zip"),
		}}, dc)
		h.a.NoError(err)

		ng, err = svc.Prepare(ctx, ng)
		h.a.NoError(err)

		ng, err = svc.Exec(ctx, ng)
		h.a.NoError(err)

		tmp = ng.Output

		h.a.Len(tmp, 2)
		h.a.Equal("text/yaml; charset=utf-8", tmp[0].MimeType())
		h.a.Equal("text/yaml; charset=utf-8", tmp[1].MimeType())
	})
}
