package samlsp

import (
	"log"
	"net/http"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/logger"
)

// ErrorFunction is a callback that is invoked to return an error to the
// web user.
type ErrorFunction func(w http.ResponseWriter, r *http.Request, err error)

// DefaultOnError is the default ErrorFunction implementation. It prints
// an message via the standard log package and returns a simple text
// "Forbidden" message to the user.
func DefaultOnError(w http.ResponseWriter, r *http.Request, err error) {
	if parseErr, ok := err.(*saml.InvalidResponseError); ok {
		log.Printf("WARNING: received invalid saml response: %s (now: %s) %s",
			parseErr.Response, parseErr.Now, parseErr.PrivateErr)
	} else {
		log.Printf("ERROR: %s", err)
	}
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

// defaultOnErrorWithLogger is like DefaultOnError but accepts a custom logger.
// This is a bridge for backward compatibility with people use provide the
// deprecated Logger options field to New().
func defaultOnErrorWithLogger(log logger.Interface) ErrorFunction {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if parseErr, ok := err.(*saml.InvalidResponseError); ok {
			log.Printf("WARNING: received invalid saml response: %s (now: %s) %s",
				parseErr.Response, parseErr.Now, parseErr.PrivateErr)
		} else {
			log.Printf("ERROR: %s", err)
		}
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
}
