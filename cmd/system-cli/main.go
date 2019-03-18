package main

import (
	"log"
	"os"

	"github.com/crusttech/crust/internal/auth"
	system "github.com/crusttech/crust/system"
	systemService "github.com/crusttech/crust/system/service"
)

func main() {
	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flags("system", system.Flags, auth.Flags)

	system.InitDatabase()
	systemService.Init()

	setupCobra()
}
