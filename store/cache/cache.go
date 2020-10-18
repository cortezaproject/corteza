package cache

import (
	"fmt"
)

func iKey(aa ...interface{}) (key string) {
	for _, a := range aa {
		key += fmt.Sprintf(":%v", a)
	}

	return
}
