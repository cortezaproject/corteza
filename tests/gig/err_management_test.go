package gig

import (
	"context"
	"errors"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

type (
	workerStage int

	workerFail struct {
		err   error
		stage workerStage
	}
)

const (
	workerStagePrepare workerStage = iota
	workerStagePreprocess
	workerStageExec
	workerStageState
	workerStageCleanup
)

func Test_err_management(t *testing.T) {
	var (
		ctx, svc, h, _ = setup(t)
		err            error
		g              gig.Gig
		expErr         = errors.New("failed: reasons")
	)

	g, err = svc.Create(ctx, gig.UpdatePayload{
		Worker: &workerFail{
			err:   expErr,
			stage: workerStagePrepare,
		},
	})
	h.a.NoError(err)

	g, err = svc.Prepare(ctx, g)
	h.a.ErrorIs(err, expErr)
}

// gig.Worker interface definitions

func (w *workerFail) Prepare(ctx context.Context, sources ...gig.Source) error {
	if w.stage == workerStagePrepare {
		return w.err
	}

	return nil
}

func (w *workerFail) Preprocess(ctx context.Context, pp ...gig.Preprocessor) error {
	if w.stage == workerStagePreprocess {
		return w.err
	}

	return nil
}

func (w *workerFail) Exec(ctx context.Context) (out gig.SourceSet, meta gig.WorkMeta, err error) {
	if w.stage == workerStageExec {
		return nil, nil, w.err
	}
	return

}

func (w *workerFail) State(ctx context.Context) (state gig.WorkerState, err error) {
	if w.stage == workerStageState {
		return nil, w.err
	}
	return

}

func (w *workerFail) Cleanup(ctx context.Context) error {
	if w.stage == workerStageCleanup {
		return w.err
	}
	return nil

}

func (w *workerFail) Ref() string {
	return "testing"
}
