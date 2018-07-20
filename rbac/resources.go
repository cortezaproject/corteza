package rbac

import (
	"fmt"
	"github.com/pkg/errors"
)

type (
	Resources struct {
		*Client
	}

	ResourcesInterface interface {
		Create(resourceID string, operations []string) error
		Delete(resourceID string) error
	}
)

const (
	resourcesCreate = "/resources/%s"
	resourcesDelete = "/resources/%s"
)

func (u *Resources) Create(resourceID string, operations []string) error {
	body := struct {
		Operations []string `json:"operations"`
	}{operations}

	resp, err := u.Client.Post(fmt.Sprintf(resourcesCreate, resourceID), body)
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

func (u *Resources) Delete(resourceID string) error {
	resp, err := u.Client.Delete(fmt.Sprintf(resourcesCreate, resourceID))
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

var _ ResourcesInterface = &Resources{}
