package webapp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var (
	baseHrefMatcher = regexp.MustCompile(`<base\s+href="?.+?"?\s*\/?>`)
)

func MakeWebappServer(log *zap.Logger, httpSrvOpt options.HTTPServerOpt, authOpt options.AuthOpt, discoveryOpt options.DiscoveryOpt) func(r chi.Router) {
	var (
		apiBaseUrl          = options.CleanBase(httpSrvOpt.BaseUrl, httpSrvOpt.ApiBaseUrl)
		discoveryApiBaseUrl = options.CleanBase(discoveryOpt.BaseUrl)
		apps                = strings.Split(httpSrvOpt.WebappList, ",")

		appIndexHTMLs = make(map[string][]byte)

		webBaseUrl string
		err        error
	)

	// Preload index files for all apps
	for _, app := range append(apps, "") {
		appIndexHTMLs[app], err = modifyIndexHTML(app, httpSrvOpt.WebappBaseDir, httpSrvOpt.BaseUrl)
		if err != nil {
			log.Error("could not preload application index HTML", zap.Error(err))
		}
	}

	// Serves static files directly from FS
	return func(r chi.Router) {
		fs := http.StripPrefix(
			options.CleanBase(httpSrvOpt.BaseUrl, httpSrvOpt.WebappBaseUrl),
			http.FileServer(http.Dir(httpSrvOpt.WebappBaseDir)),
		)

		for _, app := range apps {
			webBaseUrl = options.CleanBase(httpSrvOpt.WebappBaseUrl, app)
			serveConfig(r, webBaseUrl, apiBaseUrl, authOpt.BaseURL, httpSrvOpt.BaseUrl, discoveryApiBaseUrl)
			r.Get(webBaseUrl+"*", serveIndex(httpSrvOpt, appIndexHTMLs[app], fs))
		}

		webBaseUrl = options.CleanBase(httpSrvOpt.WebappBaseUrl)
		serveConfig(r, webBaseUrl, apiBaseUrl, authOpt.BaseURL, httpSrvOpt.BaseUrl, discoveryApiBaseUrl)
		r.Get(webBaseUrl+"*", serveIndex(httpSrvOpt, appIndexHTMLs[""], fs))
	}
}

// Serves index.html in case the requested file isn't found (or some other os.Stat error)
func serveIndex(opt options.HTTPServerOpt, indexHTML []byte, serve http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestedFile := path.Join(
			// could be relative path
			opt.WebappBaseDir,
			strings.TrimPrefix(r.URL.Path, opt.BaseUrl),
		)

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

func serveConfig(r chi.Router, appUrl, apiBaseUrl, authBaseUrl, webappBaseUrl, discoveryApiBaseUrl string) {
	r.Get(options.CleanBase(appUrl, "config.js"), func(w http.ResponseWriter, r *http.Request) {
		const line = "window.%s = '%s';\n"
		_, _ = fmt.Fprintf(w, line, "CortezaAPI", apiBaseUrl)
		_, _ = fmt.Fprintf(w, line, "CortezaAuth", authBaseUrl)
		_, _ = fmt.Fprintf(w, line, "CortezaWebapp", webappBaseUrl)
		_, _ = fmt.Fprintf(w, line, "CortezaDiscoveryAPI", discoveryApiBaseUrl)
	})
}

// Reads and modifies index HTML for the webapp
//
// It replaces value for href in <base> tag with the virtual location of the web app
// This is controlled with HTTP_WEBAPP_BASE_URL
func modifyIndexHTML(app, dir, baseHref string) ([]byte, error) {
	if fh, err := os.Open(path.Join(dir, app, "index.html")); err != nil {
		return nil, err
	} else if buf, err := io.ReadAll(fh); err != nil {
		return nil, err
	} else {
		return replaceBaseHrefPlaceholder(buf, app, baseHref), nil
	}
}

func replaceBaseHrefPlaceholder(buf []byte, app, baseHref string) []byte {
	var (
		base = strings.TrimSuffix(options.CleanBase(baseHref, app), "/") + "/"

		warning     = []byte(`<!-- Error! Could not locate or modify <base> tag, your webapp might misbehave -->`)
		replacement = []byte(fmt.Sprintf(`<base href="%s" />`, base))
		fixed       = baseHrefMatcher.ReplaceAll(buf, replacement)
	)

	if bytes.Equal(buf, fixed) {
		return append(buf, warning...)
	}

	return fixed
}
