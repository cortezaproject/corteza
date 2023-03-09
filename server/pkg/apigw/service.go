package apigw

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/cortezaproject/corteza/server/pkg/apigw/filter"
	"github.com/cortezaproject/corteza/server/pkg/apigw/filter/proxy"
	"github.com/cortezaproject/corteza/server/pkg/apigw/pipeline"
	"github.com/cortezaproject/corteza/server/pkg/apigw/pipeline/chain"
	"github.com/cortezaproject/corteza/server/pkg/apigw/profiler"
	"github.com/cortezaproject/corteza/server/pkg/apigw/registry"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	routeServicer interface {
		LoadRoute(context.Context, string, string) ([]*types.Route, error)
		LoadRoutes(context.Context) ([]*types.Route, error)
	}

	filterServicer interface {
		LoadFilters(context.Context, uint64) ([]*types.RouteFilter, error)
	}

	apigw struct {
		log    *zap.Logger
		reg    *registry.Registry
		routes []*route
		mx     *chi.Mux
		pr     *profiler.Profiler

		rs routeServicer
		fs filterServicer

		cfg types.Config
	}
)

var (
	// global service
	apiGw *apigw
)

func Service() *apigw {
	return apiGw
}

// Setup handles the singleton service
func Setup(cfg types.Config, log *zap.Logger, rs routeServicer, fs filterServicer) {
	if apiGw != nil {
		return
	}

	apiGw = New(cfg, log, rs, fs)
}

func New(cfg types.Config, logger *zap.Logger, rs routeServicer, fs filterServicer) *apigw {
	var (
		pr = profiler.New()
	)

	reg := registry.NewRegistry(cfg)
	reg.Preload()

	return &apigw{
		log: logger.Named("http.apigw"),
		rs:  rs,
		fs:  fs,
		reg: reg,
		pr:  pr,
		cfg: cfg,
	}
}

// ServeHTTP forwards the given HTTP request to the underlying chi mux which
// then handles the heavy lifting
//
// When reloading routes, make sure to replace the original mux
func (s *apigw) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.mx == nil {
		http.Error(w, "Integration Gateway not initialized", http.StatusInternalServerError)
		return
	}

	if len(s.routes) == 0 {
		helperDefaultResponse(s.cfg, s.pr, s.log)(w, r)
		return
	}

	// Remove route context for chi
	//
	// Without this, chi can not properly handle requests
	// in API gateway's sub-router
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, nil))

	// Handle api-gw request
	s.mx.ServeHTTP(w, r)
}

// Reload reloads all routes and their filters
//
// The procedure constructs a new chi mux
func (s *apigw) Reload(ctx context.Context) (err error) {
	routes, err := s.loadRoutes(ctx)

	if err != nil {
		s.log.Error("could not reload Integration Gateway routes", zap.Error(err))
		return
	}

	s.Init(ctx, routes...)

	// Rebuild the mux
	s.mx = chi.NewMux()

	for _, r := range s.routes {
		// Register route handler on endpoint & method
		s.mx.Method(r.Method, r.Endpoint, r)
	}

	// handling missed hits
	// profiler gets the missed hit info also
	{
		var (
			defaultMethodResponse = helperMethodNotAllowed(s.cfg, s.pr, s.log)
			defaultResponse       = helperDefaultResponse(s.cfg, s.pr, s.log)
		)

		s.mx.NotFound(defaultResponse)
		s.mx.MethodNotAllowed(defaultMethodResponse)
	}

	return nil
}

// ReloadEndpoint reload a route and its filters
//
// The procedure use existing chi mux
func (s *apigw) ReloadEndpoint(ctx context.Context, method, endpoint string) (err error) {
	var (
		routes []*route
	)

	routes, err = s.loadRoute(ctx, method, endpoint)

	if err != nil {
		s.log.Error("could not reload Integration Gateway routes", zap.Error(err))
		return
	}

	// rr := append(s.routes, routes...)
	s.PrepRoutes(ctx, routes...)

	if s.mx == nil {
		// Rebuild the mux
		s.mx = chi.NewMux()
	}

	for _, r := range routes {
		// Register route handler on endpoint & method
		s.mx.Method(r.Method, r.Endpoint, r)
	}

	// Make sure to append newly registered routes
	s.AppendRoutes(routes...)

	// handling missed hits
	// profiler gets the missed hit info also
	{
		var (
			defaultMethodResponse = helperMethodNotAllowed(s.cfg, s.pr, s.log)
			defaultResponse       = helperDefaultResponse(s.cfg, s.pr, s.log)
		)

		s.mx.NotFound(defaultResponse)
		s.mx.MethodNotAllowed(defaultMethodResponse)
	}

	return nil
}

// Init all routes
func (s *apigw) Init(ctx context.Context, routes ...*route) {
	s.PrepRoutes(ctx, routes...)
	s.routes = routes
}

