package payload

import (
	"encoding/json"

	"github.com/crusttech/crust/internal/payload/incoming"
)

func Unmarshal(raw []byte) (*incoming.Payload, error) {
	var p = &incoming.Payload{}
	return p, json.Unmarshal(raw, p)
}
