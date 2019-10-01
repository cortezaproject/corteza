package payload

import (
	"encoding/json"

	"github.com/cortezaproject/corteza-server/pkg/payload/incoming"
)

func Unmarshal(raw []byte) (*incoming.Payload, error) {
	var p = &incoming.Payload{}
	return p, json.Unmarshal(raw, p)
}
