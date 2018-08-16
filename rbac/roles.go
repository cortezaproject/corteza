package rbac

import (
	"encoding/json"
	"fmt"
	"github.com/crusttech/crust/rbac/types"
	"github.com/pkg/errors"
	"strings"
)

type (
	Roles struct {
		*Client
	}

	RolesInterface interface {
		Create(rolepath string) error
		CreateNested(rolepaths ...string) error
		Get(rolepath string) (*types.Role, error)
		Delete(rolepath string) error
	}
)

const (
	rolesCreate = "/roles/%s"
	rolesGet    = "/roles/%s"
	rolesDelete = "/roles/%s"
)

func (u *Roles) Create(rolepath string) error {
	if rolepath == "" {
		return errors.New("tried creating empty role")
	}
	resp, err := u.Client.Post(fmt.Sprintf(rolesCreate, rolepath), nil)
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

func (u *Roles) CreateNested(rolepaths ...string) error {
	if len(rolepaths) == 0 {
		return errors.New("tried creating empty role")
	}
	return u.Create(strings.Join(rolepaths, "/"))
}

func (u *Roles) Get(rolepath string) (*types.Role, error) {
	resp, err := u.Client.Get(fmt.Sprintf(rolesDelete, rolepath))
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		role := &types.Role{}
		return role, errors.Wrap(json.NewDecoder(resp.Body).Decode(role), "decoding json failed")
	default:
		return nil, toError(resp)
	}
}

func (u *Roles) GetNested(rolepaths ...string) (*types.Role, error) {
	if len(rolepaths) == 0 {
		return nil, errors.New("tried creating empty role")
	}
	return u.Get(strings.Join(rolepaths, "/"))
}

func (u *Roles) Delete(rolepath string) error {
	resp, err := u.Client.Delete(fmt.Sprintf(rolesDelete, rolepath))
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

var _ RolesInterface = &Roles{}
