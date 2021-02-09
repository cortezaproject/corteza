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

func MakeWebappServer(httpSrvOpt options.HTTPServerOpt, authOpt options.AuthOpt, fed options.FederationOpt) func(r chi.Router) {
	// Serves static files directly from FS
	return func(r chi.Router) {
		fileserver := http.FileServer(http.Dir(httpSrvOpt.WebappBaseDir))

		for _, app := range strings.Split(httpSrvOpt.WebappList, ",") {
			basedir := path.Join(httpSrvOpt.WebappBaseUrl, app)
			serveConfig(r, basedir, httpSrvOpt.ApiBaseUrl, authOpt.BaseURL, fed.Enabled)
			r.Get(basedir+"*", serveIndex(httpSrvOpt.WebappBaseDir, basedir, fileserver))
		}

		serveConfig(r, httpSrvOpt.WebappBaseUrl, httpSrvOpt.ApiBaseUrl, authOpt.BaseURL, fed.Enabled)
		r.Get(httpSrvOpt.WebappBaseUrl+"*", serveIndex(httpSrvOpt.WebappBaseDir, httpSrvOpt.WebappBaseUrl, fileserver))
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

func serveConfig(r chi.Router, appUrl, apiBaseUrl, authBaseUrl string, fedEnabled bool) {
	r.Get(strings.TrimRight(appUrl, "/")+"/config.js", func(w http.ResponseWriter, r *http.Request) {
		const line = "window.%s = '%s';\n"
		_, _ = fmt.Fprintf(w, line, "CortezaAPI", apiBaseUrl)
		_, _ = fmt.Fprintf(w, line, "CortezaAuth", authBaseUrl)
	})
}
