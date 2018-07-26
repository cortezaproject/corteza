package main

import (
	"log"
	"net"
	"os"

	"net/http"

	"github.com/go-chi/chi"

	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/rbac"
	"github.com/crusttech/crust/sam/rest"
	"github.com/crusttech/crust/sam/websocket"
	"github.com/labstack/gommon/random"
	"github.com/titpetric/factory"
)

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func main() {
	config := flags("sam", rbac.Flags, websocket.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)

	// set up database connection
	factory.Database.Add("default", config.dbDSN)
	db, err := factory.Database.Get()
	handleError(err, "Can't connect to database")
	db.Profiler = &factory.Database.ProfilerStdout

	// listen socket for http server
	log.Println("Starting http server on address " + config.httpAddr)
	listener, err := net.Listen("tcp", config.httpAddr)
	handleError(err, "Can't listen on addr "+config.httpAddr)

	// route options
	routeOptions, err := RouteOptions{}.New()
	handleError(err, "Error creating RouteOptions object")

	r := chi.NewRouter()

	// JWT Auth
	jwtAuth := auth.JWT([]byte(config.jwtSecret))
	r.Use(jwtAuth.Verifier(), jwtAuth.Authenticator())

	if len(config.jwtSecret) == 0 {
		println("Environment variable JWT_SECRET not set! Add next line to your .env file:")
		println("JWR_SECRET=" + random.String(64, random.Alphabetic))
	}

	// mount REST & WS routes
	MountRoutes(r, routeOptions, rest.MountRoutes(jwtAuth), websocket.MountRoutes())

	http.Serve(listener, r)
}
