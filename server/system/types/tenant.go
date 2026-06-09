package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
)

type (
	// Tenant is the top-level isolation unit. Every user belongs to exactly one
	// tenant; projects and tenant-scoped resources nest beneath it.
	Tenant struct {
		ID     uint64       `json:"tenantID,string"`
		Handle string       `json:"handle"`
		Status TenantStatus `json:"status"`

		Config TenantConfig `json:"config"`
		Meta   TenantMeta   `json:"meta"`

		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		CreatedAt   time.Time  `json:"createdAt,omitempty"`
		CreatedBy   uint64     `json:"createdBy,string"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy   uint64     `json:"updatedBy,string,omitempty"`
		SuspendedAt *time.Time `json:"suspendedAt,omitempty"`
		DeletedAt   *time.Time `json:"deletedAt,omitempty"`
		DeletedBy   uint64     `json:"deletedBy,string,omitempty"`
	}

	TenantFilter struct {
		TenantID []string     `json:"tenantID"`
		Handle   string       `json:"handle"`
		Status   TenantStatus `json:"status"`
		Query    string       `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*Tenant) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	// TenantConfig holds behavioral settings. System acts on these.
	TenantConfig struct {
		// DALConnectionID points at the tenant's data store.
		DALConnectionID uint64 `json:"dalConnectionID,string,omitempty"`
		// AuthProviders are the allowed auth provider handles.
		AuthProviders []string `json:"authProviders,omitempty"`
		// FeatureFlags are tenant-specific feature overrides.
		FeatureFlags map[string]bool `json:"featureFlags,omitempty"`
		// Quotas cap tenant resource usage.
		Quotas TenantQuotas `json:"quotas,omitempty"`
		// Locale is the default tenant locale (BCP 47).
		Locale string `json:"locale,omitempty"`
		// Timezone is the default tenant timezone (IANA).
		Timezone string `json:"timezone,omitempty"`
	}

	// TenantQuotas are per-tenant resource caps. 0 = unlimited.
	TenantQuotas struct {
		MaxUsers    int   `json:"maxUsers,omitempty"`
		MaxProjects int   `json:"maxProjects,omitempty"`
		MaxStorage  int64 `json:"maxStorage,omitempty"`
	}

	// TenantMeta is display-only. System does not act on these.
	TenantMeta struct {
		Short       string   `json:"short,omitempty"`
		Description string   `json:"description,omitempty"`
		LogoID      uint64   `json:"logoID,string,omitempty"`
		Color       string   `json:"color,omitempty"`
		Tags        []string `json:"tags,omitempty"`
	}

	TenantStatus string
)

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusArchived  TenantStatus = "archived"
)

func (m *TenantConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m TenantConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *TenantMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m TenantMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseTenantConfig(ss []string) (p TenantConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseTenantMeta(ss []string) (p TenantMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
