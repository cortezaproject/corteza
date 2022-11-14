package server

import (
	"github.com/cortezaproject/corteza-server-discovery/pkg/auth"
	"github.com/cortezaproject/corteza-server-discovery/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server-discovery/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"path"
	"strings"
)

// routes used when server is in waiting mode
func waitingRoutes(log *zap.Logger, httpOpt options.HttpServerOpt) (r chi.Router) {
	r = chi.NewRouter()
	r.Use(handleCORS)

	mountServiceHandlers(r, log, httpOpt, waiting)

	return
}

// routes used when server in shutdown mode
func shutdownRoutes() (r chi.Router) {
	r = chi.NewRouter()

	return
}

// routes used when in active mode
func activeRoutes(log *zap.Logger, mountable []func(r chi.Router), envOpt options.EnvironmentOpt, httpOpt options.HttpServerOpt, searcherOpt options.SearcherOpt) (r chi.Router) {
	r = chi.NewRouter()
	r.Use(handleCORS)

	r.Route("/"+strings.TrimPrefix(httpOpt.BaseUrl, "/"), func(r chi.Router) {
		// Base middleware, CORS, RealIP, RequestID, context-logger
		r.Use(BaseMiddleware(envOpt.IsProduction(), log)...)

		// Verifies JWT in headers, cookies, ...
		r.Use(auth.HttpTokenVerifier)

		for _, mount := range mountable {
			mount(r)
		}

	})

	if httpOpt.BaseUrl != "/" {
		r.Handle("/", http.RedirectHandler(httpOpt.BaseUrl, http.StatusTemporaryRedirect))
	}

	mountServiceHandlers(r, log, httpOpt, active)

	return
}

func mountServiceHandlers(r chi.Router, log *zap.Logger, opt options.HttpServerOpt, state uint32) {
	if opt.EnableVersionRoute {
		mountVersionHandler(r, log, opt.BaseUrl)
	}

	if opt.EnableHealthcheckRoute {
		mountHealthCheckHandler(r, log, opt.BaseUrl)
	}
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
