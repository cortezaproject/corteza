package main

import (
	"fmt"
	"os"

	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

func main() {
	c := compose.InitCompose()
	if err := c.Command(cli.Context()).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
