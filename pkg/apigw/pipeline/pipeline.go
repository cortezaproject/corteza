package pipeline

import (
	"net/http"
	"sort"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type (
	Worker struct {
		Handler func(rw http.ResponseWriter, r *http.Request) error
		Weight  int
		Name    string
	}

	workerSet []*Worker

	Pl struct {
		workers workerSet
		err     types.ErrorHandlerFunc
		log     *zap.Logger
	}
)

func NewPipeline(log *zap.Logger) *Pl {
	var (
		defaultErrorHandler = types.NewDefaultErrorHandler(log)
	)

	return &Pl{
		log: log,
		err: defaultErrorHandler.Handler(),
	}
}

func (pp *Pl) Error() types.ErrorHandlerFunc {
	return pp.err
}

// add error handler
func (pp *Pl) ErrorHandler(ff types.ErrorHandlerFunc) {
	pp.err = ff
}

// add filter
func (pp *Pl) Add(w *Worker) {
	pp.workers = append(pp.workers, w)
	sort.Sort(pp.workers)
}

func (pp *Pl) AddHandler(h http.Handler) {}

func (pp *Pl) Handler() http.Handler {
	var (
		middleware []func(http.Handler) http.Handler
	)

	for _, wrker := range pp.workers {
		middleware = append(middleware, pp.makeHandler(*wrker))
	}

	return chi.Chain(middleware...).Handler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}))
}

func (pp *Pl) makeHandler(hh Worker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var (
				start = time.Now()
			)

			pp.log.Debug("started processing", zap.String("filter", hh.Name))

			// if w.async {
			// 	ctx = context.Background()
			// 	r.WithContext(context.Background())
			// 	go w.handler(rw, r)
			// 	next.ServeHTTP(rw, r)
			// } else {

			err := hh.Handler(rw, r)

			pp.log.Debug("finished processing",
				zap.String("filter", hh.Name),
				zap.Duration("duration", time.Since(start)))

			if err != nil {
				pp.err(rw, r, err)
				return
			} else {
				next.ServeHTTP(rw, r)
			}
			// }

		})
	}
}

func (a workerSet) Len() int { return len(a) }
func (a workerSet) Less(i, j int) bool {
	return a[i].Weight < a[j].Weight
}
func (a workerSet) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
