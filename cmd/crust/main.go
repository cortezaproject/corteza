package main

import (
	"log"
	"net"
	"os"
	"path"

	"net/http"

	context "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/namsral/flag"

	crm "github.com/crusttech/crust/crm"
	messaging "github.com/crusttech/crust/messaging"
	system "github.com/crusttech/crust/system"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/middleware"
	"github.com/crusttech/crust/internal/rbac"
	"github.com/crusttech/crust/internal/routes"
	"github.com/crusttech/crust/internal/subscription"
	"github.com/crusttech/crust/internal/version"
)

// Serves index.html in case the requested file isn't found (or some other os.Stat error)
func serveIndex(assetPath string, indexPath string, serve http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage := path.Join(assetPath, indexPath)
		requestedPage := path.Join(assetPath, r.URL.Path)
		_, err := os.Stat(requestedPage)
		if err != nil {
			http.ServeFile(w, r, indexPage)
			return
		}
		serve.ServeHTTP(w, r)
	}
}

func main() {
	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting "+os.Args[0]+", version: %v, built on: %v", version.Version, version.BuildTime)

	ctx := context.AsContext(sigctx.New())

	var flags struct {
		http    *config.HTTP
		monitor *config.Monitor
	}
	flags.http = new(config.HTTP).Init()
	flags.monitor = new(config.Monitor).Init()

	crm.Flags("crm")
	messaging.Flags("messaging")
	system.Flags("system")

	auth.Flags()
	rbac.Flags()
	subscription.Flags()

	flag.Parse()

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "help":
		flag.PrintDefaults()
	default:
		// Initialize configuration of our services
		if err := system.Init(); err != nil {
			log.Fatalf("Error initializing system: %+v", err)
		}
		if err := crm.Init(); err != nil {
			log.Fatalf("Error initializing crm: %+v", err)
		}
		if err := messaging.Init(); err != nil {
			log.Fatalf("Error initializing messaging: %+v", err)
		}

		// Checks subscription, will os.Exit(1) if there is an error
		// Disabled for now, system service is the only one that validates subscription
		// ctx = subscription.Monitor(ctx)

		log.Println("Starting http server on address " + flags.http.Addr)
		listener, err := net.Listen("tcp", flags.http.Addr)
		if err != nil {
			log.Fatalf("Can't listen on addr %s", flags.http.Addr)
		}

		if flags.monitor.Interval > 0 {
			go metrics.NewMonitor(flags.monitor.Interval)
		}

		r := chi.NewRouter()

		// logging, cors and such
		middleware.Mount(ctx, r, flags.http)

		r.Route("/api", func(r chi.Router) {
			r.Route("/crm", func(r chi.Router) {
				crm.MountRoutes(ctx, r)
			})
			r.Route("/messaging", func(r chi.Router) {
				messaging.MountRoutes(ctx, r)
			})
			r.Route("/system", func(r chi.Router) {
				system.MountRoutes(ctx, r)
			})
			middleware.MountSystemRoutes(ctx, r, flags.http)
		})

		fileserver := http.FileServer(http.Dir("public_html"))

		for _, service := range []string{"system", "messaging", "crm"} {
			r.HandleFunc("/"+service+"*", serveIndex("public_html", "crm/index.html", fileserver))
		}
		r.HandleFunc("/*", serveIndex("public_html", "index.html", fileserver))

		routes.Print(r)

		go http.Serve(listener, r)
		<-ctx.Done()
	}
}
