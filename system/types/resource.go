package types

import (
	"fmt"

	"encoding/json"
)

type Resource struct {
	ID    uint64 `json:"id,string"`
	Name  string `json:"name"`
	Scope string `json:"scope"`
}

type ResourceJSON struct {
	ID         uint64 `json:"id,string"`
	Name       string `json:"name"`
	Scope      string `json:"scope"`
	ResourceID string `json:"resource"`
}

func (r Resource) String() string {
	return fmt.Sprintf("%s:%d", r.Scope, r.ID)
}

func (r Resource) All() string {
	return fmt.Sprintf("%s:*", r.Scope)
}

func (r Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(ResourceJSON{
		r.ID,
		r.Name,
		r.Scope,
		r.String(),
	})
}

var _ fmt.Stringer = Resource{}
