package request

import (
	"net/http"
)

// RequestFiller is an interface for typed request parameters
type RequestFiller interface {
	Fill(r *http.Request) error
}
