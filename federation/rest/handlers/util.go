package handlers

import (
	"net/http"
)

// Edit this file in `codegen/common/rest/handlers/util.go`;
// any changes under [service] will be overwritten by codegen.

func serveHTTP(value interface{}, w http.ResponseWriter, r *http.Request) bool {
	switch fn := value.(type) {
	case func(http.ResponseWriter, *http.Request):
		fn(w, r)
		return true
	}
	return false
}
