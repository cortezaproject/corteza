package types

import (
	"encoding/json"
)

func ParseApigwfFilterParams(ss []string) (p ApigwFilterParams, err error) {
	p = ApigwFilterParams{}
	return p, parseStringsInput(ss, p)
}

func parseStringsInput(ss []string, p interface{}) (err error) {
	if len(ss) == 0 {
		return
	}

	return json.Unmarshal([]byte(ss[0]), &p)
}
