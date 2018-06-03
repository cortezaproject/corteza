package sam

import (
	"net/http"

	"github.com/pkg/errors"
)

// User login request parameters
type UserLoginRequest struct {
	username string
	password string
}

func (UserLoginRequest) new() *UserLoginRequest {
	return &UserLoginRequest{}
}

func (u *UserLoginRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: UserLoginRequest.Fill")
}

var _ RequestFiller = UserLoginRequest{}.new()

// User search request parameters
type UserSearchRequest struct {
	query string
}

func (UserSearchRequest) new() *UserSearchRequest {
	return &UserSearchRequest{}
}

func (u *UserSearchRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: UserSearchRequest.Fill")
}

var _ RequestFiller = UserSearchRequest{}.new()
