package gig

import (
	"testing"
)

func Test_source_management(t *testing.T) {
	// var (
	// 	ctx, svc, h, s, g = setupWithImportGig(t)
	// 	err               error
	// )
	// _ = s
	// _ = svc
	// _ = g

	// var srcA, srcB gig.Source
	// t.Run("add srcB source", func(_ *testing.T) {
	// 	srcA, err = gig.FileSourceFromURI(ctx, testSource(t, "secret-message.txt"))
	// 	h.a.NoError(err)
	// 	srcID := srcA.ID()

	// 	h.a.NotEqual(0, srcA.ID())
	// 	h.a.Equal("secret-message.txt", srcA.FileName())
	// 	h.a.Equal("text/plain; charset=utf-8", srcA.MimeType())
	// 	h.a.Equal(int64(29), srcA.Size())
	// 	h.a.Equal("d5ef1751bdb3e37742c7f44b1940b7f4800cd243af5788265faac808fda96e1b", srcA.Checksum())

	// 	g, err = svc.AddSources(ctx, g, gig.SourceSet{srcA})
	// 	h.a.NoError(err)

	// 	h.a.Len(g.Sources, 1)
	// 	h.a.Equal(srcID, g.Sources[0].ID())
	// })

	// t.Run("upsert sources", func(_ *testing.T) {
	// 	srcB, err = gig.FileSourceFromURI(ctx, testSource(t, "secret-key.txt"))
	// 	h.a.NoError(err)
	// 	srcID := srcB.ID()

	// 	g, err = svc.AddSources(ctx, g, gig.SourceSet{srcB, srcA})
	// 	h.a.NoError(err)

	// 	h.a.Len(g.Sources, 2)
	// 	// new source
	// 	h.a.Equal(srcID, g.Sources[1].ID())
	// 	// Original source
	// 	h.a.Equal(srcA.ID(), g.Sources[0].ID())
	// 	h.a.Equal("d5ef1751bdb3e37742c7f44b1940b7f4800cd243af5788265faac808fda96e1b", g.Sources[0].Checksum())
	// })

	// t.Run("remote source", func(_ *testing.T) {
	// 	ts, url := makeRemoteServer(t, "secret-key.txt")
	// 	defer ts.Close()

	// 	src, err := gig.FileSourceFromURI(ctx, url)
	// 	h.a.NoError(err)

	// 	h.a.NotEqual(0, src.ID())
	// 	h.a.Equal("secret-key.txt", src.FileName())
	// 	h.a.Equal("text/plain; charset=utf-8", src.MimeType())
	// 	h.a.Equal(int64(1), src.Size())
	// 	h.a.Equal("559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd", src.Checksum())
	// })

	// t.Run("cleanup", func(_ *testing.T) {
	// 	// Make sure the files exist
	// 	_, err = os.Open(srcA.Name())
	// 	h.a.NoError(err)
	// 	_, err = os.Open(srcB.Name())
	// 	h.a.NoError(err)

	// 	_, err = svc.Complete(ctx, g)
	// 	h.a.NoError(err)

	// 	// Make sure the files were removed
	// 	_, err = os.Open(srcA.Name())
	// 	h.a.Error(err)
	// 	_, err = os.Open(srcB.Name())
	// 	h.a.Error(err)
	// })
}
