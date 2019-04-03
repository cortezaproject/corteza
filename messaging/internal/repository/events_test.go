// +build unit

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/crusttech/crust/messaging/types"
)

func TestEvents(t *testing.T) {
	assert := func(ok bool, format string, args ...interface{}) {
		if !ok {
			t.Fatalf(format, args...)
		}
	}
	queue := Events()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	queue.Push(ctx, &types.EventQueueItem{Subscriber: "test1"})
	queue.Push(ctx, &types.EventQueueItem{Subscriber: "test2"})
	queue.Push(ctx, &types.EventQueueItem{Subscriber: "test3"})

	for i := 1; i <= 3; i++ {
		item, err := queue.Pull(ctx)
		assert(err == nil, "Expected non-error queue return, got %+v", err)
		assert(item != nil, "Expected non-empty queue item")
		expected := fmt.Sprintf("test%d", i)
		assert(item.Subscriber == expected, "Expected subscriber value doesn't match: %s != %s", expected, item.Subscriber)
	}
}
