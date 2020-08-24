package version

import (
	"net/http"
	"time"

	"github.com/titpetric/factory/resputil"
)

var (
	// BuildTime value is set at build time and served over API and CLI
	BuildTime = time.Now().Format(time.RFC3339)

	// Version is set as LDFLAG at build time:
	// -X github.com/cortezaproject/corteza-server/pkg/version.Version=....
	// See Makefile for details
	Version = "development"
)

func HttpHandler(w http.ResponseWriter, r *http.Request) {
	resputil.JSON(w, struct {
		BuildTime string `json:"buildTime"`
		Version   string `json:"version"`
	}{BuildTime, Version})
}
