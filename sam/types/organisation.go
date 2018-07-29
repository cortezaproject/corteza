package types

import (
	"time"
)

type (
	// Organisations - Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.
	Organisation struct {
		ID         uint64     `db:"id"`
		FQN        string     `db:"fqn"`
		Name       string     `db:"name"`
		CreatedAt  time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	OrganisationFilter struct {
		Query string
	}
)
