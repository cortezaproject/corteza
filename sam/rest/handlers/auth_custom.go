package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `auth.go`, `auth.util.go` or `auth_test.go` to
	implement your API calls, helper functions and tests. The file `auth.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/sam/rest/request"
	"net/url"
	"time"
)

type (
	authPayload interface {
		Token() string
	}
)

// Initializies custom auth handler that attaches cookie info
//
// Cookie with JWT is added on successful login or user creation
//
func NewAuthCustom(ah AuthAPI, cookieExp int) *Auth {
	setCookie := func(w http.ResponseWriter, reqUrl *url.URL) func(payload interface{}, err error) (interface{}, error) {
		return func(payload interface{}, err error) (interface{}, error) {
			if ap, ok := payload.(authPayload); ok && err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:  "jwt",
					Value: ap.Token(),

					HttpOnly: false, // we need this for attachments & ws!
					Secure:   reqUrl.Scheme == "https",
					Path:     "/",
					//Domain:   "localhost",

					// @todo read from the config file.
					Expires: time.Now().Add(time.Duration(cookieExp) * time.Minute),
				})
			}

			return payload, err
		}
	}

	return &Auth{
		Login: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthLogin()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return setCookie(w, r.URL)(ah.Login(r.Context(), params))
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return setCookie(w, r.URL)(ah.Create(r.Context(), params))
			})
		},
	}
}
