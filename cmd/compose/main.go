package main

import (
	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

func main() {
	cfg := compose.Configure()
	cmd := cfg.MakeCLI(cli.Context())
	cli.HandleError(cmd.Execute())
}
