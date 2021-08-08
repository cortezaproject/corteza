package pipeline

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"go.uber.org/zap"
)

type (
	Worker interface {
		types.Execer
		types.Stringer
		// Sorter
	}

	workers []Worker

	Pl struct {
		w   workers
		err types.ErrorHandler
		log *zap.Logger
	}
)

func NewPipeline(log *zap.Logger, w ...Worker) *Pl {
	return &Pl{
		w:   w,
		log: log,
		err: types.DefaultErrorHandler{},
	}
}

func (pp *Pl) Error() types.ErrorHandler {
	return pp.err
}

// Exec takes care of error handling and main
// functionality that takes place in worker
func (pp *Pl) Exec(ctx context.Context, scope *types.Scp) (err error) {
	for _, w := range pp.w {

		pp.log.Debug("executing worker", zap.Any("worker", w.String()))
		err = w.Exec(ctx, scope)

		if err != nil {
			pp.log.Debug("could not execute worker", zap.Error(err))
			return
		}
	}

	return
}

// Add registers a new worker with parameters
// fetched from store
func (pp *Pl) Add(w Worker) {
	pp.w = append(pp.w, w)
	// sort.Sort(pp.w)

	pp.log.Debug("registered worker", zap.Any("worker", w.String()))
}

// add error handler
func (pp *Pl) ErrorHandler(ff types.ErrorHandler) {
	pp.err = ff
}

// func (a workers) Len() int { return len(a) }
// func (a workers) Less(i, j int) bool {
// 	return a[i].worker.Weight() < a[j].worker.Weight()
// }
// func (a workers) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
