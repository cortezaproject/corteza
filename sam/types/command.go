package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type (
	Command struct {
		Name        string          `db:"name"        json:"name"`
		Params      CommandParamSet `db:"params"      json:"params"`
		Description string          `db:"description" json:"description"`
	}

	CommandParam struct {
		Name     string `db:"name"     json:"name"`
		Type     string `db:"type"     json:"type"`
		Required bool   `db:"required" json:"required"`
	}

	CommandSet      []*Command
	CommandParamSet []*CommandParam
)

var (
	Preset CommandSet // @todo move this to someplace safe
)

func init() {
	Preset = CommandSet{
		&Command{
			Name:        "echo",
			Description: "It does exactly what it says on the tin"},
		&Command{
			Name:        "shrug",
			Description: "It does exactly what it says on the tin"},
	}
}

func (cc CommandSet) Walk(w func(*Command) error) (err error) {
	for i := range cc {
		if err = w(cc[i]); err != nil {
			return
		}
	}

	return
}

func (pp CommandParamSet) Walk(w func(*CommandParam) error) (err error) {
	for i := range pp {
		if err = w(pp[i]); err != nil {
			return
		}
	}

	return
}

func (pp *CommandParamSet) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		*pp = CommandParamSet{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), pp); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into CommandParamSet", value)
		}
	}

	return nil
}

func (pp CommandParamSet) Value() (driver.Value, error) {
	return json.Marshal(pp)
}
