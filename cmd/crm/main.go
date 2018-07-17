package main

import (
	"log"
	"net"
	"os"

	"net/http"

	"github.com/go-chi/chi"

	"github.com/namsral/flag"

	"github.com/crusttech/crust/rbac"
	"github.com/crusttech/crust/crm/rest"
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
	config := flags("crm")

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

	// mount routes
	r := chi.NewRouter()
	MountRoutes(r, routeOptions, rest.MountRoutes)
	http.Serve(listener, r)
}
