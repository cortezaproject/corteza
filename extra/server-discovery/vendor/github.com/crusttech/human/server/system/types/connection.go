package types

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ConnectionFilter struct {
		Handle string   `json:"handle"`
		Status []string `json:"status"`
		Query  string   `json:"query"`
		Tags   []string `json:"tags"`
		Source string   `json:"source"` // "catalog", "local", or "" for all

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*Connection) (bool, error) `json:"-"`

		filter.Paging
		filter.Sorting
	}

	// Slice types for DB JSON serialization
	ConnectionResources  []ConnectionResource
	ConnectionOperations []ConnectionOperation
)

func (ct *ConnectionTemplate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		ct.Value = s
		return nil
	}
	type Alias ConnectionTemplate
	return json.Unmarshal(data, (*Alias)(ct))
}

