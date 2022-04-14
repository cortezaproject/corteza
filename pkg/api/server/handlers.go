package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/assets"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/cortezaproject/corteza-server/webconsole"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// routes used when server is in waiting mode
func waitingRoutes(log *zap.Logger, httpOpt options.HttpServerOpt) (r chi.Router) {
	r = chi.NewRouter()
	r.Use(handleCORS)

	mountServiceHandlers(r, log, httpOpt, waiting)

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
func activeRoutes(log *zap.Logger, mountable []func(r chi.Router), opts *options.Options) (r chi.Router) {
	r = chi.NewRouter()
	r.Use(handleCORS)

	httpOpt := opts.HTTPServer
	authOpt := opts.Auth
	envOpt := opts.Environment

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

	mountServiceHandlers(r, log, httpOpt, active)

	r.HandleFunc(httpOpt.ApiBaseUrl, handleStaticPages(log, httpOpt, authOpt, "api-landing.html"))
	r.HandleFunc(httpOpt.ApiBaseUrl+"/", handleStaticPages(log, httpOpt, authOpt, "api-landing.html"))
	r.NotFound(handleStaticPages(log, httpOpt, authOpt, "api-404.html"))

	return
}

func mountServiceHandlers(r chi.Router, log *zap.Logger, opt options.HttpServerOpt, state uint32) {
	if opt.WebConsoleEnabled {
		path := "/console"
		log.Info("web console enabled (HTTP_SERVER_WEB_CONSOLE_ENABLED=true): " + path)
		r.Route(path, func(r chi.Router) {
			if len(opt.WebConsolePassword) > 0 {
				credentials := map[string]string{
					opt.WebConsoleUsername: opt.WebConsolePassword,
				}
				r.Use(middleware.BasicAuth("web-console", credentials))
			} else {
				// warn only in waiting state to avoid repeated log messages
				if state == waiting {
					// warn the user regardless of what environment Corteza is running in.
					log.Warn("SECURITY RISK: web console is enabled and unprotected, set " +
						"HTTP_SERVER_WEB_CONSOLE_USERNAME, HTTP_SERVER_WEB_CONSOLE_PASSWORD " +
						"if not running in development environment!")
				}
			}

			webconsole.Mount(r)
			mountDebugLogViewer(r, log)

			// redirect from /console to /console/ui/
			r.Mount("/", http.RedirectHandler(path+"/ui", http.StatusTemporaryRedirect))
		})
	}

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

// @todo move all these routes under /console and
//       output JSON instead of plain raw text
func mountDebugHandler(r chi.Router, log *zap.Logger) {
	log.Debug("route debugger enabled: /__routes")
	r.Get("/__routes", debugRoutes(r))

	log.Debug("profiler enabled: /debug/pprof")
	r.Mount("/debug", middleware.Profiler())

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

func mountDebugLogViewer(r chi.Router, log *zap.Logger) {
	var (
		path = "/server-log-feed"
	)

	r.Get(path+".json", func(w http.ResponseWriter, r *http.Request) {
		var (
			after int = 0
			limit int = 100
			err   error
			q     = r.URL.Query()
		)

		if aux := q.Get("after"); len(aux) > 0 {
			after, err = strconv.Atoi(aux)
			if err != nil {
				errors.ProperlyServeHTTP(w, r, errors.InvalidData("invalid value format for after: %v", err), false)
				return
			}
		}

		if aux := q.Get("limit"); len(aux) > 0 {
			limit, err = strconv.Atoi(aux)
			if err != nil {
				errors.ProperlyServeHTTP(w, r, errors.InvalidData("invalid value format for limit: %v", err), false)
				return
			}
		}

		_, _ = logger.WriteLogBuffer(w, after, limit)
	})
}

func handleStaticPages(log *zap.Logger, hOpt options.HttpServerOpt, aOpt options.AuthOpt, file string) http.HandlerFunc {
	// "good-enough" for now, plan to move to templates when
	// merging with auth
	const linkTpl = `<a class="btn btn-light font-weight-bold text-dark m-2" href="%s">%s</a>`
	var (
		links = make([]string, 0)
		buf   []byte

		placeholder = []byte("<!-- links -->")
	)

	links = append(links, fmt.Sprintf(linkTpl, aOpt.BaseURL, "Login"))

	if hOpt.ApiEnabled {
		links = append(links, fmt.Sprintf(linkTpl, "https://docs.cortezaproject.org/", "Documentation"))
	}

	if hOpt.WebConsoleEnabled {
		links = append(links, fmt.Sprintf(linkTpl, "/console", "Console"))
	}

	page, err := assets.Files(log, hOpt.AssetsPath).Open(file)
	if err != nil {
		log.Warn("could not open static page", zap.String("file", file), zap.Error(err))
	}

	buf, err = io.ReadAll(page)
	if err != nil {
		log.Warn("could not prepare static page", zap.String("file", file), zap.Error(err))
	}

	buf = bytes.ReplaceAll(buf, placeholder, []byte(strings.Join(links, "")))

	return func(w http.ResponseWriter, r *http.Request) {
		if page == nil {
			// fallback to default 404 handler
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(buf)
	}
}
