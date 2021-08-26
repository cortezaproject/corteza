package apigw

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/cortezaproject/corteza-server/pkg/apigw/filter"
	"github.com/cortezaproject/corteza-server/pkg/apigw/filter/proxy"
	"github.com/cortezaproject/corteza-server/pkg/apigw/pipeline"
	"github.com/cortezaproject/corteza-server/pkg/apigw/registry"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	f "github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/options"
	st "github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type (
	storer interface {
		SearchApigwRoutes(ctx context.Context, f st.ApigwRouteFilter) (st.ApigwRouteSet, st.ApigwRouteFilter, error)
		SearchApigwFilters(ctx context.Context, f st.ApigwFilterFilter) (st.ApigwFilterSet, st.ApigwFilterFilter, error)
	}

	apigw struct {
		opts   *options.ApigwOpt
		log    *zap.Logger
		reg    *registry.Registry
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
	reg := registry.NewRegistry()
	reg.Preload()

	return &apigw{
		opts:   opts,
		log:    logger,
		storer: storer,
		reg:    reg,
		reload: make(chan bool),
	}
}

func (s *apigw) Reload(ctx context.Context) {
	go func() {
		s.reload <- true
	}()
}

func (s *apigw) loadRoutes(ctx context.Context) (rr []*route, err error) {
	routes, _, err := s.storer.SearchApigwRoutes(ctx, st.ApigwRouteFilter{Enabled: true, Deleted: f.StateExcluded})

	if err != nil {
		return
	}

	for _, r := range routes {
		route := &route{
			ID:       r.ID,
			endpoint: r.Endpoint,
			method:   r.Method,
			meta: routeMeta{
				debug: r.Meta.Debug,
				async: r.Meta.Async,
			},
		}

		rr = append(rr, route)
	}

	return
}

func (s *apigw) loadFilters(ctx context.Context, route uint64) (ff []*st.ApigwFilter, err error) {
	ff, _, err = s.storer.SearchApigwFilters(ctx, st.ApigwFilterFilter{RouteID: route})
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
	var (
		defaultPostFilter types.Handler
	)

	s.routes = route

	s.log.Debug("registering routes", zap.Int("count", len(s.routes)))

	defaultPostFilter, err := s.reg.Get("defaultJsonResponse")

	if err != nil {
		s.log.Error("could not register default filter", zap.Error(err))
	}

	for _, r := range s.routes {
		var (
			log  = s.log.With(zap.String("route", r.String()))
			pipe = pipeline.NewPipeline(log)
		)

		// pipeline needs to know how to handle
		// async processers
		pipe.Async(r.meta.async)

		r.opts = s.opts
		r.log = log

		regFilters, err := s.loadFilters(ctx, r.ID)

		if err != nil {
			log.Error("could not load filters for route", zap.Error(err))
			continue
		}

		for _, f := range regFilters {
			flog := log.With(zap.String("ref", f.Ref))

			// make sure there is only one postfilter
			// on async routes
			if r.meta.async && f.Kind == string(types.PostFilter) {
				flog.Debug("not registering filter for async route")
				continue
			}

			ff, err := s.registerFilter(f, r)

			if err != nil {
				flog.Error("could not register filter", zap.Error(err))
				continue
			}

			pipe.Add(ff)

			flog.Debug("registered filter")
		}

		// add default postfilter on async
		// routes if not present
		if r.meta.async {
			log.Info("registering default postfilter", zap.Error(err))

			pipe.Add(&pipeline.Worker{
				Handler: defaultPostFilter.Handler(),
				Name:    defaultPostFilter.String(),
				Type:    types.PostFilter,
				Weight:  math.MaxInt8,
			})
		}

		r.handler = pipe.Handler()
		r.errHandler = pipe.Error()

		log.Debug("successfuly registered route")
	}
}

func (s *apigw) registerFilter(f *st.ApigwFilter, r *route) (ff *pipeline.Worker, err error) {
	handler, err := s.reg.Get(f.Ref)

	if err != nil {
		return
	}

	enc, err := json.Marshal(f.Params)

	if err != nil {
		err = fmt.Errorf("could not load params for filter: %s", err)
		return
	}

	handler, err = s.reg.Merge(handler, enc)

	if err != nil {
		err = fmt.Errorf("could not merge params to handler: %s", err)
		return
	}

	ff = &pipeline.Worker{
		Async:   r.meta.async && f.Kind == string(types.Processer),
		Handler: handler.Handler(),
		Name:    handler.String(),
		Type:    types.FilterKind(f.Kind),
		Weight:  filter.FilterWeight(int(f.Weight), types.FilterKind(f.Kind)),
	}

	return
}

func (s *apigw) Funcs(kind string) (list types.FilterMetaList) {
	list = s.reg.All()

	if kind != "" {
		list, _ = list.Filter(func(fm *types.FilterMeta) (bool, error) {
			return string(fm.Kind) == kind, nil
		})
	}

	return
}

func (s *apigw) ProxyAuthDef() (list []*proxy.ProxyAuthDefinition) {
	list = proxy.ProxyAuthDef()
	return
}
