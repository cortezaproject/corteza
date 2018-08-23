package main

import (
	"log"
	"os"

	"github.com/crusttech/crust/crm"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/rbac"
)

func main() {
	config := flags("crm", crm.Flags, rbac.Flags, auth.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go NewMonitor(config.monitorInterval)

	if err := crm.Init(); err != nil {
		log.Fatalf("Error initializing crm: %+v", err)
	}
	if err := crm.Start(); err != nil {
		log.Fatalf("Error starting/running crm: %+v", err)
	}
}
