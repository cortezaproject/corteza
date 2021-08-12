package apigw

import (
	"context"
	"encoding/json"

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
		hasPostFilters    bool
		defaultPostFilter types.Handler
	)

	s.routes = route

	s.log.Debug("registering routes", zap.Int("count", len(s.routes)))

	defaultPostFilter, err := s.reg.Get("defaultJsonResponse")

	if err != nil {
		s.log.Error("could not register default filter", zap.Error(err))
	}

	for _, r := range s.routes {
		hasPostFilters = false
		log := s.log.With(zap.String("route", r.String()))

		r.pipe = pipeline.NewPipeline(log)
		r.opts = s.opts
		r.log = log

		regFilters, err := s.loadFilters(ctx, r.ID)

		if err != nil {
			log.Error("could not load functions for route", zap.Error(err))
			continue
		}

		r.pipe.ErrorHandler(filter.NewErrorHandler("error handler expediter", []string{}))

		for _, f := range regFilters {
			h, err := s.reg.Get(f.Ref)

			if err != nil {
				log.Error("could not register filter", zap.Error(err))
				continue
			}

			enc, err := json.Marshal(f.Params)

			if err != nil {
				log.Error("could not load params for filter", zap.String("ref", f.Ref), zap.Error(err))
				continue
			}

			h, err = s.reg.Merge(h, enc)

			if err != nil {
				log.Error("could not merge params to handler", zap.String("ref", f.Ref), zap.Error(err))
				continue
			}

			// check if it's a postfilter for async support
			if f.Kind == string(types.PostFilter) {
				hasPostFilters = true
			}

			r.pipe.Add(h)
		}

		// add default postfilter on async
		// routes if not present
		if r.meta.async && !hasPostFilters {
			log.Info("registering default postfilter", zap.Error(err))
			r.pipe.Add(defaultPostFilter)
		}

		log.Debug("successfuly registered route")
	}
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
