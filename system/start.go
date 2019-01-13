package service

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/SentimensRG/ctx/sigctx"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/internal/mail"
	migrate "github.com/crusttech/crust/system/db"
	"github.com/crusttech/crust/system/service"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/version"
)

var (
	jwtVerifier      (func(http.Handler) http.Handler)
	jwtAuthenticator (func(http.Handler) http.Handler)
	jwtEncoder       auth.TokenEncoder
)

func Init() error {
	// validate configuration
	if err := flags.Validate(); err != nil {
		return err
	}
	// JWT Auth
	if jwtAuth, err := auth.JWT(); err != nil {
		return errors.Wrap(err, "Error creating JWT Auth object")
	} else {
		jwtEncoder = jwtAuth
		jwtVerifier = jwtAuth.Verifier()
		jwtAuthenticator = jwtAuth.Authenticator()
	}

	mail.SetupDialer(flags.smtp)

	// start/configure database connection
	factory.Database.Add("default", flags.db.DSN)
	db := factory.Database.MustGet()

	// @todo: profiling as an external service?
	switch flags.db.Profiler {
	case "stdout":
		db.Profiler = &factory.Database.ProfilerStdout
	default:
		log.Println("No database query profiler selected")
	}

	// migrate database schema
	if err := migrate.Migrate(db); err != nil {
		return err
	}

	// configure resputil options
	resputil.SetConfig(resputil.Options{
		Pretty: flags.http.Pretty,
		Trace:  flags.http.Tracing,
		Logger: func(err error) {
			// @todo: error logging
		},
	})

	service.Init()

	return nil
}

func Start() error {
	log.Printf("Starting auth, version: %v, built on: %v", version.Version, version.BuildTime)
	log.Println("Starting http server on address " + flags.http.Addr)
	listener, err := net.Listen("tcp", flags.http.Addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can't listen on addr %s", flags.http.Addr))
	}

	if flags.monitor.Interval > 0 {
		go metrics.NewMonitor(flags.monitor.Interval)
	}
	go http.Serve(listener, Routes())

	var deadline = sigctx.New()
	<-deadline.Done()

	return nil
}
