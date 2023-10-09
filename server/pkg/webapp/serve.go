package webapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/bep/godartsass/v2"
	"github.com/cortezaproject/corteza/server/assets"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	webappConfig struct {
		appUrl              string
		apiBaseUrl          string
		authBaseUrl         string
		webappBaseUrl       string
		discoveryApiBaseUrl string
		sentryUrl           string
		settings            *types.AppSettings
	}
)

var (
	baseHrefMatcher = regexp.MustCompile(`<base\s+href="?.+?"?\s*\/?>`)
)

func MakeWebappServer(log *zap.Logger, httpSrvOpt options.HttpServerOpt, authOpt options.AuthOpt, discoveryOpt options.DiscoveryOpt, sentryOpt options.SentryOpt) func(r chi.Router) {
	var (
		apiBaseUrl          = options.CleanBase(httpSrvOpt.BaseUrl, httpSrvOpt.ApiBaseUrl)
		webappSentryUrl     = sentryOpt.WebappDSN
		discoveryApiBaseUrl = discoveryOpt.BaseUrl
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
			serveConfig(r, webappConfig{
				appUrl:              webBaseUrl,
				apiBaseUrl:          apiBaseUrl,
				authBaseUrl:         authOpt.BaseURL,
				webappBaseUrl:       httpSrvOpt.BaseUrl,
				discoveryApiBaseUrl: discoveryApiBaseUrl,
				sentryUrl:           webappSentryUrl,
				settings:            service.CurrentSettings,
			})
			r.Get(webBaseUrl+"*", serveIndex(httpSrvOpt, appIndexHTMLs[app], fs))
		}

		webBaseUrl = options.CleanBase(httpSrvOpt.WebappBaseUrl)
		serveConfig(r, webappConfig{
			appUrl:              webBaseUrl,
			apiBaseUrl:          apiBaseUrl,
			authBaseUrl:         authOpt.BaseURL,
			webappBaseUrl:       httpSrvOpt.BaseUrl,
			discoveryApiBaseUrl: discoveryApiBaseUrl,
			sentryUrl:           webappSentryUrl,
			settings:            service.CurrentSettings,
		})
		r.Get(webBaseUrl+"*", serveIndex(httpSrvOpt, appIndexHTMLs[""], fs))
	}
}

// Serves index.html in case the requested file isn't found (or some other os.Stat error)
func serveIndex(opt options.HttpServerOpt, indexHTML []byte, serve http.Handler) http.HandlerFunc {
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

func serveConfig(r chi.Router, config webappConfig) {
	r.Get(options.CleanBase(config.appUrl, "config.js"), func(w http.ResponseWriter, r *http.Request) {

		// Assure the content-type
		// The presence of the X-Content-Type-Options: nosniff header breaks web applications
		w.Header().Add("Content-Type", "text/javascript")

		const line = "window.%s = '%s';\n"
		_, _ = fmt.Fprintf(w, line, "CortezaAPI", config.apiBaseUrl)
		_, _ = fmt.Fprintf(w, line, "CortezaAuth", config.authBaseUrl)
		_, _ = fmt.Fprintf(w, line, "CortezaWebapp", config.webappBaseUrl)
		if len(config.discoveryApiBaseUrl) > 0 {
			_, _ = fmt.Fprintf(w, line, "CortezaDiscoveryAPI", config.discoveryApiBaseUrl)
		}
		if len(config.sentryUrl) > 0 {
			_, _ = fmt.Fprintf(w, line, "SentryDSN", config.sentryUrl)
		}
	})

	r.Get(options.CleanBase(config.appUrl, "custom.css"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")

		var stringsBuilder strings.Builder

		assetFiles := assets.Files(logger.Default(), "")

		// process BrandingSASS
		if config.settings.UI.BrandingSASS != "" {
			brandingSassString, err := jsonToSass(config.settings.UI.BrandingSASS)
			if err != nil {
				logger.Default().Error("failed to parse branding SASS", zap.Error(err))
				return
			}
			stringsBuilder.WriteString(brandingSassString)
		}

		// Get SASS variables from the custom CSS editor and give them precedence over branding variables
		variablesPattern := regexp.MustCompile(`(\$[a-zA-Z_-]+):\s*([^;]+);`)
		customVariables := variablesPattern.FindAllString(config.settings.UI.CustomCSS, -1)

		for _, customVariable := range customVariables {
			stringsBuilder.WriteString(fmt.Sprintf("%s \n", customVariable))
		}

		// Boostrap, bootstrap-vue and custom variables scss files
		fileNames := assets.FileNames("scss")
		for _, fileName := range fileNames {
			open, err := assetFiles.Open("scss/" + fileName)
			if err != nil {
				logger.Default().Error(fmt.Sprintf("failed to open asset %s file", fileName), zap.Error(err))
				return
			}
			fileContent, err := io.ReadAll(open)
			if err != nil {
				return
			}
			stringsBuilder.WriteString(string(fileContent))
		}

		//Custom CSS editor selector block
		selectorBlock := variablesPattern.ReplaceAllString(config.settings.UI.CustomCSS, "")
		stringsBuilder.WriteString(selectorBlock)

		// Process Sass
		args := godartsass.Args{
			Source: stringsBuilder.String(),
		}

		t, err := godartsass.Start(godartsass.Options{
			DartSassEmbeddedFilename: "sass",
		})

		if err != nil {
			panic(err.Error())
		}

		execute, err := t.Execute(args)
		if err != nil {
			logger.Default().Error("Sass syntax error", zap.Error(err))
			return
		}

		_, _ = fmt.Fprint(w, execute.CSS)
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

// jsonToSass converts JSON string to SASS variable assignment string
func jsonToSass(jsonStr string) (string, error) {
	var (
		colorMap       map[string]string
		stringsBuilder strings.Builder
	)

	err := json.Unmarshal([]byte(jsonStr), &colorMap)
	if err != nil {
		return "", err
	}

	for key, value := range colorMap {
		stringsBuilder.WriteString(fmt.Sprintf("$%s: %s;\n", key, value))
	}

	return stringsBuilder.String(), nil
}
