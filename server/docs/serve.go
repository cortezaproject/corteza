package docs

import (
	"embed"
	"net/http"
)

//go:embed *.yaml
//go:embed *.html
var docs embed.FS

func GetFS() http.FileSystem {
	return http.FS(docs)
}
