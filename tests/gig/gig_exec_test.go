package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_gig_exec(t *testing.T) {
	var (
		ctx, svc, h, s, g = setupWithNoopGig(t)
		err               error
	)
	_ = s
	_ = svc
	_ = g

	var ng gig.Gig
	t.Run("prepare -> run", func(_ *testing.T) {
		ng, err = svc.Prepare(ctx, g)
		h.a.NoError(err)
		h.a.NotNil(ng.PreparedAt)

		ng, err = svc.Exec(ctx, ng)
		h.a.NoError(err)
		h.a.NotNil(ng.Status.StartedAt)
		h.a.NotNil(ng.Status.CompletedAt)
		h.a.NotNil(ng.Status.Elapsed)
	})

	t.Run("prepare without a worker", func(_ *testing.T) {
		ng = gig.Gig{ID: 1}
		ng, err = svc.Prepare(ctx, ng)
		h.a.Error(err)
	})

	t.Run("re-prepare worker", func(_ *testing.T) {
		ng, err = svc.Prepare(ctx, g)
		h.a.NoError(err)
		ng, err = svc.Prepare(ctx, ng)
		h.a.Error(err)
	})

	t.Run("implicit prepare", func(_ *testing.T) {
		ng, err = svc.Exec(ctx, g)
		h.a.NoError(err)
		h.a.NotNil(ng.Status.StartedAt)
		h.a.NotNil(ng.Status.CompletedAt)
		h.a.NotNil(ng.Status.Elapsed)
	})
}
