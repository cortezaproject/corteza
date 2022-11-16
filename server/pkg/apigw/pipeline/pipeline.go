package pipeline

import (
	"context"
	"net/http"
	"sort"
	"time"

	actx "github.com/cortezaproject/corteza/server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza/server/pkg/apigw/pipeline/chain"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"go.uber.org/zap"
)

type (
	Worker struct {
		Weight  int
		Async   bool
		Name    string
		Type    types.FilterKind
		Handler func(rw http.ResponseWriter, r *http.Request) error
	}

	workerSet []*Worker

	Pl struct {
		workers workerSet
		ch      chain.ChainHandler
		err     types.ErrorHandlerFunc
		log     *zap.Logger
		async   bool
	}
)

func NewPipeline(log *zap.Logger, ch chain.ChainHandler) *Pl {
	var (
		defaultErrorHandler = types.NewDefaultErrorHandler(log)
	)

	return &Pl{
		ch:  ch,
		log: log,
		err: defaultErrorHandler.Handler(),
	}
}

func (pp *Pl) Async(a bool) {
	pp.async = a
}

func (pp *Pl) Error() types.ErrorHandlerFunc {
	return pp.err
}

// ErrorHandler adds error handler
func (pp *Pl) ErrorHandler(ff types.ErrorHandlerFunc) {
	pp.err = ff
}

// Add filter
func (pp *Pl) Add(w *Worker) {
	pp.workers = append(pp.workers, w)
	sort.Sort(pp.workers)
}

// Handler is the main operating entry point for requests
// that handles filter groups
func (pp *Pl) Handler() http.Handler {
	// use the chi implementation of chains
	pp.ch.Chain(pp.makeMiddleware(pp.workers...))
	return pp.ch.Handler()
}

// makeMiddleware creates a list of handlers from workers
// it is used in chaining
func (pp *Pl) makeMiddleware(hh ...*Worker) (middleware []func(http.Handler) http.Handler) {
	for _, wrker := range hh {
		middleware = append(middleware, pp.makeHandler(*wrker))
	}

	return middleware
}

// makeHandler creates a handler from worker, it also
// wraps the error handling
func (pp *Pl) makeHandler(hh Worker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var (
				err   error
				start = time.Now()
				log   = pp.log.With(zap.String("filter", hh.Name))
			)

			log.Debug("started processing", zap.Bool("async", hh.Async))

			fn := func() (err error) {
				err = hh.Handler(rw, r)
				log.Debug("finished processing", zap.Duration("duration", time.Since(start)))
				return
			}

			if hh.Async {
				go func() {
					var (
						newCtx = context.Background()

						ident = auth.GetIdentityFromContext(r.Context())
						scope = actx.ScopeFromContext(r.Context())
					)

					newCtx = actx.ScopeToContext(context.Background(), scope)
					newCtx = auth.SetIdentityToContext(newCtx, ident)

					r = r.WithContext(newCtx)

					// only log error, do not call error handler,
					// since we do not reply back the response (it was already sent)
					if err = fn(); err != nil {
						log.Error(err.Error())
					}
				}()
			} else {
				err = fn()
			}

			if err != nil {
				pp.err(rw, r, err)
				return
			} else {
				next.ServeHTTP(rw, r)
			}
		})
	}
}

func (a workerSet) Len() int { return len(a) }
func (a workerSet) Less(i, j int) bool {
	return a[i].Weight < a[j].Weight
}
func (a workerSet) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a workerSet) Filter(f func(*Worker) (bool, error)) (out workerSet, err error) {
	var ok bool
	out = workerSet{}

	for i := range a {
		if ok, err = f(a[i]); err != nil {
			return
		} else if ok {
			out = append(out, a[i])
		}
	}

	return
}
