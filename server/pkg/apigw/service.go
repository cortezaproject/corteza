package apigw

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"

	"github.com/crusttech/human/server/pkg/apigw/filter"
	"github.com/crusttech/human/server/pkg/apigw/filter/proxy"
	"github.com/crusttech/human/server/pkg/apigw/pipeline"
	"github.com/crusttech/human/server/pkg/apigw/pipeline/chain"
	"github.com/crusttech/human/server/pkg/apigw/profiler"
	"github.com/crusttech/human/server/pkg/apigw/registry"
	"github.com/crusttech/human/server/pkg/apigw/types"
	f "github.com/crusttech/human/server/pkg/filter"
	st "github.com/crusttech/human/server/system/types"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	storer interface {
		SearchApigwRoutes(ctx context.Context, f st.ApigwRouteFilter) (st.ApigwRouteSet, st.ApigwRouteFilter, error)
		SearchApigwFilters(ctx context.Context, f st.ApigwFilterFilter) (st.ApigwFilterSet, st.ApigwFilterFilter, error)
	}

	apigw struct {
		log    *zap.Logger
		reg    *registry.Registry
		routes []*route
		mx     *chi.Mux
		pr     *profiler.Profiler
		storer storer

		cfg types.Config

		// reloadMu serializes every structural change (Reload, ReloadEndpoint,
		// NotFound, UpdateSettings) so two goroutines never mutate the chi tree
		// at once. mu guards the mx/routes fields for the ServeHTTP read side
		// against the swap.
		reloadMu sync.Mutex
		mu       sync.RWMutex
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
func Setup(cfg types.Config, log *zap.Logger, storer storer) {
	if apiGw != nil {
		return
	}

	apiGw = New(cfg, log, storer)
}

func New(cfg types.Config, logger *zap.Logger, storer storer) *apigw {
	var (
		pr = profiler.New()
	)

	reg := registry.NewRegistry(cfg)
	reg.Preload()

	return &apigw{
		log:    logger.Named("http.apigw"),
		storer: storer,
		reg:    reg,
		pr:     pr,
		cfg:    cfg,
	}
}

// ServeHTTP forwards the given HTTP request to the underlying chi mux which
// then handles the heavy lifting
//
// When reloading routes, make sure to replace the original mux
func (s *apigw) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Capture the current mux; a reload only ever swaps in a fresh, immutable
	// mux, so the captured pointer stays safe to serve even mid-reload.
	s.mu.RLock()
	mx := s.mx
	empty := len(s.routes) == 0
	s.mu.RUnlock()

	if mx == nil {
		http.Error(w, "Integration Gateway not initialized", http.StatusInternalServerError)
		return
	}

	if empty {
		helperDefaultResponse(s.cfg, s.pr, s.log)(w, r)
		return
	}

	// Remove route context for chi
	//
	// Without this, chi can not properly handle requests
	// in API gateway's sub-router
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, nil))

	// Handle api-gw request
	mx.ServeHTTP(w, r)
}

// Reload reloads all routes and their filters
//
// The procedure constructs a new chi mux
func (s *apigw) Reload(ctx context.Context) error {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()
	return s.reloadLocked(ctx)
}

// reloadLocked rebuilds every route into a fresh mux and publishes it. The
// caller must hold reloadMu.
func (s *apigw) reloadLocked(ctx context.Context) error {
	routes, err := s.loadRoutes(ctx)
	if err != nil {
		s.log.Error("could not reload Integration Gateway routes", zap.Error(err))
		return err
	}

	s.PrepRoutes(ctx, routes...)
	s.publish(s.buildMux(routes), routes)
	return nil
}

// buildMux constructs a fresh chi mux from the given routes plus the miss
// handlers. The returned mux is never mutated afterwards, so ServeHTTP can hold
// a captured reference to it safely.
func (s *apigw) buildMux(routes []*route) *chi.Mux {
	mx := chi.NewMux()
	for _, r := range routes {
		mx.Method(r.method, r.endpoint, r)
	}
	mx.NotFound(helperDefaultResponse(s.cfg, s.pr, s.log))
	mx.MethodNotAllowed(helperMethodNotAllowed(s.cfg, s.pr, s.log))
	return mx
}

