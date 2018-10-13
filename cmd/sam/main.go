package main

import (
	"log"
	"os"

	service "github.com/crusttech/crust/sam"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rbac"
	"github.com/crusttech/crust/internal/version"
)

func main() {
	config := flags("sam", service.Flags, auth.Flags, rbac.Flags)

	log.Printf("Starting sam, version: %v, built on: %v", version.Version, version.BuildTime)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go NewMonitor(config.monitorInterval)

	if err := service.Init(); err != nil {
		log.Fatalf("Error initializing sam: %+v", err)
	}
	if err := service.Start(); err != nil {
		log.Fatalf("Error starting/running sam: %+v", err)
	}
}
