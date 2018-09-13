package rbac

import (
	"encoding/json"
	"fmt"
	"github.com/crusttech/crust/internal/rbac/types"
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

		AddRole(username string, roles ...string) error
		RemoveRole(username string, roles ...string) error
	}
)

const (
	usersCreate = "/users/%s"
	usersGet    = "/users/%s"
	usersDelete = "/users/%s"
	// @todo: plural for users, but singular for sessions
	usersAddRole    = "/users/%s/assignRoles"
	usersRemoveRole = "/users/%s/deassignRoles"
)

func (u *Users) Create(username, password string) error {
	body := struct {
		Password string `json:"password"`
	}{password}

	resp, err := u.Client.Post(fmt.Sprintf(usersCreate, username), body)
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

func (u *Users) AddRole(username string, roles ...string) error {
	body := struct {
		Roles []string `json:"roles"`
	}{roles}

	resp, err := u.Client.Patch(fmt.Sprintf(usersAddRole, username), body)
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

func (u *Users) RemoveRole(username string, roles ...string) error {
	body := struct {
		Roles []string `json:"roles"`
	}{roles}

	resp, err := u.Client.Patch(fmt.Sprintf(usersRemoveRole, username), body)
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
	resp, err := u.Client.Get(fmt.Sprintf(usersGet, username))
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
	resp, err := u.Client.Delete(fmt.Sprintf(usersDelete, username))
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
