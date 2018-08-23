package main

import (
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/rbac"

	"github.com/crusttech/crust/sam"

	"log"
	"os"
)

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func main() {
	config := flags("sam", sam.Flags, auth.Flags, rbac.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go NewMonitor(config.monitorInterval)

	if err := sam.Init(); err != nil {
		log.Fatalf("Error initializing sam: %+v", err)
	}
	if err := sam.Start(); err != nil {
		log.Fatalf("Error starting/running sam: %+v", err)
	}
}
