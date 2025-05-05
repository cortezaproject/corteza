package rbac

import "strings"

// permuteResource returns the given identifier lvl and all lower levels
func permuteResource(res string) (out []string) {
	out = append(out, res)
	rr := strings.Split(res, "/")
	for i := len(rr) - 1; i >= 1; i-- {
		if rr[i] == "*" {
			continue
		}

		rr[i] = "*"
		out = append(out, strings.Join(rr, "/"))
	}

	return
}
