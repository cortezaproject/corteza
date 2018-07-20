package rbac

import (
	"encoding/json"
	"github.com/crusttech/crust/rbac/types"
	"github.com/pkg/errors"
)

type (
	Users struct {
		*Client
	}

	UsersInterface interface {
		Create(username, password string) error
		Get(username string) (*types.User, error)
		Delete(username string) error
	}
)

func (u *Users) Create(username, password string) error {
	body := struct {
		Password string `json:"password"`
	}{password}

	resp, err := u.Client.Post("/users/"+username, body)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil
	default:
		return toError(resp)
	}
}

func (u *Users) Get(username string) (*types.User, error) {
	resp, err := u.Client.Get("/users/" + username)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		user := &types.User{}
		return user, errors.Wrap(json.NewDecoder(resp.Body).Decode(user), "decoding json failed")
	default:
		return nil, toError(resp)
	}
}

func (u *Users) Delete(username string) error {
	resp, err := u.Client.Delete("/users/" + username)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil
	default:
		return toError(resp)
	}
}

var _ UsersInterface = &Users{}
