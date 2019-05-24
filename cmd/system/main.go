package main

import (
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system"
)

func main() {
	cfg := system.Configure()
	cmd := cfg.MakeCLI(cli.Context())
	cli.HandleError(cmd.Execute())
}
