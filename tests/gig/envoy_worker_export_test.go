package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_envoy_worker_export(t *testing.T) {
	var (
		ctx, svc, h, s, g = setupWithExportGig(t)
		err               error
	)
	loadScenario(ctx, s, t, h)

	p1 := gig.PreprocessorResourceLoadByHandle(gig.ComposeNamespaceResourceType, "ns1")
	p2 := gig.PreprocessorResourceRemove(gig.ResourceTranslationResourceType, "*")

	g, err = svc.SetPreprocessors(ctx, g, p1, p2)
	h.a.NoError(err)

	g, err = svc.Exec(ctx, g)
	h.a.NoError(err)

	h.a.Len(g.Output, 1)
}
