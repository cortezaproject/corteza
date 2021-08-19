package pipeline

import (
	"context"
	"sort"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"go.uber.org/zap"
)

type (
	Worker interface {
		types.Execer
		types.Stringer
		types.Sorter
	}

	workerSet []Worker

	workers struct {
		prefilter  workerSet
		processer  workerSet
		postfilter workerSet
	}

	Pl struct {
		w   workers
		err types.ErrorHandler
		log *zap.Logger
	}
)

func NewPipeline(log *zap.Logger) *Pl {
	return &Pl{
		log: log,
		w:   workers{},
		err: types.DefaultErrorHandler{},
	}
}

func (pp *Pl) Error() types.ErrorHandler {
	return pp.err
}

// Exec takes care of error handling and main
// functionality that takes place in worker
func (pp *Pl) Exec(ctx context.Context, scope *types.Scp, async bool) (err error) {
	err = pp.process(ctx, scope, pp.w.prefilter...)

	if err != nil {
		return
	}

	if async {
		go pp.process(ctx, scope, pp.w.processer...)
	} else {
		err = pp.process(ctx, scope, pp.w.processer...)

		if err != nil {
			return
		}
	}

	err = pp.process(ctx, scope, pp.w.postfilter...)

	if err != nil {
		return
	}

	return
}

// Add registers a new worker with parameters
// fetched from store
func (pp *Pl) Add(w Worker) {
	var pipe *workerSet

	switch w.Type() {
	case types.PreFilter:
		pipe = &pp.w.prefilter
	case types.Processer:
		pipe = &pp.w.processer
	case types.PostFilter:
		pipe = &pp.w.postfilter
	}

	*pipe = append(*pipe, w)
	sort.Sort(pipe)

	pp.log.Debug("registered worker", zap.Any("worker", w.String()))
}

// add error handler
func (pp *Pl) ErrorHandler(ff types.ErrorHandler) {
	pp.err = ff
}

func (pp *Pl) process(ctx context.Context, scope *types.Scp, w ...Worker) (err error) {
	for _, w := range w {
		pp.log.Debug("started worker", zap.Any("worker", w.String()))

		start := time.Now()
		err = w.Exec(ctx, scope)
		elapsed := time.Since(start)

		pp.log.Debug("finished worker", zap.Any("worker", w.String()), zap.Duration("duration", elapsed))

		if err != nil {
			pp.log.Debug("could not execute worker", zap.Error(err))
			return
		}
	}

	return
}

func (a workerSet) Len() int { return len(a) }
func (a workerSet) Less(i, j int) bool {
	return a[i].Weight() < a[j].Weight()
}
func (a workerSet) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
