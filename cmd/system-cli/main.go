package main

import (
	"log"
	"os"

	context "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"

	"github.com/crusttech/crust/internal/auth"
	system "github.com/crusttech/crust/system"
)

func main() {
	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flags("system", system.Flags, auth.Flags)

	system.Init(context.AsContext(sigctx.New()))

	setupCobra()
}
