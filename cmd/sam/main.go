package main

import (
	"flag"
	"log"
	"net"
	"os"

	"net/http"

	"github.com/go-chi/chi"

	project "github.com/crusttech/crust/sam"
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
	// set up flags
	var (
		addr = flag.String("addr", ":3000", "Listen address for HTTP server")
		dsn = flag.String("dsn", "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	)
	flag.Parse()

	// log to stdout not stderr
	log.SetOutput(os.Stdout)

	// set up database connection
	factory.Database.Add("default", dsn)
	db, err := factory.Database.Get()
	handleError(err, "Can't connect to database")
	db.Profiler = &factory.Database.ProfilerStdout

	// listen socket for http server
	log.Println("Starting http server on address " + *addr)
	listener, err := net.Listen("tcp", *addr)
	handleError(err, "Can't listen on addr "+*addr)

	// route options
	routeOptions, err := RouteOptions{}.New()
	handleError(err, "Error creating RouteOptions object")

	// mount routes
	r := chi.NewRouter()
	MountRoutes(r, routeOptions, project.MountRoutes)
	http.Serve(listener, r)
}
