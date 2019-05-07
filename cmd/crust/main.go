package main

import (
	"net"
	"net/http"
	"os"
	"path"

	context "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/namsral/flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	compose "github.com/crusttech/crust/compose"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/logger"
	messaging "github.com/crusttech/crust/messaging"
	system "github.com/crusttech/crust/system"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/middleware"
	"github.com/crusttech/crust/internal/subscription"
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
	// Initialize default logger
	logger.Init(zapcore.DebugLevel)
	log := logger.Default().Named("crust")

	// New signal-bond context that we will use and
	// will get terminated (Done()) on SIGINT or SIGTERM
	ctx := context.AsContext(sigctx.New())

	// Bind default logger to context
	ctx = logger.ContextWithValue(ctx, log)

	var flags struct {
		http    *config.HTTP
		monitor *config.Monitor
	}

	flags.http = new(config.HTTP).Init()
	flags.monitor = new(config.Monitor).Init()

	compose.Flags("compose")
	messaging.Flags("messaging")
	system.Flags("system")

	authJwtFlags := new(config.JWT).Init()

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
		if err := system.Init(ctx); err != nil {
			log.Fatal("failed to initialize system", zap.Error(err))
		}
		if err := compose.Init(ctx); err != nil {
			log.Fatal("failed to initialize compose", zap.Error(err))
		}
		if err := messaging.Init(ctx); err != nil {
			log.Fatal("failed to initialize messaging", zap.Error(err))
		}
		// Checks subscription, will os.Exit(1) if there is an error
		// Disabled for now, system service is the only one that validates subscription
		// ctx = subscription.Monitor(ctx)

		log.Info("Starting http server on address " + flags.http.Addr)
		listener, err := net.Listen("tcp", flags.http.Addr)
		if err != nil {
			log.Info("Can't listen on addr " + flags.http.Addr)
		}

		if flags.monitor.Interval > 0 {
			go metrics.NewMonitor(flags.monitor.Interval)
		}

		r := chi.NewRouter()

		// logging, cors and such
		middleware.Mount(ctx, r, flags.http)

		// Use JWT secret for hmac signer for now
		auth.DefaultSigner = auth.HmacSigner(authJwtFlags.Secret)
		auth.DefaultJwtHandler, err = auth.JWT(authJwtFlags.Secret, authJwtFlags.Expiry)
		if err != nil {
			log.Fatal("Error creating JWT Auth", zap.Error(err))
		}

		r.Route("/api", func(r chi.Router) {
			r.Route("/compose", func(r chi.Router) {
				compose.MountRoutes(ctx, r)
			})
			r.Route("/messaging", func(r chi.Router) {
				messaging.MountRoutes(ctx, r)
			})
			r.Route("/system", func(r chi.Router) {
				system.MountRoutes(ctx, r)
			})
			middleware.MountSystemRoutes(ctx, r, flags.http)
		})

		fileserver := http.FileServer(http.Dir("webapp"))

		for _, service := range []string{"admin", "auth", "messaging", "compose"} {
			r.HandleFunc("/"+service+"*", serveIndex("webapp", "compose/index.html", fileserver))
		}
		r.HandleFunc("/*", serveIndex("webapp", "index.html", fileserver))

		go http.Serve(listener, r)
		<-ctx.Done()
	}
}
