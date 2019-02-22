package main

import (
	"log"
	"os"

	context "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"

	"github.com/crusttech/crust/internal/subscription"
	service "github.com/crusttech/crust/system"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rbac"
)

func main() {
	ctx := context.AsContext(sigctx.New())

	flags(
		"system",
		service.Flags,
		auth.Flags,
		rbac.Flags,
		subscription.Flags,
	)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := service.Init(); err != nil {
		log.Fatalf("Error initializing: %+v", err)
	}

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "help":
	case "merge-users":
	default:
		log.Println("Validating subscription")
		// Checks subscription & runs internal checker that runs every 24h
		if err := subscription.Check(ctx); err != nil {
			log.Printf("Subscription could not be validated, reason: %v", err)
			os.Exit(-1)
		} else {
			log.Println("Subscription valdiated")
		}

		if err := service.StartRestAPI(ctx); err != nil {
			log.Fatalf("Error starting/running: %+v", err)
		}
	}
}
