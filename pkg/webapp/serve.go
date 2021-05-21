package webapp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var (
	htmlIndex struct {
		body []byte
		mod  time.Time
	}
)

func MakeWebappServer(log *zap.Logger, httpSrvOpt options.HTTPServerOpt, authOpt options.AuthOpt) func(r chi.Router) {
	var (
		apiBaseUrl = options.CleanBase(httpSrvOpt.BaseUrl, httpSrvOpt.ApiBaseUrl)
		apps       = strings.Split(httpSrvOpt.WebappList, ",")

		appIndexHTMLs = make(map[string][]byte)

		webBaseUrl string
		err        error
	)

	// Preload index files for all apps
	for _, app := range append(apps, "") {
		pathPrefix := path.Join(httpSrvOpt.BaseUrl, httpSrvOpt.WebappBaseUrl, app)
		if !strings.HasSuffix(pathPrefix, "/") {
			pathPrefix += "/"
		}

		appIndexHTMLs[app], err = modifyIndexHTML(path.Join(httpSrvOpt.WebappBaseDir, app), pathPrefix)
		if err != nil {
			log.Error("could not preload application index HTML", zap.Error(err))
		}
	}

	// Serves static files directly from FS
	return func(r chi.Router) {
		fileserver := http.StripPrefix(
			path.Join(httpSrvOpt.BaseUrl, httpSrvOpt.WebappBaseUrl),
			http.FileServer(http.Dir(httpSrvOpt.WebappBaseDir)),
		)

		for _, app := range apps {
			webBaseUrl = "/" + path.Join(httpSrvOpt.WebappBaseUrl, app)
			serveConfig(r, webBaseUrl, apiBaseUrl, authOpt.BaseURL, httpSrvOpt.BaseUrl)

			r.Get(webBaseUrl+"*", serveIndex(httpSrvOpt, appIndexHTMLs[app], fileserver))
		}

		webBaseUrl = "/" + path.Join(httpSrvOpt.WebappBaseUrl)
		serveConfig(r, webBaseUrl, apiBaseUrl, authOpt.BaseURL, httpSrvOpt.BaseUrl)

		r.Get(webBaseUrl+"*", serveIndex(httpSrvOpt, appIndexHTMLs[""], fileserver))
	}
}

// Serves index.html in case the requested file isn't found (or some other os.Stat error)
func serveIndex(opt options.HTTPServerOpt, indexHTML []byte, serve http.Handler) http.HandlerFunc {
	//indexPage := path.Join(opt.WebappBaseDir, indexPath, "index.html")

	return func(w http.ResponseWriter, r *http.Request) {
		requestedFile := path.Join(opt.WebappBaseDir, strings.TrimPrefix(r.URL.Path, opt.BaseUrl))

		f, err := os.Stat(requestedFile)
		// When file does not exist or is a directory, serve app's index
		if os.IsNotExist(err) || f.IsDir() || strings.HasSuffix(r.URL.String(), "/") {
			// Make sure index is not cached
			// In the big scheme of things, this couple of kilobytes does not make any difference
			// and it's really important that users get fresh index files on each full refresh
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(indexHTML)
			return
		} else if err == nil {
			// Serve the file requested
			serve.ServeHTTP(w, r)
			return
		}

		logger.Default().Error(
			"failed to serve static file",
			zap.Error(err),
			zap.Stringer("url", r.URL),
		)

		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func serveConfig(r chi.Router, appUrl, apiBaseUrl, authBaseUrl, webappBaseUrl string) {
	r.Get(strings.TrimSuffix(appUrl, "/")+"/config.js", func(w http.ResponseWriter, r *http.Request) {
		const line = "window.%s = '%s';\n"
		_, _ = fmt.Fprintf(w, line, "CortezaAPI", apiBaseUrl)
		_, _ = fmt.Fprintf(w, line, "CortezaAuth", authBaseUrl)
		_, _ = fmt.Fprintf(w, line, "CortezaWebapp", webappBaseUrl)
	})
}

// Reads and modifies index HTML for the webapp
//
// It replaces <base> tag with the exact location of the web app
func modifyIndexHTML(dir, baseHref string) (buf []byte, err error) {
	var (
		warning     = []byte("\n\n<!--\n\nError!\n\nFailed could not locate or modify <base> tag, your webapp might misbehave\n\n-->\n")
		placeholder = []byte("<base href=/ >")
		replacement = []byte(fmt.Sprintf(`<base href="%s" />`, baseHref))
		fh          *os.File
	)

	fh, err = os.Open(path.Join(dir, "index.html"))
	if err != nil {
		return
	}

	buf, err = io.ReadAll(fh)
	if err != nil {
		return
	}

	if bytes.Contains(buf, placeholder) {
		return bytes.Replace(buf, placeholder, replacement, 1), nil
	} else {
		return append(buf, warning...), nil
	}
}
