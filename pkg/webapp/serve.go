package webapp

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"strings"
)

func MakeWebappServer(opt options.HTTPServerOpt, fed options.FederationOpt) func(r chi.Router) {
	// Serves static files directly from FS
	return func(r chi.Router) {
		fileserver := http.FileServer(http.Dir(opt.WebappBaseDir))

		for _, app := range strings.Split(opt.WebappList, ",") {
			basedir := path.Join(opt.WebappBaseUrl, app)
			serveConfig(r, basedir, opt.ApiBaseUrl, fed.Enabled)
			r.Get(basedir+"*", serveIndex(opt.WebappBaseDir, basedir, fileserver))
		}

		serveConfig(r, opt.WebappBaseUrl, opt.ApiBaseUrl, fed.Enabled)
		r.Get(opt.WebappBaseUrl+"*", serveIndex(opt.WebappBaseDir, opt.WebappBaseUrl, fileserver))
	}
}

// Serves index.html in case the requested file isn't found (or some other os.Stat error)
func serveIndex(assetPath string, indexPath string, serve http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage := path.Join(assetPath, indexPath, "index.html")
		requestedPage := path.Join(assetPath, r.URL.Path)

		_, err := os.Stat(requestedPage)

		if err == nil {
			if strings.HasSuffix(r.URL.String(), "/") {
				// If request ends with a slash we want to prevent
				// serving of directory index (list of fileS)
				http.ServeFile(w, r, indexPage)
				return
			}

			// Serve the file requested
			serve.ServeHTTP(w, r)
			return
		}

		if os.IsNotExist(err) {
			// Forcefully serve index page on whatever
			http.ServeFile(w, r, indexPage)
			return
		}

		logger.Default().WithOptions(zap.AddStacktrace(zap.PanicLevel)).Error(
			"failed to serve static file",
			zap.Error(err),
			zap.Stringer("url", r.URL),
		)

		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func serveConfig(r chi.Router, appUrl, apiBaseUrl string, fedEnabled bool) {
	r.Get(strings.TrimRight(appUrl, "/")+"/config.js", func(w http.ResponseWriter, r *http.Request) {
		const line = "window.%sAPI = '%s/%s';\n"
		_, _ = fmt.Fprintf(w, line, "System", apiBaseUrl, "system")
		_, _ = fmt.Fprintf(w, line, "Messaging", apiBaseUrl, "messaging")
		_, _ = fmt.Fprintf(w, line, "Compose", apiBaseUrl, "compose")
		_, _ = fmt.Fprintf(w, line, "Automation", apiBaseUrl, "automation")

		if fedEnabled {
			_, _ = fmt.Fprintf(w, line, "Federation", apiBaseUrl, "federation")
		}
	})
}
