package main

import (
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/pkg/cli"
	"github.com/crusttech/human/server/pkg/logger"
)

func main() {
	// Initialize logger before any other action
	logger.Init()

	cli.HandleError(app.New().Execute())
}
