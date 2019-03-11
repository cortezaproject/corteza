package main

import (
	"log"
	"os"

	context "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rbac"
	"github.com/crusttech/crust/internal/subscription"
	"github.com/crusttech/crust/internal/version"
	sub "github.com/crusttech/crust/system"
)

func main() {
	log.Printf("Starting "+os.Args[0]+", version: %v, built on: %v", version.Version, version.BuildTime)

	ctx := context.AsContext(sigctx.New())

	flags(
		"system",
		sub.Flags,
		auth.Flags,
		rbac.Flags,
		subscription.Flags,
	)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := sub.Init(); err != nil {
		log.Fatalf("Error initializing: %+v", err)
	}

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "help":
	default:
		// Checks subscription, will os.Exit(1) if there is an error
		ctx = subscription.Monitor(ctx)

		if err := sub.StartRestAPI(ctx); err != nil {
			log.Fatalf("Error starting/running: %+v", err)
		}
	}
}
