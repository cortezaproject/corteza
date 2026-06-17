package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type DataPrivacyRequest struct {
	ID          uint64                       `json:"requestID,string"`
	TenantID    uint64                       `json:"tenantID,string,omitempty"`
	ProjectID   uint64                       `json:"projectID,string,omitempty"`
	Kind        RequestKind                  `json:"kind"`
	Status      RequestStatus                `json:"status"`
	Payload     DataPrivacyRequestPayloadSet `json:"payload,omitempty"`
	RequestedAt time.Time                    `json:"requestedAt,omitempty"`
	RequestedBy uint64                       `json:"requestedBy,string"`
	CompletedAt *time.Time                   `json:"completedAt,omitempty"`
	CompletedBy uint64                       `json:"completedBy,string,omitempty"`
	CreatedAt   time.Time                    `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time                   `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time                   `json:"deletedAt,omitempty"`
	CreatedBy   uint64                       `json:"createdBy,string"`
	UpdatedBy   uint64                       `json:"updatedBy,string,omitempty"`
	DeletedBy   uint64                       `json:"deletedBy,string,omitempty"`
}

func (m *DataPrivacyRequestPayloadSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DataPrivacyRequestPayloadSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDataPrivacyRequestPayloadSet(ss []string) (p DataPrivacyRequestPayloadSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
