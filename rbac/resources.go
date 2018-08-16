package rbac

import (
	"encoding/json"
	"fmt"
	"github.com/crusttech/crust/rbac/types"
	"github.com/pkg/errors"
)

type (
	Resources struct {
		*Client
	}

	ResourcesRole struct {
		Rolepath  string `json:"role"`
		Operation string `json:"operation"`
	}

	ResourcesInterface interface {
		Create(resourceID string, operations []string) error
		Get(resourceID string) (*types.Resource, error)
		Delete(resourceID string) error

		Grant(resourceID string, rolepath string, operations []string) error
		GrantMultiple(resourceID string, roles []ResourcesRole) error

		CheckAccess(resourceID string, operation string, sessionID string) error
	}
)

const (
	resourcesCreate      = "/resources/%s"
	resourcesGet         = "/resources/%s"
	resourcesDelete      = "/resources/%s"
	resourcesGrant       = "/resources/%s/grantPermission"
	resourcesCheckAccess = "/resources/%s/checkAccess?operation=%s&session=%s"
)

func (u *Resources) CheckAccess(resourceID string, operation string, sessionID string) error {
	resp, err := u.Client.Get(fmt.Sprintf(resourcesCheckAccess, resourceID, operation, sessionID))
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

func (u *Resources) Grant(resourceID string, rolepath string, operations []string) error {
	body := make([]ResourcesRole, len(operations))
	for index, operation := range operations {
		body[index] = ResourcesRole{rolepath, operation}
	}
	return u.GrantMultiple(resourceID, body)
}

func (u *Resources) GrantMultiple(resourceID string, roles []ResourcesRole) error {
	body := struct {
		Permissions []ResourcesRole `json:"permissions"`
	}{roles}

	resp, err := u.Client.Patch(fmt.Sprintf(resourcesGrant, resourceID), body)
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

func (u *Resources) Get(resourceID string) (*types.Resource, error) {
	resp, err := u.Client.Get(fmt.Sprintf(resourcesGet, resourceID))
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		resource := &types.Resource{}
		return resource, errors.Wrap(json.NewDecoder(resp.Body).Decode(resource), "decoding json failed")
	default:
		return nil, toError(resp)
	}
}

func (u *Resources) Delete(resourceID string) error {
	resp, err := u.Client.Delete(fmt.Sprintf(resourcesDelete, resourceID))
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
