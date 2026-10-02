package main

import (
	// Embed the IANA timezone database into the binary so timezone-aware
	// features (e.g. record export) work in minimal deploy images that
	// don't ship /usr/share/zoneinfo.
	_ "time/tzdata"

	"github.com/cortezaproject/corteza/server/app"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/logger"
)

func main() {
	// Initialize logger before any other action
	logger.Init()

	cli.HandleError(app.New().Execute())
}
