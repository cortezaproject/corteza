package basicauth

import (
	"fmt"
	"net/http"
)

// New returns a piece of middleware that will allow access only
// if the provided credentials match within the given service
// otherwise it will return a 401 and not call the next handler.
func New(realm string, credentials map[string][]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				unauthorized(w, realm)
				return
			}

			validPasswords, userFound := credentials[username]
			if !userFound {
				unauthorized(w, realm)
				return
			}

			for _, validPassword := range validPasswords {
				if password == validPassword {
					next.ServeHTTP(w, r)
					return
				}
			}

			unauthorized(w, realm)
		})
	}
}

func unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}
