package types

import "encoding/json"

func parseStringsInput(ss []string, p interface{}) (err error) {
	if len(ss) == 0 {
		return
	}

	return json.Unmarshal([]byte(ss[0]), &p)
}
