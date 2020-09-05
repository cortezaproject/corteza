package types

import (
	"encoding/json"
)

type (
	EventQueueItem struct {
		ID         uint64
		Origin     uint64
		SubType    EventQueueItemSubType
		Subscriber string
		Payload    json.RawMessage
	}

	EventQueueItemSubType string
)

const (
	EventQueueItemSubTypeUser    = "user"
	EventQueueItemSubTypeChannel = "channel"
)
