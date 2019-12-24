package main

import (
	"github.com/cortezaproject/corteza-server/corteza"
	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system"
)

func main() {
	logger.Init()

	app.Run(
		logger.Default(),
		app.NewOptions(system.SERVICE),
		&corteza.App{},
		&system.App{},
	)
}
