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

	t.Run(".tar.gz", func(_ *testing.T) {
		dc := decoderSafe(gig.DecoderArchiveSource(0))

		g, err = svc.AddSources(ctx, g, gig.SourceWrapSet{{
			Uri: testSource(t, "archive.tar.gz"),
		}}, dc)
		h.a.NoError(err)

		g, err = svc.Prepare(ctx, g)
		h.a.NoError(err)

		g, err = svc.Exec(ctx, g)
		h.a.NoError(err)

		tmp = g.Output

		h.a.Len(tmp, 2)
		h.a.Equal("text/yaml; charset=utf-8", tmp[0].MimeType())
		h.a.Equal("text/yaml; charset=utf-8", tmp[1].MimeType())
	})
}
