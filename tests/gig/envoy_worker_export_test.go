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
	_ = svc
	_ = err
	_ = g

	p1 := preprocessorSafe(gig.PreprocessorExperimentalExportParams(map[string]interface{}{
		"handle":           "ns1",
		"inclRBAC":         true,
		"inclTranslations": true,
	}))
	g, err = svc.SetPreprocessors(ctx, g, p1)
	h.a.NoError(err)

	g, err = svc.Exec(ctx, g)
	h.a.NoError(err)

	expect := []string{
		"corteza::compose:namespace.yaml",
		"corteza::compose:module.yaml",
		"corteza::compose:page.yaml",
		"corteza::compose:chart.yaml",
		"resource-translation.yaml",
	}

	for _, e := range expect {
		h.a.NotNil(g.Output.GetByName(e))
	}
}
