package app

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server-discovery/indexer"
	"github.com/cortezaproject/corteza-server-discovery/pkg/auth"
	"github.com/cortezaproject/corteza-server-discovery/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server-discovery/searcher"
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
func (app *CortezaDiscoveryApp) Setup() (err error) {
	app.lvl = bootLevelSetup

	hcd := healthcheck.Defaults()
	if app.Opt.Searcher.Enabled {
		hcd.Add(searcher.Healthcheck, "OpenSearch")
	}

	return nil
}

// InitStore initializes open search store and runs upgrade procedures
func (app *CortezaDiscoveryApp) InitStore(ctx context.Context) (err error) {
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
func (app *CortezaDiscoveryApp) Provision(ctx context.Context) (err error) {
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
func (app *CortezaDiscoveryApp) InitServices(ctx context.Context) (err error) {
	if app.lvl >= bootLevelServicesInitialized {
		return nil
	}

	if err = app.Provision(ctx); err != nil {
		return err
	}

	if auth.HttpTokenVerifier, err = auth.TokenVerifierMiddlewareWithSecretSigner(string(app.Opt.Searcher.JwtSecret)); err != nil {
		return fmt.Errorf("could not set token verifier")
	}

	if app.Opt.Indexer.Enabled {
		err = indexer.Initialize(ctx, app.Log, indexer.Config{
			Corteza: app.Opt.Corteza,
			ES:      app.Opt.ES,
			Indexer: app.Opt.Indexer,
		})
		if err != nil {
			return
		}
	}

	if app.Opt.Searcher.Enabled {
		err = searcher.Initialize(ctx, app.Log, searcher.Config{
			Corteza:    app.Opt.Corteza,
			ES:         app.Opt.ES,
			HttpServer: app.Opt.HTTPServer,
			Searcher:   app.Opt.Searcher,
		})
		if err != nil {
			return
		}
	}

	app.lvl = bootLevelServicesInitialized
	return
}

// Activate start all internal services and watchers
func (app *CortezaDiscoveryApp) Activate(ctx context.Context) (err error) {
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
