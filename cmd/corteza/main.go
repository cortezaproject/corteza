package main

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func main() {
	// Initialize logger before any other action
	logger.Init()

	cli.HandleError(app.New().Execute())
}
