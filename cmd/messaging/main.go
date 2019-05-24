package main

import (
	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

func main() {
	cfg := messaging.Configure()
	cmd := cfg.MakeCLI(cli.Context())
	cli.HandleError(cmd.Execute())
}
