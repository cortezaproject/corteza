package apigw

import (
	"context"

	"go.uber.org/zap"
)

type (
	Execer interface {
		Exec(context.Context, *scp) error
	}

	Sorter interface {
		Weight() int
	}

	ErrorHandler interface {
		Exec(context.Context, *scp, error)
	}

	Worker interface {
		Execer
		// Sorter
		Stringer
	}

	Stringer interface {
		String() string
	}

	workers []Worker

	pl struct {
		w   workers
		err ErrorHandler
		log *zap.Logger
	}
)

func NewPipeline(log *zap.Logger, w ...Worker) *pl {
	return &pl{
		w:   w,
		log: log,
		err: defaultErrorHandler{},
	}
}

// Exec takes care of error handling and main
// functionality that takes place in worker
func (pp *pl) Exec(ctx context.Context, scope *scp) (err error) {
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
func (pp *pl) Add(w Worker) {
	pp.w = append(pp.w, w)
	// sort.Sort(pp.w)

	pp.log.Debug("registered worker", zap.Any("worker", w.String()))
}

// add error handler
func (pp *pl) ErrorHandler(ff ErrorHandler) {
	pp.err = ff
}

// func (a workers) Len() int { return len(a) }
// func (a workers) Less(i, j int) bool {
// 	return a[i].worker.Weight() < a[j].worker.Weight()
// }
// func (a workers) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
