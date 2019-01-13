package service

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/SentimensRG/ctx/sigctx"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"github.com/titpetric/factory/resputil"

	migrate "github.com/crusttech/crust/crm/db"
	"github.com/crusttech/crust/crm/rest"
	crmService "github.com/crusttech/crust/crm/service"
	systemService "github.com/crusttech/crust/system/service"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/mail"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/version"
)

func Init() error {
	// validate configuration
	if err := flags.Validate(); err != nil {
		return err
	}

	mail.SetupDialer(flags.smtp)

	// start/configure database connection
	factory.Database.Add("default", flags.db.DSN)
	db, err := factory.Database.Get()
	if err != nil {
		return err
	}

	// @todo: profiling as an external service?
	switch flags.db.Profiler {
	case "stdout":
		db.Profiler = &factory.Database.ProfilerStdout
	default:
		fmt.Println("No database query profiler selected")
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

	systemService.Init()
	crmService.Init()

	return nil
}

func Start() error {
	var deadline = sigctx.New()

	log.Printf("Starting crm, version: %v, built on: %v", version.Version, version.BuildTime)
	log.Println("Starting http server on address " + flags.http.Addr)
	listener, err := net.Listen("tcp", flags.http.Addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can't listen on addr %s", flags.http.Addr))
	}

	// JWT Auth
	jwtAuth, err := auth.JWT()
	if err != nil {
		return errors.Wrap(err, "Error creating JWT Auth object")
	}

	r := chi.NewRouter()
	r.Use(handleCORS)

	// Only protect application routes with JWT
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth.Verifier(), jwtAuth.Authenticator())
		mountRoutes(r, flags.http, rest.MountRoutes(jwtAuth))
	})

	printRoutes(r, flags.http)
	mountSystemRoutes(r, flags.http)

	if flags.monitor.Interval > 0 {
		go metrics.NewMonitor(flags.monitor.Interval)
	}

	go http.Serve(listener, r)
	<-deadline.Done()

	return nil
}

// Sets up default CORS rules to use as a middleware
func handleCORS(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler(next)
}
