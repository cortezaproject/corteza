package types

import (
	"time"
)

type (
	// Teams - An organisation may have many teams. Teams may have many channels available. Access to channels may be shared between teams.
	Team struct {
		ID         uint64     `db:"id"`
		Name       string     `db:"name"`
		Handle     string     `db:"handle"`
		CreatedAt  time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	TeamFilter struct {
		Query string
	}
)
