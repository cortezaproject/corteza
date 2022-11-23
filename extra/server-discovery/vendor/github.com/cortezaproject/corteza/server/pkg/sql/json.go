package sql

import (
	"encoding/json"
	"fmt"
)

func ParseJSON(raw any, dest any) error {
	var (
		data []byte
	)

	if b, ok := raw.([]byte); ok {
		data = b
	} else if s, ok := raw.(string); ok {
		data = []byte(s)
	} else if raw == nil {
		return nil
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("can not scan JSON into %T: %w", dest, err)
	}

	return nil
}
