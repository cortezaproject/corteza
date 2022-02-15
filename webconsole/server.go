package webconsole

import (
	"embed"
	"io/fs"

	"github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/go-chi/chi/v5"
)

var (
	//go:embed dist
	assets embed.FS
)

func Mount(r chi.Router) error {
	assets, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}

	return http.MountSPA(r, "/ui", assets, http.UrlPrefix("/console/ui"))
}
