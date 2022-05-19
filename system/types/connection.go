package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
)

type (
	Connection struct {
		ID     uint64 `json:"id,string"`
		Handle string `json:"handle"`

		DSN       string `json:"dsn"`
		Location  string `json:"location"`
		Ownership string `json:"ownership"`
		Sensitive bool   `json:"sensitive"`

		Config       ConnectionConfig       `json:"config"`
		Capabilities ConnectionCapabilities `json:"capabilities"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ConnectionCapabilities struct {
		Enforced    capabilities.Set `json:"enforced"`
		Supported   capabilities.Set `json:"supported"`
		Unsupported capabilities.Set `json:"unsupported"`
		Enabled     capabilities.Set `json:"enabled"`
	}

	ConnectionConfig struct {
		DefaultModelIdent     string `json:"defaultModelIdent"`
		DefaultAttributeIdent string `json:"defaultAttributeIdent"`

		DefaultPartitionFormat string `json:"defaultPartitionFormat"`
	}
)
