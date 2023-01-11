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

func parseMapStringString(ss []string) (map[string]string, error) {
	if len(ss) == 0 {
		return nil, nil
	}

	out := make(map[string]string, 8)

	return out, json.Unmarshal([]byte(ss[0]), out)
}
