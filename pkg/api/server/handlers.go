package server

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// routes used when server is in waiting mode
func waitingRoutes(log *zap.Logger, httpOpt options.HttpServerOpt) (r chi.Router) {
	r = chi.NewRouter()
	mountServiceHandlers(r, log, httpOpt)

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		// For non GET requests, return 503 (service unavailable)
		errors.ServeHTTPWithCode(w, r,
			http.StatusServiceUnavailable,
			fmt.Errorf("corteza server initializing"),
			true,
		)
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Refresh the page in 15 seconds
		w.Header().Set("Refresh", "15; url=/")
		_, _ = fmt.Fprint(w, "Corteza server initializing\n\n")
		if httpOpt.EnableHealthcheckRoute {
			healthcheck.Defaults().Run(r.Context()).WriteTo(w)
		}
	})

	return
}

// routes used when server in shutdown mode
func shutdownRoutes() (r chi.Router) {
	r = chi.NewRouter()

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		// For non GET requests, return 503 (service unavailable)
		errors.ServeHTTPWithCode(w, r,
			http.StatusServiceUnavailable,
			fmt.Errorf("corteza server shutting down"),
			true,
		)
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Refresh the page in 15 seconds
		w.Header().Set("Refresh", "15; url=/")
		_, _ = fmt.Fprint(w, "corteza server shutting down")
	})
	return
}

// routes used when in active mode
func activeRoutes(log *zap.Logger, mountable []func(r chi.Router), envOpt options.EnvironmentOpt, httpOpt options.HttpServerOpt) (r chi.Router) {
	r = chi.NewRouter()

	r.Route("/"+strings.TrimPrefix(httpOpt.BaseUrl, "/"), func(r chi.Router) {
		// Reports error to Sentry if enabled
		if httpOpt.EnablePanicReporting {
			r.Use(sentryMiddleware())
		}

		if httpOpt.EnableMetrics {
			// Metrics tracking middleware
			r.Use(metricsMiddleware(httpOpt.MetricsServiceLabel))
		}

		// Handle panic (sets 500 server error headers)
		//r.Use(handlePanic)

		// Base middleware, CORS, RealIP, RequestID, context-logger
		r.Use(BaseMiddleware(envOpt.IsProduction(), log)...)

		// Logging request if enabled
		if httpOpt.LogRequest {
			r.Use(LogRequest)
		}

		// Logging response if enabled
		if httpOpt.LogResponse {
			r.Use(LogResponse)
		}

		// Verifies JWT in headers, cookies, ...
		r.Use(auth.HttpTokenVerifier)

		for _, mount := range mountable {
			mount(r)
		}

		if httpOpt.EnableMetrics {
			metricsMount(r, httpOpt.MetricsUsername, httpOpt.MetricsPassword)
		}

	})

	if httpOpt.BaseUrl != "/" {
		r.Handle("/", http.RedirectHandler(httpOpt.BaseUrl, http.StatusTemporaryRedirect))

	}

	mountServiceHandlers(r, log, httpOpt)
	return
}

func mountServiceHandlers(r chi.Router, log *zap.Logger, opt options.HttpServerOpt) {
	if opt.EnableDebugRoute {
		mountDebugHandler(r, log)
	}

	if opt.EnableVersionRoute {
		mountVersionHandler(r, log, opt.BaseUrl)
	}

	if opt.EnableHealthcheckRoute {
		mountHealthCheckHandler(r, log, opt.BaseUrl)
	}
}

func mountDebugHandler(r chi.Router, log *zap.Logger) {
	log.Debug("route debugger enabled: /__routes")
	r.Get("/__routes", debugRoutes(r))

	log.Debug("profiler enabled: /__profiler")
	r.Mount("/__profiler", middleware.Profiler())

	log.Debug("eventbus handlers debug enabled: /__eventbus")
	r.Get("/__eventbus", debugEventbus())

	log.Debug("corredor service debug enabled: /__corredor")
	r.Get("/__corredor", debugCorredor())
}

func mountVersionHandler(r chi.Router, log *zap.Logger, basePath string) {
	var (
		dPath   = "/version"
		sPath   = path.Join(basePath, dPath)
		handler = func(w http.ResponseWriter, r *http.Request) {
			api.Send(w, r, struct {
				BuildTime string `json:"buildTime"`
				Version   string `json:"version"`
			}{version.BuildTime, version.Version})
		}
	)

	r.Get(dPath, handler)
	log.Debug("version route enabled: " + dPath)

	if dPath != sPath {
		r.Get(sPath, handler)
		log.Debug("version route enabled: " + sPath)
	}

}

func mountHealthCheckHandler(r chi.Router, log *zap.Logger, basePath string) {
	// default & sub path for health-check endpoint
	var (
		dPath = "/healthcheck"
		sPath = path.Join(basePath, dPath)
	)

	r.Get(dPath, healthcheck.HttpHandler())
	log.Debug("health check route enabled: " + sPath)

	if dPath != sPath {
		r.Get(sPath, healthcheck.HttpHandler())
		log.Debug("health check route enabled: " + sPath)
	}

}
