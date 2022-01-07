package gig

import (
	"os"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_source_management(t *testing.T) {
	var (
		ctx, svc, h, s, g = setupWithNoopGig(t)
		// err               error
	)
	_ = s
	_ = svc
	_ = g

	// var srcA, srcB gig.Source

	t.Run("add source", func(_ *testing.T) {
		ng, err := svc.AddSources(ctx, g, gig.SourceWrapSet{{
			Uri: testSource(t, "secret-message.txt"),
		}})
		h.a.NoError(err)
		h.a.Len(ng.Sources, 1)

		src := ng.Sources[0]
		h.a.NotEqual(0, src.ID())
		h.a.Equal("secret-message.txt", src.FileName())
		h.a.Equal("text/plain; charset=utf-8", src.MimeType())
		h.a.Equal(int64(29), src.Size())
		h.a.Equal("d5ef1751bdb3e37742c7f44b1940b7f4800cd243af5788265faac808fda96e1b", src.Checksum())

		ng, err = svc.AddSources(ctx, ng, gig.SourceWrapSet{{
			Uri: testSource(t, "secret-key.txt"),
		}})
		h.a.NoError(err)
		h.a.Len(ng.Sources, 2)
	})

	t.Run("set source", func(_ *testing.T) {
		ng, err := svc.SetSources(ctx, g, gig.SourceWrapSet{{
			Uri: testSource(t, "secret-message.txt"),
		}})
		h.a.NoError(err)
		h.a.Len(ng.Sources, 1)

		ng, err = svc.SetSources(ctx, ng, gig.SourceWrapSet{{
			Uri: testSource(t, "secret-key.txt"),
		}})
		h.a.NoError(err)
		h.a.Len(ng.Sources, 1)
	})

	t.Run("upsert sources", func(_ *testing.T) {
		ng, err := svc.AddSources(ctx, g, gig.SourceWrapSet{{
			Uri: testSource(t, "secret-message.txt"),
		}})
		h.a.NoError(err)
		h.a.Len(ng.Sources, 1)
		src := ng.Sources[0]

		ss := gig.SourceWrapSet{{
			ID: src.ID(),
		}, {
			Uri: testSource(t, "secret-key.txt"),
		}}

		ng, err = svc.AddSources(ctx, ng, ss)
		h.a.NoError(err)
		h.a.Len(ng.Sources, 2)
	})

	t.Run("remote source", func(_ *testing.T) {
		ts, url := makeRemoteServer(t, "secret-key.txt")
		defer ts.Close()

		ng, err := svc.AddSources(ctx, g, gig.SourceWrapSet{{
			Uri: url,
		}})
		h.a.NoError(err)
		h.a.Len(ng.Sources, 1)

		src := ng.Sources[0]
		h.a.NotEqual(0, src.ID())
		h.a.Equal("secret-key.txt", src.FileName())
		h.a.Equal("text/plain; charset=utf-8", src.MimeType())
		h.a.Equal(int64(1), src.Size())
		h.a.Equal("559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd", src.Checksum())

	})

	t.Run("cleanup", func(_ *testing.T) {
		ng, err := svc.SetSources(ctx, g, gig.SourceWrapSet{{
			Uri: testSource(t, "secret-message.txt"),
		}})
		h.a.NoError(err)
		src := ng.Sources[0]

		// Make sure the files exist
		_, err = os.Open(src.Name())
		h.a.NoError(err)

		_, err = svc.Complete(ctx, ng)
		h.a.NoError(err)

		// Make sure the files were removed
		_, err = os.Open(src.Name())
		h.a.Error(err)
	})
}
