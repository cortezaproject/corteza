package helpers

import "encoding/json"

func JSON(i interface{}) string {
	if b, err := json.Marshal(i); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}
