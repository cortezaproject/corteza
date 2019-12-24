package main

import (
	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/corteza"
	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func main() {
	logger.Init()

	app.Run(
		logger.Default(),
		app.NewOptions(compose.SERVICE),
		&corteza.App{},
		&compose.App{},
	)
}
