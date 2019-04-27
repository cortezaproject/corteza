package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	migrate "github.com/crusttech/crust/crm/db"
	"github.com/crusttech/crust/crm/internal/service"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/db"
	"github.com/crusttech/crust/internal/mail"
	"github.com/crusttech/crust/internal/metrics"
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
		return err
	}

	// Don't change this to init(), it needs Database
	return service.Init()
}

func InitDatabase(ctx context.Context) error {
	// start/configure database connection
	db, err := db.TryToConnect(ctx, "crm", flags.db.DSN, flags.db.Profiler)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}

	// migrate database schema
	if err := migrate.Migrate(db); err != nil {
		return err
	}

	return nil
}

func StartRestAPI(ctx context.Context) error {
	log.Println("Starting http server on address " + flags.http.Addr)
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
