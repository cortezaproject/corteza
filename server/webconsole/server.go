package webconsole

import (
	"embed"
	"io/fs"

	"github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/go-chi/chi/v5"
)

var (
	// we need combination of go:embed dist/* and .placeholder
	// file inside. If only dist (w/o) wildcard is used,
	// dot-file (.placeholder) will be ignored

	//go:embed dist/*
	assets embed.FS
)

func Mount(r chi.Router) error {
	assets, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}

	return http.MountSPA(r, "/ui", assets, http.UrlPrefix("/console/ui"))
}
