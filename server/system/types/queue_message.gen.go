package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/pkg/revisions"
	"reflect"
	"time"
)

type (
	QueueMessage struct {
		ID        uint64     `json:"messageID"`
		Queue     string     `json:"queue"`
		Payload   []byte     `json:"payload"`
		Created   *time.Time `json:"created"`
		Processed *time.Time `json:"processed"`
	}
)

func (r QueueMessage) Clone() *QueueMessage {
	dup := r
	if r.Payload != nil {
		dup.Payload = make([]byte, len(r.Payload))
		copy(dup.Payload, r.Payload)
	}

	if r.Created != nil {
		v := *r.Created
		dup.Created = &v
	}

	if r.Processed != nil {
		v := *r.Processed
		dup.Processed = &v
	}

	return &dup
}

func (r QueueMessage) Diff(cmp *QueueMessage) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &QueueMessage{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "messageID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Queue != cmp.Queue {
		out = append(out, &revisions.Change{Key: "queue", Old: []any{cmp.Queue}, New: []any{r.Queue}})
	}

	if !reflect.DeepEqual(r.Payload, cmp.Payload) {
		out = append(out, &revisions.Change{Key: "payload", Old: []any{cmp.Payload}, New: []any{r.Payload}})
	}

	if !reflect.DeepEqual(r.Created, cmp.Created) {
		out = append(out, &revisions.Change{Key: "created", Old: []any{cmp.Created}, New: []any{r.Created}})
	}

	if !reflect.DeepEqual(r.Processed, cmp.Processed) {
		out = append(out, &revisions.Change{Key: "processed", Old: []any{cmp.Processed}, New: []any{r.Processed}})
	}

	return out
}
