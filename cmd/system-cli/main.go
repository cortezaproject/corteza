package main

import (
	"log"
	"os"

	"github.com/crusttech/crust/internal/auth"
	system "github.com/crusttech/crust/system"
)

func main() {
	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flags("system", system.Flags, auth.Flags)

	system.Init()

	setupCobra()
}
