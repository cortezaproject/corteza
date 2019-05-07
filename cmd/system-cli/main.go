package main

import (
	context "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/crusttech/crust/internal/logger"
	system "github.com/crusttech/crust/system"
	"github.com/crusttech/crust/system/cli"
)

func main() {

	// Initialize default logger
	logger.Init(zapcore.DebugLevel)
	log := logger.Default().Named("system-cli")

	// New signal-bond context that we will use and
	// will get terminated (Done()) on SIGINT or SIGTERM
	ctx := context.AsContext(sigctx.New())

	// Bind default logger to context
	ctx = logger.ContextWithValue(ctx, log)

	flags("system", system.Flags)
	if err := system.Init(ctx); err != nil {
		log.Fatal("failed to initialize system", zap.Error(err))
	}

	cli.StartCLI(ctx)
}
