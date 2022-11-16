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
	f "github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/options"
	st "github.com/cortezaproject/corteza/server/system/types"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	storer interface {
		SearchApigwRoutes(ctx context.Context, f st.ApigwRouteFilter) (st.ApigwRouteSet, st.ApigwRouteFilter, error)
		SearchApigwFilters(ctx context.Context, f st.ApigwFilterFilter) (st.ApigwFilterSet, st.ApigwFilterFilter, error)
	}

	apigw struct {
		opts   options.ApigwOpt
		log    *zap.Logger
		reg    *registry.Registry
		routes []*route
		mx     *chi.Mux
		pr     *profiler.Profiler
		storer storer
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
func Setup(opts options.ApigwOpt, log *zap.Logger, storer storer) {
	if apiGw != nil {
		return
	}

	apiGw = New(opts, log, storer)
}

func New(opts options.ApigwOpt, logger *zap.Logger, storer storer) *apigw {
	var (
		pr  = profiler.New()
		reg = registry.NewRegistry(opts)
	)

	reg.Preload()

	return &apigw{
		opts:   opts,
		log:    logger.Named("http.apigw"),
		storer: storer,
		reg:    reg,
		pr:     pr,
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
		helperDefaultResponse(s.opts, s.pr, s.log)(w, r)
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

// Reload reloads routes and their filters
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
		s.mx.Method(r.method, r.endpoint, r)
	}

	// handling missed hits
	// profiler gets the missed hit info also
	{
		var (
			defaultMethodResponse = helperMethodNotAllowed(s.opts, s.pr, s.log)
			defaultResponse       = helperDefaultResponse(s.opts, s.pr, s.log)
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
		s.mx.Method(r.method, r.endpoint, r)
	}

	// Make sure to append newly registered routes
	s.AppendRoutes(routes...)

	// handling missed hits
	// profiler gets the missed hit info also
	{
		var (
			defaultMethodResponse = helperMethodNotAllowed(s.opts, s.pr, s.log)
			defaultResponse       = helperDefaultResponse(s.opts, s.pr, s.log)
		)

		s.mx.NotFound(defaultResponse)
		s.mx.MethodNotAllowed(defaultMethodResponse)
	}

	return nil
}

// Init all routes
func (s *apigw) Init(ctx context.Context, routes ...*route) {
	var (
		defaultPostFilter types.Handler
	)

	s.routes = routes

	s.loadInfo()
	s.log.Debug("registering routes", zap.Int("count", len(s.routes)))

	defaultPostFilter, err := s.reg.Get("defaultJsonResponse")

	if err != nil {
		s.log.Error("could not register default filter", zap.Error(err))
	}

	for _, r := range s.routes {
		var (
			log  = s.log.With(zap.String("route", r.String()))
			pipe = pipeline.NewPipeline(log, chain.NewDefault())
		)

		// pipeline needs to know how to handle
		// async processers
		pipe.Async(r.meta.async)

		r.opts = s.opts
		r.log = log
		r.pr = s.pr

		regFilters, err := s.loadFilters(ctx, r.ID)

		if err != nil {
			log.Error("could not load filters for route", zap.Error(err))
			continue
		}

		for _, rf := range regFilters {
			flog := log.With(zap.String("ref", rf.Ref))

			// make sure there is only one postfilter
			// on async routes
			if r.meta.async && rf.Kind == string(types.PostFilter) {
				flog.Debug("not registering filter for async route")
				continue
			}

			ff, err := s.registerFilter(rf, r)

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

		log.Debug("successfully registered route")
	}
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

			regFilters []*st.ApigwFilter
		)

		// pipeline needs to know how to handle
		// async processers
		pipe.Async(r.meta.async)

		r.opts = s.opts
		r.log = log
		r.pr = s.pr

		regFilters, err = s.loadFilters(ctx, r.ID)
		if err != nil {
			log.Error("could not load filters for route", zap.Error(err))
			continue
		}

		for _, rf := range regFilters {
			flog := log.With(zap.String("ref", rf.Ref))

			// make sure there is only one postfilter
			// on async routes
			if r.meta.async && rf.Kind == string(types.PostFilter) {
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
			return r.method + r.endpoint
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

	return
}

func (s *apigw) NotFound(_ context.Context, method, endpoint string) {
	if s.mx == nil || len(method) == 0 || len(endpoint) == 0 {
		return
	}

	var (
		defaultResponse = helperDefaultResponse(s.opts, s.pr, s.log)
	)

	// Attach 404 handler
	s.mx.Method(method, endpoint, defaultResponse)
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

func (s *apigw) loadRoutes(ctx context.Context) (rr []*route, err error) {
	var (
		routes st.ApigwRouteSet
		agwf   = st.ApigwRouteFilter{
			Deleted:  f.StateExcluded,
			Disabled: f.StateExcluded,
		}
	)

	routes, _, err = s.storer.SearchApigwRoutes(ctx, agwf)
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

func (s *apigw) loadRoute(ctx context.Context, method, endpoint string) (rr []*route, err error) {
	var (
		routes st.ApigwRouteSet
		agwf   = st.ApigwRouteFilter{
			Endpoint: endpoint,
			Method:   method,
			Deleted:  f.StateExcluded,
			Disabled: f.StateExcluded,
		}
	)

	routes, _, err = s.storer.SearchApigwRoutes(ctx, agwf)
	if err != nil {
		return
	}

	for _, r := range routes {
		rr = append(rr, &route{
			ID:       r.ID,
			endpoint: r.Endpoint,
			method:   r.Method,
			meta: routeMeta{
				debug: r.Meta.Debug,
				async: r.Meta.Async,
			},
		})
	}

	return
}

func (s *apigw) loadFilters(ctx context.Context, route uint64) (ff []*st.ApigwFilter, err error) {
	ff, _, err = s.storer.SearchApigwFilters(ctx, st.ApigwFilterFilter{
		RouteID:  route,
		Deleted:  f.StateExcluded,
		Disabled: f.StateExcluded,
	})

	return
}

func (s *apigw) loadInfo() {
	s.log.Info("loading Integration Gateway", zap.Bool("debug", s.opts.Debug), zap.Bool("log", s.opts.LogEnabled))

	if s.opts.ProfilerEnabled {
		if s.opts.LogRequestBody {
			s.log.Info("profiler and request body logging is enabled, profiler use is prefered",
				zap.Bool("APIGW_PROFILER_ENABLED", s.opts.ProfilerEnabled),
				zap.Bool("APIGW_LOG_REQUEST_BODY", s.opts.LogRequestBody))
		} else {
			s.log.Info("request body logging is enabled, profiler use is prefered (APIGW_PROFILER_ENABLED)",
				zap.Bool("APIGW_LOG_REQUEST_BODY", s.opts.LogRequestBody))
		}

		if !s.opts.ProfilerGlobal {
			s.log.Warn("profiler enabled only for routes with a profiler prefilter, use global setting to enable for all (APIGW_PROFILER_GLOBAL)")
		}
	} else {
		if s.opts.ProfilerGlobal {
			s.log.Warn("profiler global is enabled, but profiler disabled, no routes will be profiled",
				zap.Bool("APIGW_PROFILER_ENABLED", s.opts.ProfilerEnabled),
				zap.Bool("APIGW_PROFILER_GLOBAL", s.opts.ProfilerGlobal))
		}
	}
}

func (s *apigw) Profiler() *profiler.Profiler {
	return s.pr
}
