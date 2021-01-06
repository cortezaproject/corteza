package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type (
	// Used for expression steps, arguments/results mapping and for input validation
	WorkflowExpression struct {
		Name  string `json:"name,omitempty"`
		Expr  string `json:"expr"`
		Error string `json:"error,omitempty"`
	}

	// Variables uses same structure as wfexec.Variables
	// but implements scan/value to support serialization into store
	Variables map[string]interface{}
)

func ParseWorkflowVariables(ss []string) (p Variables, err error) {
	p = Variables{}
	return p, parseStringsInput(ss, &p)
}

func (vv *Variables) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = Variables{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("can not scan '%v' into Variables: %w", string(b), err)
		}
	}

	return nil
}

func (vv Variables) Value() (driver.Value, error) {
	return json.Marshal(vv)
}
