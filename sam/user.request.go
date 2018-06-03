package sam

import (
	"net/http"
)

// User login request parameters
type userLoginRequest struct {
	username string
	password string
}

func (userLoginRequest) new() *userLoginRequest {
	return &userLoginRequest{}
}

func (u *userLoginRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	u.username = post["username"]

	u.password = post["password"]
	return nil
}

var _ RequestFiller = userLoginRequest{}.new()

// User search request parameters
type userSearchRequest struct {
	query string
}

func (userSearchRequest) new() *userSearchRequest {
	return &userSearchRequest{}
}

func (u *userSearchRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	u.query = get["query"]
	return nil
}

var _ RequestFiller = userSearchRequest{}.new()
