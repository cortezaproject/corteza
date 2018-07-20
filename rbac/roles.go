package rbac

import (
	"fmt"
	"github.com/pkg/errors"
)

type (
	Roles struct {
		*Client
	}

	RolesInterface interface {
		Create(role string) error
		Delete(role string) error
	}
)

const (
	rolesCreate = "/roles/%s"
	rolesDelete = "/roles/%s"
)

func (u *Roles) Create(role string) error {
	resp, err := u.Client.Post(fmt.Sprintf(rolesCreate, role), nil)
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

func (u *Roles) Delete(role string) error {
	resp, err := u.Client.Delete(fmt.Sprintf(rolesDelete, role))
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
