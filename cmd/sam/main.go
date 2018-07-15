package main

import (
	"flag"
	"log"
	"net"
	"os"

	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"

	"github.com/crusttech/crust/sam/rest"
	"github.com/titpetric/factory"
)

const (
	defaultAddr = ":3000"
	defaultDsn  = "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci"

	envVarKey_HTTP_ADDR = "SAM_HTTP_ADDR"
	envVarKey_DB_DSN    = "SAM_DB_DSN"
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
	var envHttpAddr, envDbDsn string
	var has bool

	if envHttpAddr, has = os.LookupEnv(envVarKey_HTTP_ADDR); !has {
		envHttpAddr = defaultAddr
	}

	if envDbDsn, has = os.LookupEnv(envVarKey_DB_DSN); !has {
		envDbDsn = defaultDsn
	}

	var (
		// set up flags
		addr = flag.String("addr", envHttpAddr, "Listen address for HTTP server")
		dsn  = flag.String("dsn", envDbDsn, "DSN for database connection")
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
	MountRoutes(r, routeOptions, rest.MountRoutes)
	http.Serve(listener, r)
}
