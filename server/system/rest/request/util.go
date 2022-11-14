package request

import (
	"encoding/json"
)

func parseMapStringInterface(ss []string) (map[string]interface{}, error) {
	if len(ss) == 0 {
		return nil, nil
	}

	out := make(map[string]interface{})

	return out, json.Unmarshal([]byte(ss[0]), out)
}
