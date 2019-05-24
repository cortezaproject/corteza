package main

import (
	"fmt"
	"os"

	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

func main() {
	m := messaging.InitMessaging()
	if err := m.Command(cli.Context()).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