func (s *apigw) PrepRoutes(ctx context.Context, routes ...*route) {
	var (
		err               error
		defaultPostFilter types.Handler
	)

	s.loadInfo()
	s.log.Debug("preparing routes", zap.Int("count", len(routes)))

	defaultPostFilter, err = s.reg.Get("defaultJsonResponse")

	if err != nil {
		s.log.Error("could not register default filter", zap.Error(err))
	}

	for _, r := range routes {
		var (
			log  = s.log.With(zap.String("route", r.String()))
			pipe = pipeline.NewPipeline(log, chain.NewDefault())

			regFilters []*types.RouteFilter
		)

		// pipeline needs to know how to handle
		// async processers
		pipe.Async(r.Meta.Async)

		r.cfg = s.cfg
		r.log = log
		r.pr = s.pr

		regFilters, err = s.fs.LoadFilters(ctx, r.ID)

		if err != nil {
			log.Error("could not load filters for route", zap.Error(err))
			continue
		}

		for _, rf := range regFilters {
			flog := log.With(zap.String("ref", rf.Ref))

			// make sure there is only one postfilter
			// on async routes
			if r.Meta.Async && rf.Kind == string(types.PostFilter) {
				flog.Debug("not registering filter for async route")
				continue
			}

			var ff *pipeline.Worker
			ff, err = s.registerFilter(rf, r)
			if err != nil {
				flog.Error("could not register filter", zap.Error(err))
				continue
			}

			pipe.Add(ff)

			flog.Debug("registered filter")
		}

		// add default postfilter on async
		// routes if not present
		if r.Meta.Async {
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

		log.Debug("successfully registered route")
	}
}

func (s *apigw) AppendRoutes(routes ...*route) {
	var (
		rMap = make(map[string]*route)
		uniq = func(r *route) string {
			if routes == nil {
				return ""
			}
			return r.Method + r.Endpoint
		}
	)

	for _, r := range routes {
		if r == nil {
			continue
		}
		rMap[uniq(r)] = r
	}

	// update existing routes
	for i, r := range s.routes {
		if val, ok := rMap[uniq(r)]; ok && val != nil {
			s.routes[i] = val
			rMap[uniq(r)] = nil
		}
	}

	// add new routes
	for _, r := range rMap {
		if r == nil {
			continue
		}
		s.routes = append(s.routes, r)
	}
}

func (s *apigw) NotFound(_ context.Context, method, endpoint string) {
	if s.mx == nil || len(method) == 0 || len(endpoint) == 0 {
		return
	}

	var (
		defaultResponse = helperDefaultResponse(s.cfg, s.pr, s.log)
	)

	// Attach 404 handler
	s.mx.Method(method, endpoint, defaultResponse)
}

func (s *apigw) registerFilter(f *types.RouteFilter, r *route) (ff *pipeline.Worker, err error) {
	handler, err := s.reg.Get(f.Ref)

	if err != nil {
		return
	}

	enc, err := json.Marshal(f.Params)

	if err != nil {
		err = fmt.Errorf("could not load params for filter: %s", err)
		return
	}

	handler, err = s.reg.Merge(handler, enc, s.cfg)

	if err != nil {
		err = fmt.Errorf("could not merge params to handler: %s", err)
		return
	}

	ff = &pipeline.Worker{
		Async:   r.Meta.Async && f.Kind == string(types.Processer),
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

func (s *apigw) UpdateSettings(ctx context.Context, cfg types.Config) {
	s.cfg = cfg

	s.reg = registry.NewRegistry(cfg)
	s.reg.Preload()

	s.Reload(ctx)
}

func (s *apigw) loadRoutes(ctx context.Context) (rr []*route, err error) {
	routes, err := s.rs.LoadRoutes(ctx)

	if err != nil {
		return
	}

	rr = mapRoute(routes)

	return
}

func (s *apigw) loadRoute(ctx context.Context, method, endpoint string) (rr []*route, err error) {
	routes, err := s.rs.LoadRoute(ctx, method, endpoint)

	if err != nil {
		return
	}

	rr = mapRoute(routes)

	return
}

func mapRoute(routes []*types.Route) (rr []*route) {
	for _, r := range routes {
		route := &route{}

		route.ID = r.ID
		route.Endpoint = r.Endpoint
		route.Method = r.Method
		route.Meta = types.RouteMeta{
			Debug: r.Meta.Debug,
			Async: r.Meta.Async,
		}

		rr = append(rr, route)
	}

	return
}

func (s *apigw) loadInfo() {
	s.log.Info("loading Integration Gateway")

	if s.cfg.Profiler.Enabled {
		if !s.cfg.Profiler.Global {
			s.log.Warn("profiler enabled only for routes with a profiler prefilter, use global setting to enable for all (APIGW_PROFILER_GLOBAL)")
		}
	} else {
		if s.cfg.Profiler.Global {
			s.log.Warn("profiler global is enabled, but profiler disabled, no routes will be profiled",
				zap.Bool("APIGW_PROFILER_ENABLED", s.cfg.Profiler.Enabled),
				zap.Bool("APIGW_PROFILER_GLOBAL", s.cfg.Profiler.Global))
		}
	}
}

func (s *apigw) Profiler() *profiler.Profiler {
	return s.pr
}
