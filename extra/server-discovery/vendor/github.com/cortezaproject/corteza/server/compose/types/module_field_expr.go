package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza/server/pkg/sql"
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
		ValidatorID uint64 `json:"validatorID,string,omitempty"`
		Test        string `json:"test,omitempty"`
		Error       string `json:"error,omitempty"`
	}
)

func (opt *ModuleFieldExpr) Scan(src any) error          { return sql.ParseJSON(src, opt) }
func (opt ModuleFieldExpr) Value() (driver.Value, error) { return json.Marshal(opt) }
