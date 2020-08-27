package server

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/docs"
	"github.com/goware/statik/fs"
	"net/http"
)

func serveDocs(w http.ResponseWriter, r *http.Request) {
	root, err := fs.New(docs.Asset)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read embeded filesystem: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	http.StripPrefix("/docs", http.FileServer(root)).ServeHTTP(w, r)
}
