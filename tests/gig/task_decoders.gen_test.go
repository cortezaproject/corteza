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

func Test_decoder_tasks(t *testing.T) {
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
			Worker: test_decoder_tasks_worker(t, h, s),
		})
		h.a.NoError(err)

		test_decoder_tasks_noop(ctx, t, h, svc, s, g, "noop")
	})
	t.Run("archive", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_decoder_tasks_worker(t, h, s),
		})
		h.a.NoError(err)

		test_decoder_tasks_archive(ctx, t, h, svc, s, g, "archive")
	})
}
