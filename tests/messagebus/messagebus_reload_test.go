package messagebus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

var (
	testQueueNewlyCreated = &messagebus.QueueSettings{
		ID:      2,
		Queue:   "new_queue_created_after_init",
		Handler: string(messagebus.HandlerSql),
		Meta: messagebus.QueueSettingsMeta{
			PollDelay:      makeDelay(time.Second),
			DispatchEvents: true,
		},
	}
)

func TestMessageReload(t *testing.T) {
	var (
		h   = newHelper(t)
		ctx = context.Background()
	)

	h.prepareRBAC()
	h.prepareQueues(ctx, testQueueDispatched)

	// reinit the messagebus
	h.initMessagebus(ctx)

	messagebus.Service().Watch(ctx, service.DefaultStore)

	h.noError(testApp.Store.CreateMessagebusQueuesetting(ctx, testQueueNewlyCreated))

	w := sync.WaitGroup{}
	w.Add(1)

	// send the signal to reload all queues
	messagebus.Service().ReloadQueues()

	go func() {
		sent := false
		for {
			select {
			default:
				if !sent {
					// first, check if the new queue has been already registered
					queue := messagebus.Service().Write("new_queue_created_after_init")

					if queue == nil {
						break
					}

					queue <- []byte(`this is a test`)
					queue <- []byte(`foo bar`)

					// set sent to true, so the above would not call
					// again and block again
					sent = true
				}

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
