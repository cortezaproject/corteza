package http

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	// special file-server meant to be used for serving single page applications
	// or anything that needs a bit of extra attention like fallback handling
	// 404 handler, error handler and URL prefixing
	fileServer struct {
		// fileServer files
		files http.FileSystem

		urlPrefix string

		fallbacks []string

		// Final not-found handler
		notFound http.HandlerFunc

		// how errors are handled
		errHandler func(w http.ResponseWriter, error string, code int)

		logger *zap.Logger
	}

	configurator func(*fileServer) error
)

// MountSPA helper function, preconfigures FileServer for SPA serving
// and mounts it to chi Router
func MountSPA(r chi.Router, path string, root fs.FS, cc ...configurator) error {
	path = "/" + strings.Trim(strings.TrimRight(path, "*"), "/") + "/"
	handler, err := FileServer(root, append(cc, UrlPrefix(path), Fallbacks("index.html"))...)
	if err != nil {
		return err
	}

	r.Handle(
		strings.TrimRight(path, "/"),
		http.RedirectHandler("."+path, http.StatusTemporaryRedirect),
	)

	r.Handle(path+"*", handler)
	return nil
}

func FileServer(files fs.FS, cc ...configurator) (h *fileServer, err error) {
	h = &fileServer{
		files:      http.FS(files),
		notFound:   http.NotFound,
		errHandler: http.Error,
		logger:     zap.NewNop(),
	}

	for _, configure := range cc {
		if err = configure(h); err != nil {
			return
		}
	}

	return
}

func UrlPrefix(prefix string) configurator {
	return func(s *fileServer) error {
		s.urlPrefix = prefix
		return nil
	}
}

func Fallbacks(ff ...string) configurator {
	return func(s *fileServer) error {
		s.fallbacks = ff
		return nil
	}
}

func NotFound(h http.HandlerFunc) configurator {
	return func(s *fileServer) error {
		s.notFound = h
		return nil
	}
}

func Logger(l *zap.Logger) configurator {
	return func(s *fileServer) error {
		s.logger = l
		return nil
	}
}

// Serves the single-page-application
//
// This is file-server with some special logic for handling missing
// files (404s) and directories.
// In both cases we serve index file directly
func (h *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// catch requests for non-existing files and redirect to index.html
	if h.files == nil {
		h.errHandler(w, "configured without files", http.StatusInternalServerError)
		return
	}

	trimmed := path.Clean(strings.TrimPrefix(r.URL.Path, h.urlPrefix))
	h.logger.Debug(r.URL.Path, zap.String("trimmed", trimmed), zap.String("urlPrefix", h.urlPrefix))
	r.URL.Path = trimmed

	var (
		err error
		fh  http.File
		st  fs.FileInfo
	)

	for _, candidate := range append([]string{r.URL.Path}, h.fallbacks...) {
		if len(candidate) == 0 {
			continue
		}

		if fh, err = h.files.Open(candidate); err != nil {
			continue
		} else if st, err = fh.Stat(); err != nil {
			continue
		} else if st.IsDir() {
			// index
			continue
		}

		break
	}

	if fh == nil || st == nil {
		h.notFound(w, r)
		return
	}

	http.ServeContent(w, r, st.Name(), st.ModTime(), fh)
}
