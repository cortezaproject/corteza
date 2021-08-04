package apigw

import (
	"context"
	"encoding/json"

	as "github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type (
	storer interface {
		SearchApigwRoutes(ctx context.Context, f types.ApigwRouteFilter) (types.ApigwRouteSet, types.ApigwRouteFilter, error)
		SearchApigwFunctions(ctx context.Context, f types.ApigwFunctionFilter) (types.ApigwFunctionSet, types.ApigwFunctionFilter, error)
	}

	apigw struct {
		opts   *options.ApigwOpt
		log    *zap.Logger
		reg    *registry
		routes []*route
		storer storer
		reload chan bool
	}
)

var (
	// global service
	apiGw *apigw
)

func Service() *apigw {
	return apiGw
}

func Set(a *apigw) {
	apiGw = a
}

// Setup handles the singleton service
func Setup(opts *options.ApigwOpt, log *zap.Logger, storer storer) {
	if apiGw != nil {
		return
	}

	apiGw = New(opts, log, storer)
}

func New(opts *options.ApigwOpt, logger *zap.Logger, storer storer) *apigw {
	reg := NewRegistry()
	reg.Preload()

	return &apigw{
		opts:   opts,
		log:    logger,
		storer: storer,
		reload: make(chan bool),
		reg:    reg,
	}
}

func (s *apigw) Reload(ctx context.Context) {
	go func() {
		s.reload <- true
	}()
}

func (s *apigw) loadRoutes(ctx context.Context) (rr []*route, err error) {
	routes, _, err := s.storer.SearchApigwRoutes(ctx, types.ApigwRouteFilter{Enabled: true, Deleted: filter.StateExcluded})

	if err != nil {
		return
	}

	for _, r := range routes {
		route := &route{
			ID:       r.ID,
			endpoint: r.Endpoint,
			method:   r.Method,
		}

		rr = append(rr, route)
	}

	return
}

func (s *apigw) loadFunctions(ctx context.Context, route uint64) (ff []*types.ApigwFunction, err error) {
	ff, _, err = s.storer.SearchApigwFunctions(ctx, types.ApigwFunctionFilter{RouteID: route})
	return
}

func (s *apigw) Router(r chi.Router) {
	var (
		ctx = context.Background()
	)

	r.HandleFunc("/", helperDefaultResponse(s.opts))

	routes, err := s.loadRoutes(ctx)

	if err != nil {
		s.log.Error("could not load routes", zap.Error(err))
		return
	}

	s.Init(ctx, routes...)

	for _, route := range s.routes {
		r.Handle(route.endpoint, route)
	}

	go func() {
		for {
			select {
			case <-s.reload:
				routes, err := s.loadRoutes(ctx)

				if err != nil {
					s.log.Error("could not reload API Gateway routes", zap.Error(err))
					return
				}

				s.log.Debug("reloading API Gateway routes and functions", zap.Int("count", len(routes)))

				s.Init(ctx, routes...)

				for _, route := range s.routes {
					r.Handle(route.endpoint, route)
				}

			case <-ctx.Done():
				s.log.Debug("shutting down API Gateway service")
				return
			}
		}
	}()
}

// init all the routes
func (s *apigw) Init(ctx context.Context, route ...*route) {
	s.routes = route

	s.log.Debug("registering routes", zap.Int("count", len(s.routes)))

	for _, r := range s.routes {
		log := s.log.With(zap.String("route", r.String()))

		r.pipe = NewPipeline(log)

		r.opts = s.opts
		r.log = log

		regFuncs, err := s.loadFunctions(ctx, r.ID)

		if err != nil {
			log.Error("could not load functions for route", zap.Error(err))
			continue
		}

		r.pipe.ErrorHandler(errorHandler{
			name:   "error handler expediter",
			args:   []string{},
			weight: 5,
			step:   0,
		})

		for _, f := range regFuncs {
			h, err := s.reg.Get(f.Ref)

			if err != nil {
				log.Error("could not register function for route", zap.Error(err))
				continue
			}

			enc, err := json.Marshal(f.Params)

			if err != nil {
				log.Error("could not load params for function", zap.String("ref", f.Ref), zap.Error(err))
				continue
			}

			h, err = s.reg.Merge(h, enc)

			if err != nil {
				log.Error("could not merge params to handler", zap.String("ref", f.Ref), zap.Error(err))
				continue
			}

			r.pipe.Add(h)
		}

		log.Debug("successfuly registered route")
	}
}

func (s *apigw) Funcs(kind string) (list functionMetaList) {
	list = s.reg.All()

	if kind != "" {
		list, _ = list.Filter(func(fm *functionMeta) (bool, error) {
			return string(fm.Kind) == kind, nil
		})
	}

	return
}

func (s *apigw) ProxyAuthDef() (list []*proxyAuthDefinition) {
	list = ProxyAuthDef()
	return
}

func NewWorkflow() (wf WfExecer) {
	return as.Workflow(&zap.Logger{}, options.CorredorOpt{})
}
