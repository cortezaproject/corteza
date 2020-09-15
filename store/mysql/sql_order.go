package mysql

import "fmt"

func SqlSortHandler(exp string, desc bool) string {
	if desc {
		return fmt.Sprintf("%s DESC", exp)
	} else {
		return fmt.Sprintf("%s ASC", exp)
	}
}
