package main

import (
	"fmt"
	"os"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system"
)

func main() {
	s := system.InitSystem()
	if err := s.Command(cli.Context()).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
