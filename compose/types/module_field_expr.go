package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type (
	ModuleFieldExpr struct {
		ValueExpr string `json:"value,omitempty"`

		Sanitizers []string `json:"sanitizers,omitempty"`

		Validators               []ModuleFieldValidator `json:"validators,omitempty"`
		DisableDefaultValidators bool                   `json:"disableDefaultValidators,omitempty"`

		Formatters               []string `json:"formatters,omitempty"`
		DisableDefaultFormatters bool     `json:"disableDefaultFormatters,omitempty"`
	}

	ModuleFieldValidator struct {
		Test  string `json:"test,omitempty"`
		Error string `json:"error,omitempty"`
	}
)

func (opt *ModuleFieldExpr) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*opt = ModuleFieldExpr{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, opt); err != nil {
			return fmt.Errorf("cannot scan '%v' into ModuleFieldExpr: %v", string(b), err)
		}
	}

	return nil
}

func (opt ModuleFieldExpr) Value() (driver.Value, error) {
	return json.Marshal(opt)
}
