package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"time"
)

type ProjectGroupEntry struct {
	ID             uint64    `json:"-"`
	ProjectGroupID uint64    `json:"projectGroupID,string"`
	ResourceRef    string    `json:"resourceRef"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}
