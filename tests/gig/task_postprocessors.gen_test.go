package gig

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_postprocessor_tasks(t *testing.T) {
	var (
		ctx, svc, h, s = setup(t)
		err            error
		g              gig.Gig
	)
	_ = s
	_ = svc
	_ = err

	t.Run("noop", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_postprocessor_tasks_worker(t, h, s),
		})
		h.a.NoError(err)

		test_postprocessor_tasks_noop(ctx, t, h, svc, s, g, "noop")
	})
	t.Run("discard", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_postprocessor_tasks_worker(t, h, s),
		})
		h.a.NoError(err)

		test_postprocessor_tasks_discard(ctx, t, h, svc, s, g, "discard")
	})
	t.Run("save", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_postprocessor_tasks_worker(t, h, s),
		})
		h.a.NoError(err)

		test_postprocessor_tasks_save(ctx, t, h, svc, s, g, "save")
	})
	t.Run("archive", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_postprocessor_tasks_worker(t, h, s),
		})
		h.a.NoError(err)

		test_postprocessor_tasks_archive(ctx, t, h, svc, s, g, "archive")
	})
}
