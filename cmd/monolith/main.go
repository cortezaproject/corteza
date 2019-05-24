package main

import (
	"github.com/cortezaproject/corteza-server/monolith"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

func main() {
	cfg := monolith.Configure()
	cmd := cfg.MakeCLI(cli.Context())
	cli.HandleError(cmd.Execute())
}
