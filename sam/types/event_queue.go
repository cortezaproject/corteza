package types

import (
	"encoding/json"
)

type (
	EventQueueItem struct {
		ID         uint64          `db:"id"`
		Origin     uint64          `db:"origin"`
		Subscriber string          `db:"subscriber"`
		Payload    json.RawMessage `db:"payload"`
	}
)
