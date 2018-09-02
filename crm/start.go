package crm

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/SentimensRG/ctx/sigctx"

	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/crm/rest"

	"github.com/go-chi/cors"
	"github.com/titpetric/factory"
	"github.com/titpetric/factory/resputil"
)

func Init() error {
	// validate configuration
	if err := config.Validate(); err != nil {
		return err
	}

	// start/configure database connection
	factory.Database.Add("default", config.db.dsn)
	db, err := factory.Database.Get()
	if err != nil {
		return err
	}
	// @todo: profiling as an external service?
	switch config.db.profiler {
	case "stdout":
		db.Profiler = &factory.Database.ProfilerStdout
	default:
		fmt.Println("No database query profiler selected")
	}

	// configure resputil options
	resputil.SetConfig(resputil.Options{
		Pretty: config.http.pretty,
		Trace:  config.http.tracing,
		Logger: func(err error) {
			// @todo: error logging
		},
	})

	return nil
}

func Start() error {
	var ctx = sigctx.New()

	log.Println("Starting http server on address " + config.http.addr)
	listener, err := net.Listen("tcp", config.http.addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can't listen on addr %s", config.http.addr))
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
		mountRoutes(r, config, rest.MountRoutes(jwtAuth))
	})

	printRoutes(r, config)
	mountSystemRoutes(r, config)

	go http.Serve(listener, r)
	<-ctx.Done()

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
