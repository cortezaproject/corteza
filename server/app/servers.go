package app

import (
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/crusttech/human/server/assets"
	automationRest "github.com/crusttech/human/server/automation/rest"
	composeRest "github.com/crusttech/human/server/compose/rest"
	discoveryRest "github.com/crusttech/human/server/discovery/rest"
	"github.com/crusttech/human/server/docs"
	federationRest "github.com/crusttech/human/server/federation/rest"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/logger"

	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/webapp"
	systemRest "github.com/crusttech/human/server/system/rest"
	widgetRest "github.com/crusttech/human/server/system/rest/widget"
	"github.com/crusttech/human/server/system/scim"
	systemService "github.com/crusttech/human/server/system/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (app *HumanApp) mountHttpRoutes(r chi.Router) {
	var (
		ho = app.Opt.HTTPServer

		// Shared by the two halves of the OAuth discovery handshake, which are
		// mounted far apart: the challenge goes on the MCP route below, and the
		// document it points at goes on the root router, because RFC 9728 puts
		// it on the origin with the resource path appended.
		mcpResource = auth.OAuthProtectedResource{
			BasePath:      options.CleanBase(ho.BaseUrl, ho.ApiBaseUrl, "mcp"),
			Issuer:        app.Opt.Auth.BaseURL,
			Scopes:        []string{"api", "profile"},
			SslTerminated: ho.SslTerminated,
		}
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
			"client web application enabled",
			zap.String("baseUrl", options.CleanBase(ho.BaseUrl, ho.WebappBaseUrl)),
			zap.String("baseDir", ho.WebappBaseDir),
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

			// Public chatbot widget API (/api/widget/v1/*) — no admin token
			// validator, auth is per-session JWT. Mounted here so every
			// front-end origin configured on an Agent can reach it.
			widgetCtrl := widgetRest.New(
				systemService.DefaultStore,
				systemService.DefaultObsBus,
				systemService.DefaultAgenticRuntime,
				systemService.DefaultAiConversation,
				app.Opt.Auth.Secret,
				systemService.DefaultChatbotSession,
			)
			widgetCtrl.MountRoutes(r)
			// MCP tool surface (/api/mcp). Both middlewares are needed and they
			// do different jobs: HttpTokenValidator rejects a token minted for
			// another scope, HttpAuthenticatedOnly rejects a caller with no
			// token at all. The REST surface gets away with only the former
			// because its handlers reach RBAC, which denies the anonymous role.
			// MCP tool discovery never reaches RBAC — it lists every registered
			// tool with its full description straight from the registry — so
			// without the latter the whole tool catalogue is public.
			if app.McpServer != nil {
				r.Route("/mcp", func(r chi.Router) {
					// Outermost, so it sees the 401 either validator below
					// produces and can tell the client where to authenticate.
					r.Use(mcpResource.Challenge())
					r.Use(auth.HttpTokenValidator("api"))
					r.Use(auth.HttpAuthenticatedOnly())
					app.McpServer.MountRoutes(r)
				})
			}

			if app.Opt.Discovery.Enabled {
				r.Route("/discovery", discoveryRest.MountRoutes(app.Opt.Discovery))
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

		// clients discover the document under the issuer URL (AUTH_BASE_URL) that ends with /auth
		r.Handle("/auth/.well-known/openid-configuration", app.AuthService.WellKnownOpenIDConfiguration())

		// RFC 8414 puts the issuer's path *after* the well-known segment rather
		// than before it, and an OAuth client walks those variants before the
		// one above. They have to be answered here rather than left to fall
		// through: an unrouted path on this server returns an empty 200, not a
		// 404, so a client probing them finds a response that is neither a
		// document nor a miss, and may never try the form that works.
		//
		// The document is the same either way — an OpenID configuration is a
		// superset of RFC 8414 metadata — and it describes the one
		// authorization server this deployment has, whatever path was asked
		// for, so the wildcards disclose nothing the root path does not.
		r.Handle("/.well-known/oauth-authorization-server", app.AuthService.WellKnownOpenIDConfiguration())
		r.Handle("/.well-known/oauth-authorization-server/*", app.AuthService.WellKnownOpenIDConfiguration())
		r.Handle("/.well-known/openid-configuration/*", app.AuthService.WellKnownOpenIDConfiguration())

		if app.McpServer != nil {
			// Both forms: a client that knows which resource it wants appends
			// the resource path, one that does not asks for the bare document.
			r.Handle(auth.WellKnownProtectedResource, mcpResource.MetadataHandler())
			r.Handle(auth.WellKnownProtectedResource+"/*", mcpResource.MetadataHandler())
		}
	}()
}
