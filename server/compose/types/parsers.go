package types

import "encoding/json"

func ParsePageLayoutMeta(ss []string) (p PageLayoutMeta, err error) {
	p = PageLayoutMeta{}
	return p, parseStringsInput(ss, &p)
}

func ParsePageMeta(ss []string) (p PageMeta, err error) {
	p = PageMeta{}
	return p, parseStringsInput(ss, &p)
}

func parseStringsInput(ss []string, p interface{}) (err error) {
	if len(ss) == 0 {
		return
	}

	return json.Unmarshal([]byte(ss[0]), &p)
}
