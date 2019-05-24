package main

import (
	"fmt"
	"os"

	"github.com/cortezaproject/corteza-server/monolith"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

func main() {
	c := monolith.InitMonolith()
	if err := c.Command(cli.Context()).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
