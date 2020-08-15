package app

import (
	"context"
	composeRest "github.com/cortezaproject/corteza-server/compose/rest"
	messagingRest "github.com/cortezaproject/corteza-server/messaging/rest"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/webapp"
	systemRest "github.com/cortezaproject/corteza-server/system/rest"
	"github.com/go-chi/chi"
	"strings"
	"sync"
)

func (app *CortezaApp) Serve(ctx context.Context) (err error) {
	wg := &sync.WaitGroup{}

	{
		// @todo refactor wait-for out of HTTP API server.
		app.HttpServer = api.New(app.Log, app.Opt.HTTPServer, app.Opt.WaitFor)
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
			r.Route("/messaging", messagingRest.MountRoutes)
		})
	}

	if app.Opt.HTTPServer.WebappEnabled {
		r.Route("/"+webappBaseUrl, webapp.MakeWebappServer(app.Opt.HTTPServer))
	}
}
