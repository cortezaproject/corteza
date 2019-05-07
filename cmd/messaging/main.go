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
	messaging "github.com/crusttech/crust/messaging"
	system "github.com/crusttech/crust/system"
)

func main() {
	// Initialize default logger
	logger.Init(zapcore.DebugLevel)
	log := logger.Default().Named("messaging")

	// New signal-bond context that we will use and
	// will get terminated (Done()) on SIGINT or SIGTERM
	ctx := context.AsContext(sigctx.New())

	// Bind default logger to context
	ctx = logger.ContextWithValue(ctx, log)

	messaging.Flags("messaging")
	system.Flags("system")

	subscription.Flags()

	flag.Parse()

	if err := system.Init(ctx); err != nil {
		log.Fatal("failed to initialize system", zap.Error(err))
	}
	if err := messaging.Init(ctx); err != nil {
		log.Fatal("failed to initialize messaging", zap.Error(err))
	}

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "help":
		flag.PrintDefaults()
	default:
		// Checks subscription, will os.Exit(1) if there is an error
		// Disabled for now, system service is the only one that validates subscription
		// ctx = subscription.Monitor(ctx)

		if err := messaging.StartRestAPI(ctx); err != nil {
			log.Fatal("failed to start messaging REST API", zap.Error(err))
		}
	}
}
