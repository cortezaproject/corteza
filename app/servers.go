package app

import (
	"context"
	composeRest "github.com/cortezaproject/corteza-server/compose/rest"
	federationRest "github.com/cortezaproject/corteza-server/federation/rest"
	messagingRest "github.com/cortezaproject/corteza-server/messaging/rest"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/webapp"
	systemRest "github.com/cortezaproject/corteza-server/system/rest"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"strings"
	"sync"
)

func (app *CortezaApp) Serve(ctx context.Context) (err error) {
	wg := &sync.WaitGroup{}

	{
		// @todo refactor wait-for out of HTTP API server.
		app.HttpServer = server.New(app.Log, app.Opt.Environment, app.Opt.HTTPServer, app.Opt.WaitFor)
		app.HttpServer.MountRoutes(app.mountHttpRoutes)

		wg.Add(1)
		go func() {
			app.HttpServer.Serve(actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_API_REST))
			wg.Done()
		}()
	}

	{
		//wg.Add(1)
		//go func(ctx context.Context) {
		//	grpcApi.Serve(actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_API_GRPC))
		//	wg.Done()
		//}(ctx)
	}

	// Wait for all servers to be done
	wg.Wait()

	return nil
}

func (app *CortezaApp) mountHttpRoutes(r chi.Router) {
	var (
		apiBaseUrl    = strings.Trim(app.Opt.HTTPServer.ApiBaseUrl, "/")
		webappBaseUrl = strings.Trim(app.Opt.HTTPServer.WebappBaseUrl, "/")
	)

	if app.Opt.HTTPServer.ApiEnabled {
		r.Route("/"+apiBaseUrl, func(r chi.Router) {
			r.Route("/system", systemRest.MountRoutes)
			r.Route("/compose", composeRest.MountRoutes)
			r.Route("/messaging", func(r chi.Router) {
				messagingRest.MountRoutes(r)
				app.WsServer.ApiServerRoutes(r)
			})

			r.HandleFunc("/docs*", server.ServeDocs("/"+apiBaseUrl+"/docs"))
		})
	}

	if app.Opt.HTTPServer.WebappEnabled {
		app.Log.Debug(
			"serving web applications",
			zap.String("baseUrl", webappBaseUrl),
			zap.String("baseDir", app.Opt.HTTPServer.WebappBaseDir),
		)
		r.Route("/"+webappBaseUrl, webapp.MakeWebappServer(app.Opt.HTTPServer))
	}

	if app.Opt.Federation.Enabled {
		r.Route("/federation", federationRest.MountRoutes)
	}
}
