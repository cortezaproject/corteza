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

func Test_preprocessor_tasks(t *testing.T) {
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
			Worker: test_preprocessor_tasks_worker_noop(t, h, s),
		})
		h.a.NoError(err)

		test_preprocessor_tasks_noop_noop(ctx, t, h, svc, s, g, "noop", "noop")
	})

	t.Run("attachmentRemove", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_preprocessor_tasks_worker_attachment(t, h, s),
		})
		h.a.NoError(err)

		test_preprocessor_tasks_attachment_attachmentRemove(ctx, t, h, svc, s, g, "attachment", "attachmentRemove")
	})
	t.Run("attachmentTransform", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_preprocessor_tasks_worker_attachment(t, h, s),
		})
		h.a.NoError(err)

		test_preprocessor_tasks_attachment_attachmentTransform(ctx, t, h, svc, s, g, "attachment", "attachmentTransform")
	})
	t.Run("noop", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_preprocessor_tasks_worker_attachment(t, h, s),
		})
		h.a.NoError(err)

		test_preprocessor_tasks_attachment_noop(ctx, t, h, svc, s, g, "attachment", "noop")
	})

	t.Run("noop", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_preprocessor_tasks_worker_import(t, h, s),
		})
		h.a.NoError(err)

		test_preprocessor_tasks_import_noop(ctx, t, h, svc, s, g, "import", "noop")
	})

	t.Run("experimentalExport", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_preprocessor_tasks_worker_export(t, h, s),
		})
		h.a.NoError(err)

		test_preprocessor_tasks_export_experimentalExport(ctx, t, h, svc, s, g, "export", "experimentalExport")
	})
	t.Run("noop", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: test_preprocessor_tasks_worker_export(t, h, s),
		})
		h.a.NoError(err)

		test_preprocessor_tasks_export_noop(ctx, t, h, svc, s, g, "export", "noop")
	})

}
