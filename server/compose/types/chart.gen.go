package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"time"
)

type Chart struct {
	ID          uint64                           `json:"chartID,string"`
	Handle      string                           `json:"handle"`
	TenantID    uint64                           `json:"tenantID,string,omitempty"`
	ProjectID   uint64                           `json:"projectID,string,omitempty"`
	NamespaceID uint64                           `json:"namespaceID,string"`
	Name        string                           `json:"name"`
	Config      ChartConfig                      `json:"config"`
	CreatedAt   time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time                       `json:"deletedAt,omitempty"`
	Labels      map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r Chart) Clone() *Chart {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}
