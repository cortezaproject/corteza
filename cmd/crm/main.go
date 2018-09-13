package main

import (
	"log"
	"os"

	service "github.com/crusttech/crust/crm"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rbac"
)

func main() {
	config := flags("crm", service.Flags, auth.Flags, rbac.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go NewMonitor(config.monitorInterval)

	if err := service.Init(); err != nil {
		log.Fatalf("Error initializing crm: %+v", err)
	}
	if err := service.Start(); err != nil {
		log.Fatalf("Error starting/running crm: %+v", err)
	}
}
