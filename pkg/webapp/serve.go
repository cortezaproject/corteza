package webapp

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/go-chi/chi"
	"net/http"
	"os"
	"path"
	"strings"
)

func MakeWebappServer(opt options.HTTPServerOpt) func(r chi.Router) {
	// Serves static files directly from FS
	return func(r chi.Router) {
		fileserver := http.FileServer(http.Dir(opt.WebappBaseDir))

		for _, app := range strings.Split(opt.WebappList, ",") {
			basedir := path.Join(opt.WebappBaseUrl, app)
			serveConfig(r, basedir, opt.ApiBaseUrl)
			r.Get(basedir+"*", serveIndex(opt.WebappBaseDir, basedir+"/index.html", fileserver))
		}

		serveConfig(r, opt.WebappBaseUrl, opt.ApiBaseUrl)
		r.Get(opt.WebappBaseUrl+"*", serveIndex(opt.WebappBaseDir, opt.WebappBaseUrl+"/index.html", fileserver))
	}
}

// Serves index.html in case the requested file isn't found (or some other os.Stat error)
func serveIndex(assetPath string, indexPath string, serve http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage := path.Join(assetPath, indexPath)
		requestedPage := path.Join(assetPath, r.URL.Path)
		_, err := os.Stat(requestedPage)

		if err != nil {
			http.ServeFile(w, r, indexPage)
			return
		}
		serve.ServeHTTP(w, r)
	}
}

func serveConfig(r chi.Router, appUrl, apiBaseUrl string) {
	r.Get(strings.TrimRight(appUrl, "/")+"/config.js", func(w http.ResponseWriter, r *http.Request) {
		const line = "window.%sAPI = '%s/%s';\n"
		_, _ = fmt.Fprintf(w, line, "System", apiBaseUrl, "system")
		_, _ = fmt.Fprintf(w, line, "Messaging", apiBaseUrl, "messaging")
		_, _ = fmt.Fprintf(w, line, "Compose", apiBaseUrl, "compose")
	})
}
