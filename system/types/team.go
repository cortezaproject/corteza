package types

import (
	"time"
)

type (
	// Teams - An organisation may have many teams. Teams may have many channels available. Access to channels may be shared between teams.
	Team struct {
		ID         uint64     `json:"id,string" db:"id"`
		Name       string     `json:"name" db:"name"`
		Handle     string     `json:"handle" db:"handle"`
		CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	TeamFilter struct {
		Query string
	}
)
