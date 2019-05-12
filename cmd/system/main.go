package main

import (
	"os"

	context "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/namsral/flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/subscription"
	system "github.com/crusttech/crust/system"
)

func main() {

	// Initialize default logger
	logger.Init(zapcore.DebugLevel)
	log := logger.Default().Named("system")

	// New signal-bond context that we will use and
	// will get terminated (Done()) on SIGINT or SIGTERM
	ctx := context.AsContext(sigctx.New())

	// Bind default logger to context
	ctx = logger.ContextWithValue(ctx, log)

	system.Flags("system")

	subscription.Flags()

	flag.Parse()

	if err := system.Init(ctx); err != nil {
		log.Fatal("failed to initialize system", zap.Error(err))
	}

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "help":
		flag.PrintDefaults()
	case "provision":
		if err := system.Provision(ctx); err != nil {
			println("Failed to provision system: ", err.Error())
			os.Exit(1)
		}
	default:
		// Checks subscription, will os.Exit(1) if there is an error
		ctx = subscription.Monitor(ctx)

		if err := system.StartRestAPI(ctx); err != nil {
			log.Fatal("failed to start system REST API", zap.Error(err))
		}
	}
}
