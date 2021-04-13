package messagebus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

type (
	dispatcher interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Dispatch(ctx context.Context, ev eventbus.Event)
	}
)

var (
	testQueueDispatched = &messagebus.QueueSettings{
		ID:      1,
		Queue:   "test",
		Handler: string(messagebus.HandlerSql),
		Meta: messagebus.QueueSettingsMeta{
			PollDelay:      makeDelay(time.Second),
			DispatchEvents: true,
		},
	}

	testQueueMessage = messagebus.QueueMessage{
		ID:      1,
		Queue:   "test",
		Payload: []byte(`{"foo": "bar"}`),
		Created: now(),
	}
)

func TestMessageWrite(t *testing.T) {
	var (
		h   = newHelper(t)
		ctx = context.Background()
	)

	h.prepareRBAC()
	h.prepareQueues(ctx, testQueueDispatched)

	// reinit the messagebus
	h.initMessagebus(ctx)

	queue := messagebus.Service().Write("test")
	timeout := time.After(time.Second * 5)

	w := sync.WaitGroup{}
	w.Add(1)

	go func() {
		queue <- []byte(`this is a test`)
		queue <- []byte(`foo bar`)

		for {
			select {
			case <-timeout:
				w.Done()
				t.Fail()
				return

			default:
				set := h.checkPersistedMessages(ctx, messagebus.QueueMessageFilter{})

				// success, will eventually get persisted
				if len(set) == 2 {
					h.a.NotEmpty(set)
					h.a.Equal([]byte(`this is a test`), set[0].Payload)
					h.a.Equal([]byte(`foo bar`), set[1].Payload)
					w.Done()
					return
				}
			}
		}
	}()

	w.Wait()
}
