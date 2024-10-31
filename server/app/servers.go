package app

import (
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza/server/assets"
	automationRest "github.com/cortezaproject/corteza/server/automation/rest"
	composeRest "github.com/cortezaproject/corteza/server/compose/rest"
	discoveryRest "github.com/cortezaproject/corteza/server/discovery/rest"
	"github.com/cortezaproject/corteza/server/docs"
	federationRest "github.com/cortezaproject/corteza/server/federation/rest"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/pkg/webapp"
	systemRest "github.com/cortezaproject/corteza/server/system/rest"
	"github.com/cortezaproject/corteza/server/system/scim"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (app *CortezaApp) mountHttpRoutes(r chi.Router) {
	var (
		ho = app.Opt.HTTPServer
	)

	func() {
		// asset serving has some overlap with auth assets, web-console and webapp serving
		// and might be joined with one or more of them in the later version

		var (
			url   = options.CleanBase(ho.BaseUrl, "assets")
			aPath = ho.AssetsPath
			files = assets.Files(app.Log, aPath)
		)

		r.Handle(url+"/*", http.StripPrefix(url+"/", http.FileServer(http.FS(files))))

		if aPath != "" {
			app.Log.Info("custom web assets mounted", zap.String("url", url), zap.String("path", aPath))
		} else {
			app.Log.Info("embedded web assets mounted", zap.String("url", url))
		}
	}()

	func() {
		if ho.WebappEnabled && ho.ApiEnabled && ho.ApiBaseUrl == ho.WebappBaseUrl {
			app.Log.
				Warn("client web applications and api can not use the same base URL: '" + ho.WebappBaseUrl + "'")
			ho.WebappEnabled = false
		}

		if !ho.WebappEnabled {
			app.Log.Info("client web applications disabled")
			return
		}

		r.Route(options.CleanBase(ho.WebappBaseUrl), webapp.MakeWebappServer(app.Log, ho, app.Opt.Auth, app.Opt.Discovery, app.Opt.Sentry))

		app.Log.Info(
			"client web applications enabled",
			zap.String("baseUrl", options.CleanBase(ho.BaseUrl, ho.WebappBaseUrl)),
			zap.String("baseDir", ho.WebappBaseDir),
			zap.Strings("apps", strings.Split(ho.WebappList, ",")),
		)
	}()

	// Auth server
	app.AuthService.MountHttpRoutes(ho.BaseUrl, r)

	func() {
		if !ho.ApiEnabled {
			app.Log.Info("JSON REST API disabled")
			return
		}

		r.Route(options.CleanBase(ho.ApiBaseUrl), func(r chi.Router) {
			var fullpathAPI = "/" + strings.TrimPrefix(options.CleanBase(ho.BaseUrl, ho.ApiBaseUrl), "/")

			app.Log.Info(
				"JSON REST API enabled",
				zap.String("baseUrl", fullpathAPI),
			)

			r.Route("/system", systemRest.MountRoutes())
			r.Route("/automation", automationRest.MountRoutes())
			r.Route("/compose", composeRest.MountRoutes())
			r.Route("/websocket", app.WsServer.MountRoutes)

			if app.Opt.Discovery.Enabled {
				r.Route("/discovery", discoveryRest.MountRoutes())
			}

			if app.Opt.Federation.Enabled {
				r.Route("/federation", federationRest.MountRoutes(app.Opt.Limit))
			}

			var fullpathDocs = options.CleanBase(ho.BaseUrl, ho.ApiBaseUrl, "docs")
			app.Log.Info(
				"API docs enabled",
				zap.String("baseUrl", fullpathDocs),
			)

			r.Handle("/docs", http.RedirectHandler(fullpathDocs+"/", http.StatusPermanentRedirect))
			r.Handle("/docs*", http.StripPrefix(fullpathDocs, http.FileServer(docs.GetFS())))

			var fullpathGateway = options.CleanBase(ho.BaseUrl, ho.ApiBaseUrl, "gateway")
			r.Handle("/gateway*", http.StripPrefix(fullpathGateway, app.ApigwService))
		})
	}()

	func() {
		if !app.Opt.SCIM.Enabled {
			return
		}

		if app.Opt.SCIM.Secret == "" {
			app.Log.
				Error("SCIM secret empty")
		}

		var (
			baseUrl         = app.Opt.SCIM.BaseURL
			extIdValidation *regexp.Regexp
			err             error
		)

		if len(app.Opt.SCIM.ExternalIdValidation) > 0 {
			extIdValidation, err = regexp.Compile(app.Opt.SCIM.ExternalIdValidation)
		}

		if err != nil {
			app.Log.Error("failed to compile SCIM external ID validation", zap.Error(err))
			return
		}

		app.Log.Debug(
			"SCIM enabled",
			zap.String("baseUrl", path.Join(app.Opt.HTTPServer.BaseUrl, baseUrl)),
			logger.Mask("secret", app.Opt.SCIM.Secret),
		)

		r.Route(baseUrl, func(r chi.Router) {
			if !app.Opt.Environment.IsDevelopment() {
				r.Use(scim.Guard(app.Opt.SCIM))
			}

			scim.Routes(r, scim.Config{
				ExternalIdAsPrimary: app.Opt.SCIM.ExternalIdAsPrimary,
				ExternalIdValidator: extIdValidation,
			})
		})
	}()

	func() {
		r.Handle("/.well-known/openid-configuration", app.AuthService.WellKnownOpenIDConfiguration())
	}()
}
