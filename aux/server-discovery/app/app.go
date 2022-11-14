package app

import (
	"context"
	"github.com/cortezaproject/corteza-server-discovery/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/davecgh/go-spew/spew"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	httpApiServer interface {
		//MountRoutes(mm ...func(chi.Router))
		//Serve(ctx context.Context)
		Serve(ctx context.Context)
		Activate(mm ...func(chi.Router))
		Shutdown()
	}

	CortezaDiscoveryApp struct {
		Opt *options.Options
		lvl int
		Log *zap.Logger

		// Servers
		HttpServer httpApiServer
	}
)

var (
	_ *spew.ConfigState = nil
	_ esutil.BulkIndexer
)

func New() (app *CortezaDiscoveryApp, err error) {
	app = &CortezaDiscoveryApp{
		Log: logger.MakeDebugLogger().WithOptions(zap.AddStacktrace(zap.PanicLevel)),
	}
	app.Opt, err = options.Init()

	if err != nil {
		return
	}

	return
}

func (app CortezaDiscoveryApp) Serve(ctx context.Context) (err error) {
	return
}

//func (app CortezaDiscoveryApp) InitService(ctx context.Context) (err error) {
//	// Initialize indexer service
//	err = indexer.Initialize(ctx, app.Log, indexer.Config{
//		ES:      app.Opt.ES,
//		Indexer: app.Opt.Indexer,
//	})
//	if err != nil {
//		return err
//	}
//
//	return
//}
