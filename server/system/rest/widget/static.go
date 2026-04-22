package widget

import (
	"net/http"
	"os"
	"path/filepath"
)

// WidgetJSHandler returns an http.Handler that serves the built widget.js
// bundle from `baseDir/chatbot/widget.js`. Missing file -> 404.
//
// Served with `Access-Control-Allow-Origin: *` because the third-party
// host page loads this as a classic <script src>.
func WidgetJSHandler(baseDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Join(baseDir, "chatbot", "widget.js")
		f, err := os.Open(p)
		if err != nil {
			http.Error(w, "widget.js not built", http.StatusNotFound)
			return
		}
		defer f.Close()

		st, err := f.Stat()
		if err != nil {
			http.Error(w, "widget.js unreadable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "public, max-age=300")
		http.ServeContent(w, r, "widget.js", st.ModTime(), f)
	})
}
