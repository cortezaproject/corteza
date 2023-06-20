package ddl

import "strings"

func ParseColumnTypes(c *Column) (original, name string, meta []string) {
	original = strings.ToLower(c.Type.Name)

	pp := strings.Split(original, "(")
	name = pp[0]
	if len(pp) > 1 {
		meta = strings.Split(strings.TrimRight(pp[1], ")"), ",")
		for i, m := range meta {
			meta[i] = strings.TrimSpace(m)
		}
	}

	return
}