// publish swaps in a freshly built mux and its route set under the field lock.
func (s *apigw) publish(mx *chi.Mux, routes []*route) {
	s.mu.Lock()
	s.mx = mx
	s.routes = routes
	s.mu.Unlock()
}

// currentRoutes returns a copy of the live route set for a caller that is about
// to build the next one. Caller must hold reloadMu so the set cannot change
// between the read and the publish.
func (s *apigw) currentRoutes() []*route {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*route, len(s.routes))
	copy(out, s.routes)
	return out
}

// mergeRoutes returns existing with incoming applied: a route replaces one with
// the same method+endpoint, otherwise it is appended.
func mergeRoutes(existing, incoming []*route) []*route {
	out := make([]*route, len(existing))
	copy(out, existing)

	idx := make(map[string]int, len(out))
	for i, r := range out {
		if r != nil {
			idx[r.method+r.endpoint] = i
		}
	}
	for _, r := range incoming {
		if r == nil {
			continue
		}
		key := r.method + r.endpoint
		if i, ok := idx[key]; ok {
			out[i] = r
		} else {
			idx[key] = len(out)
			out = append(out, r)
		}
	}
	return out
}

// ReloadEndpoint reload a route and its filters
//
// The procedure use existing chi mux
func (s *apigw) ReloadEndpoint(ctx context.Context, method, endpoint string) error {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()

	routes, err := s.loadRoute(ctx, method, endpoint)
	if err != nil {
		s.log.Error("could not reload Integration Gateway routes", zap.Error(err))
		return err
	}

	s.PrepRoutes(ctx, routes...)

	// Merge the reloaded route into the current set and rebuild a fresh mux,
	// rather than mutating the one requests are being served from.
	merged := mergeRoutes(s.currentRoutes(), routes)
	s.publish(s.buildMux(merged), merged)
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

			regFilters []*st.ApigwFilter
		)

		// pipeline needs to know how to handle
		// async processers
		pipe.Async(r.meta.async)

		r.cfg = s.cfg
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

		// In case of webhooks, check if webhookAuth prefilter is defined
		// If not, deny all requests, else use what is defined.
		//
		// If no auth should be used, explicitly opt out.
		if _, isWebhook := r.meta.labels["human.webhookEvent"]; isWebhook {
			hasAuth := false
			for _, rf := range regFilters {
				if rf.Ref == "webhookAuth" {
					hasAuth = true
					break
				}
			}
			if !hasAuth {
				log.Warn("webhook route has no webhookAuth filter — blocking all requests")
				pipe.Add(&pipeline.Worker{
					Handler: func(rw http.ResponseWriter, r *http.Request) error {
						http.Error(rw, "", http.StatusUnauthorized)
						return nil
					},
					Name:   "webhookAuthDeny",
					Type:   types.PreFilter,
					Weight: filter.FilterWeight(0, types.PreFilter),
				})
			}
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
	if len(method) == 0 || len(endpoint) == 0 {
		return
	}

	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()

	// Rebuild from the current set, then override this endpoint with the default
	// 404 handler, and publish — never mutate the live mux in place.
	routes := s.currentRoutes()
	mx := s.buildMux(routes)
	mx.Method(method, endpoint, helperDefaultResponse(s.cfg, s.pr, s.log))
	s.publish(mx, routes)
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

	handler, err = s.reg.Merge(handler, enc, s.cfg)

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

func (s *apigw) UpdateSettings(ctx context.Context, cfg types.Config) {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()

	s.cfg = cfg

	s.reg = registry.NewRegistry(cfg)
	s.reg.Preload()

	_ = s.reloadLocked(ctx)
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
		labels := make(map[string]string, len(r.Meta.Labels))
		for k, v := range r.Meta.Labels {
			labels[k] = v.Val
		}

		route := &route{
			ID:       r.ID,
			endpoint: r.Endpoint,
			method:   r.Method,
			meta: routeMeta{
				debug:  r.Meta.Debug,
				async:  r.Meta.Async,
				labels: labels,
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
		labels := make(map[string]string, len(r.Meta.Labels))
		for k, v := range r.Meta.Labels {
			labels[k] = v.Val
		}

		rr = append(rr, &route{
			ID:       r.ID,
			endpoint: r.Endpoint,
			method:   r.Method,
			meta: routeMeta{
				debug:  r.Meta.Debug,
				async:  r.Meta.Async,
				labels: labels,
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
