package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type AuthClient struct {
	ID          uint64                           `json:"authClientID,string"`
	TenantID    uint64                           `json:"tenantID,string,omitempty"`
	ProjectID   uint64                           `json:"projectID,string,omitempty"`
	Handle      string                           `json:"handle"`
	Meta        *AuthClientMeta                  `json:"meta,omitempty"`
	Secret      string                           `json:"secret,omitempty"`
	Scope       string                           `json:"scope"`
	ValidGrant  string                           `json:"validGrant"`
	RedirectURI string                           `json:"redirectURI"`
	Enabled     bool                             `json:"enabled"`
	Trusted     bool                             `json:"trusted"`
	ValidFrom   *time.Time                       `json:"validFrom,omitempty"`
	ExpiresAt   *time.Time                       `json:"expiresAt,omitempty"`
	Security    *AuthClientSecurity              `json:"security"`
	OwnedBy     uint64                           `json:"ownedBy"`
	CreatedAt   time.Time                        `json:"createdAt"`
	UpdatedAt   *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy   uint64                           `json:"createdBy"`
	UpdatedBy   uint64                           `json:"updatedBy,omitempty"`
	DeletedBy   uint64                           `json:"deletedBy,omitempty"`
	Labels      map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r AuthClient) Clone() *AuthClient {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *AuthClientMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AuthClientMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAuthClientMeta(ss []string) (p *AuthClientMeta, err error) {
	p = &AuthClientMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *AuthClientSecurity) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AuthClientSecurity) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAuthClientSecurity(ss []string) (p *AuthClientSecurity, err error) {
	p = &AuthClientSecurity{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}
