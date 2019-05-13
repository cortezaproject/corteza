package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"
	"go.uber.org/zap"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/db"
	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/mail"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/settings"
	migrate "github.com/crusttech/crust/system/db"
	"github.com/crusttech/crust/system/internal/auth/external"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/service"
)

func Init(ctx context.Context) (err error) {
	// validate configuration
	if err = flags.Validate(); err != nil {
		return
	}

	mail.SetupDialer(flags.smtp)

	if err = InitDatabase(ctx); err != nil {
		return
	}

	if flags.http.ClientTSLInsecure {
		// This will allow HTTPS requests to insecure hosts (expired, wrong host, self signed, untrusted root...)
		// With this enabled, features like OIDC auto-discovery should work on any of examples found on badssl.com.
		//
		// With SYSTEM_HTTP_CLIENT_TSL_INSECURE=0 (default) next command returns 404 error (expected)
		// > ./system-cli external-auth auto-discovery foo-tsl-1 https://expired.badssl.com/
		//
		// Without SYSTEM_HTTP_CLIENT_TSL_INSECURE=1 next command returns "x509: certificate has expired or is not yet valid"
		// > ./system-cli external-auth auto-discovery foo-tsl-1 https://expired.badssl.com/
		//
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// configure resputil options
	resputil.SetConfig(resputil.Options{
		Pretty: flags.http.Pretty,
		Trace:  flags.http.Tracing,
		Logger: func(err error) {
			// @todo: error logging
		},
	})

	// Use JWT secret for hmac signer for now
	auth.DefaultSigner = auth.HmacSigner(flags.jwt.Secret)
	auth.DefaultJwtHandler, err = auth.JWT(flags.jwt.Secret, flags.jwt.Expiry)
	if err != nil {
		return
	}

	// Don't change this, it needs database connection
	if err = service.Init(ctx); err != nil {
		return
	}

	return nil
}

func InitDatabase(ctx context.Context) error {
	// start/configure database connection
	db, err := db.TryToConnect(ctx, "system", flags.db.DSN, flags.db.Profiler)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}

	// migrate database schema
	if err := migrate.Migrate(db); err != nil {
		return err
	}

	return nil
}

func StartWatchers(ctx context.Context) {
	service.Watchers(ctx)
}

func StartRestAPI(ctx context.Context) error {
	// Load settings from the database,
	// for now, only at start-up time.
	settingService := settings.NewService(settings.NewRepository(repository.DB(ctx), "sys_settings")).With(ctx)

	// Setup goth/external authentication
	external.Init(settingService)

	logger.Default().Info("Starting HTTP server", zap.String("address", flags.http.Addr))
	listener, err := net.Listen("tcp", flags.http.Addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can't listen on addr %s", flags.http.Addr))
	}

	if flags.monitor.Interval > 0 {
		go metrics.NewMonitor(flags.monitor.Interval)
	}

	go http.Serve(listener, Routes(ctx))
	<-ctx.Done()

	return nil
}
