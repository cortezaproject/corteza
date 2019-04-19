package main

import (
	"log"
	"os"

	context "github.com/SentimensRG/ctx"
	"github.com/SentimensRG/ctx/sigctx"

	system "github.com/crusttech/crust/system"
	"github.com/crusttech/crust/system/cli"
)

func main() {
	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.AsContext(sigctx.New())

	flags("system", system.Flags)
	if err := system.Init(ctx); err != nil {
		log.Fatalf("Error initializing system: %+v", err)
	}

	cli.Run(ctx)
}
