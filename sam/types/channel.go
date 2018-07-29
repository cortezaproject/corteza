package types

import (
	"encoding/json"
	"time"
)

type (
	Channel struct {
		ID            uint64          `db:"id"`
		Name          string          `db:"name"`
		Topic         string          `db:"-"`
		Meta          json.RawMessage `db:"meta"`
		LastMessageID uint64          `json:",omitempty" db:"rel_last_message"`
		CreatedAt     time.Time       `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt     *time.Time      `json:"updated_at,omitempty" db:"updated_at"`
		ArchivedAt    *time.Time      `json:"archived_at,omitempty" db:"archived_at"`
		DeletedAt     *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	ChannelFilter struct {
		Query string
	}
)
