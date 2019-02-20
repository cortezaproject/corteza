package rules

import (
	"fmt"

	"encoding/json"
)

type Resource struct {
	ID      uint64 `json:"id,string"`
	Name    string `json:"name"`
	Scope   string `json:"scope"`
	Service string `json:"service"`
}

type ResourceJSON struct {
	ID         uint64 `json:"id,string"`
	Name       string `json:"name"`
	Scope      string `json:"scope"`
	Service    string `json:"service"`
	ResourceID string `json:"resource"`
}

func (r Resource) String() string {
	if r.ID > 0 {
		return fmt.Sprintf("%s:%s:%d", r.Service, r.Scope, r.ID)
	}
	return ""
}

func (r Resource) All() string {
	return fmt.Sprintf("%s:%s:*", r.Service, r.Scope)
}

func (r Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(ResourceJSON{
		r.ID,
		r.Name,
		r.Scope,
		r.Service,
		r.String(),
	})
}

var _ fmt.Stringer = Resource{}
