package main

import (
	"github.com/cortezaproject/corteza-server/corteza"
	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func main() {
	logger.Init()

	app.Run(
		logger.Default(),
		app.NewOptions(messaging.SERVICE),
		&corteza.App{},
		&messaging.App{},
	)
}
