package main

import (
	"log"
	"os"

	"github.com/crusttech/crust/system"
	systemService "github.com/crusttech/crust/system/service"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rbac"
)

func main() {
	flags("system", service.Flags, auth.Flags, rbac.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	service.InitDatabase()
	systemService.Init()

	var commands []string
	if len(os.Args) > 0 {
		// @todo migrate to a proper solution (eg: https://github.com/spf13/cobra)
		commands = os.Args[1:]
		for a, arg := range os.Args {
			if arg == "--" && a+1 < len(os.Args) {
				commands = os.Args[a+1:]
			}
		}
	}

	cliRouter(commands...)
}

func cliRouter(commands ...string) {
	if len(commands) == 0 {
		return
	}

	switch commands[0] {
	case "users":
		cliExecUsers(commands[1:]...)
	default:
	}
}
