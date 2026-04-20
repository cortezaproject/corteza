package app

import (
	"context"
	"fmt"

	"github.com/crusttech/human/extra/server-discovery/indexer"
	"github.com/crusttech/human/extra/server-discovery/pkg/auth"
	"github.com/crusttech/human/extra/server-discovery/pkg/es"
	"github.com/crusttech/human/extra/server-discovery/pkg/healthcheck"
	"github.com/crusttech/human/extra/server-discovery/searcher"
)

const (
	bootLevelWaiting = iota
	bootLevelSetup
	bootLevelStoreInitialized
	bootLevelProvisioned
	bootLevelServicesInitialized
	bootLevelActivated
)

// Setup configures all required services
func (app *HumanDiscoveryApp) Setup() (err error) {
	app.lvl = bootLevelSetup

	hcd := healthcheck.Defaults()
	if app.Opt.Searcher.Enabled {
		hcd.Add(searcher.Healthcheck, "OpenSearch")
	}

	return nil
}

// InitStore initializes open search store and runs upgrade procedures
func (app *HumanDiscoveryApp) InitStore(ctx context.Context) (err error) {
	if app.lvl >= bootLevelStoreInitialized {
		// Is store already initialised?
		return nil
	} else if err = app.Setup(); err != nil {
		// Initialize previous level
		return err
	}

	app.lvl = bootLevelStoreInitialized
	return nil
}

// Provision instance with configuration and settings
// by importing preset configurations and running autodiscovery procedures
func (app *HumanDiscoveryApp) Provision(ctx context.Context) (err error) {
	if app.lvl >= bootLevelProvisioned {
		return
	}

	if err = app.InitStore(ctx); err != nil {
		return err
	}

	app.lvl = bootLevelProvisioned
	return
}

// InitServices initializes all services used
func (app *HumanDiscoveryApp) InitServices(ctx context.Context) (err error) {
	if app.lvl >= bootLevelServicesInitialized {
		return nil
	}

	if err = app.Provision(ctx); err != nil {
		return err
	}

	if auth.HttpTokenVerifier, err = auth.TokenVerifierMiddlewareWithSecretSigner(string(app.Opt.Searcher.JwtSecret)); err != nil {
		return fmt.Errorf("could not set token verifier")
	}

	//initialize Elastic search client
	esClient, err := es.Connect(app.Opt.ES)
	if err != nil {
		return fmt.Errorf("failed to initialize Elasticsearch client: %s", err.Error())
	}

	if app.Opt.Indexer.Enabled {
		err = indexer.Initialize(ctx, app.Log, indexer.Config{
			Human:      app.Opt.Human,
			ES:           app.Opt.ES,
			Indexer:      app.Opt.Indexer,
			VectorSearch: app.Opt.VectorSearch,
		}, esClient)
		if err != nil {
			return
		}
	}

	if app.Opt.Searcher.Enabled {
		err = searcher.Initialize(ctx, app.Log, searcher.Config{
			Human:    app.Opt.Human,
			ES:         app.Opt.ES,
			HttpServer: app.Opt.HTTPServer,
			Searcher:   app.Opt.Searcher,
		}, esClient)
		if err != nil {
			return
		}
	}

	app.lvl = bootLevelServicesInitialized
	return
}

// Activate start all internal services and watchers
func (app *HumanDiscoveryApp) Activate(ctx context.Context) (err error) {
	if app.lvl >= bootLevelActivated {
		return
	}

	if err := app.InitServices(ctx); err != nil {
		return err
	}

	if app.Opt.Indexer.Enabled {
		indexer.Watchers(ctx)
	}

	app.lvl = bootLevelActivated

	return nil
}
