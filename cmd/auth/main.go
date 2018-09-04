package main

import (
	"log"
	"os"

	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/rbac"
)

func main() {
	config := flags("auth", rbac.Flags, auth.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go NewMonitor(config.monitorInterval)

	if err := auth.Init(); err != nil {
		log.Fatalf("Error initializing auth: %+v", err)
	}
	if err := auth.Start(); err != nil {
		log.Fatalf("Error starting/running auth: %+v", err)
	}
}
