# httperr

Package httperr provides a golang structure Error which implements the
error interface and captures an HTTP status code and a private error. 

    // Error represents an error that can be modeled as an
    // http status code.
    type Error struct {
        StatusCode   int    // If not supplied, http.StatusInternalServerError is used.
        Status       string // If not supplied, http.StatusText(StatusCode) is used.
        PrivateError error  // An additional error that is not displayed to the user, but may be logged.
    }

It transforms code full of redundant log/error calls like this:

    func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
        remoteUser, err := s.Auth.RequireUser(w, r)
        if err != nil {
            http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
            return
        }

        user, err := s.Storage.Get(remoteUser.Name)
        if err != nil {
            log.Printf("ERROR: cannot fetch user: %s", err)
            http.Error(w, http.StatusText(http.StatusInternalServerErorr), http.StatusInternalServerErorr)
            return
        }
        json.NewEncoder(w).Encode(user)
    }

Into this:

    func (s *Server) getUser(w http.ResponseWriter, r *http.Request) error {
        remoteUser, err := s.Auth.RequireUser(w, r)
        if err != nil {
            return httperr.Unauthorized
        }

        user, err := s.Storage.Get(remoteUser.Name)
        if err != nil {
            return err
        }
        json.NewEncoder(w).Encode(user)
        return nil
    }

Life changing? Probably not, but it seems to remove a lot of redundancy and make control flow in web servers simpler.

