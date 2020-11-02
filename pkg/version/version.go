package version

import (
	"time"
)

var (
	// BuildTime value is set at build time and served over API and CLI
	BuildTime = time.Now().Format(time.RFC3339)

	// Version is set as LDFLAG at build time:
	// -X github.com/cortezaproject/corteza-server/pkg/version.Version=....
	// See Makefile for details
	Version = "development"
)
