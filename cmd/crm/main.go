package main

import (
	"log"
	"os"

	service "github.com/crusttech/crust/crm"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rbac"
)

func main() {
	flags("crm", service.Flags, auth.Flags, rbac.Flags)

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
		if err := service.Start(); err != nil {
			log.Fatalf("Error starting/running: %+v", err)
		}
	}
}
