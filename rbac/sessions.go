package rbac

import (
	"encoding/json"
	"fmt"
	"github.com/crusttech/crust/rbac/types"
	"github.com/pkg/errors"
)

type (
	Sessions struct {
		*Client
	}

	SessionsInterface interface {
		Create(sessionID, username string, roles ...string) error
		Get(sessionID string) (*types.Session, error)
		Delete(sessionID string) error

		ActivateRole(sessionID string, roles ...string) error
		DeactivateRole(sessionID string, roles ...string) error
	}
)

const (
	sessionsCreate         = "/sessions/%s"
	sessionsGet            = "/sessions/%s"
	sessionsDelete         = "/sessions/%s"
	sessionsActivateRole   = "/sessions/%s/activateRole"
	sessionsDeactivateRole = "/sessions/%s/deactivateRole"
)

func (u *Sessions) Create(sessionID, username string, roles ...string) error {
	body := struct {
		Username string   `json:"username"`
		Roles    []string `json:"roles,omitempty"`
	}{username, roles}

	resp, err := u.Client.Post(fmt.Sprintf(sessionsCreate, sessionID), body)
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

func (u *Sessions) Get(sessionID string) (*types.Session, error) {
	resp, err := u.Client.Get(fmt.Sprintf(sessionsGet, sessionID))
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		session := &types.Session{}
		return session, errors.Wrap(json.NewDecoder(resp.Body).Decode(session), "decoding json failed")
	default:
		return nil, toError(resp)
	}
}

func (u *Sessions) Delete(sessionID string) error {
	resp, err := u.Client.Delete(fmt.Sprintf(sessionsDelete, sessionID))
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

func (u *Sessions) ActivateRole(sessionID string, roles ...string) error {
	body := struct {
		Roles []string `json:"roles"`
	}{roles}

	resp, err := u.Client.Patch(fmt.Sprintf(sessionsActivateRole, sessionID), body)
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

func (u *Sessions) DeactivateRole(sessionID string, roles ...string) error {
	body := struct {
		Roles []string `json:"roles"`
	}{roles}

	resp, err := u.Client.Patch(fmt.Sprintf(sessionsDeactivateRole, sessionID), body)
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

var _ SessionsInterface = &Sessions{}
