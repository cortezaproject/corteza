package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"time"
)

type QueueMessage struct {
	ID        uint64     `json:"messageID"`
	Queue     string     `json:"queue"`
	Payload   []byte     `json:"payload"`
	Created   *time.Time `json:"created"`
	Processed *time.Time `json:"processed"`
}
