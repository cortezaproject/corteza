package rules

import (
	"encoding/json"
	"errors"
)

type Access int

const (
	Allow   Access = 2
	Deny           = 1
	Inherit        = 0
)

type Rules struct {
	TeamID    uint64 `db:"rel_team"`
	Resource  string `db:"resource"`
	Operation string `db:"operation"`
	Value     Access `db:"value"`
}

func (a *Access) UnmarshalJSON(data []byte) error {
	var i interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	s, ok := i.(string)
	if !ok {
		return errors.New("Type assertion .(string) failed.")
	}

	switch s {
	case "allow":
		*a = Allow
	case "deny":
		*a = Deny
	default:
		*a = Inherit
	}
	return nil
}
