package messagebus

import (
	"context"
	"sync"
	"testing"

	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/system/service/event"
)

func TestMessageRead(t *testing.T) {
	var (
		h                  = newHelper(t)
		ctx                = context.Background()
		messageProcessed   = make(chan bool)
		messageProcessedFn = func(c chan bool) eventbus.HandlerFn {
			return func(ctx context.Context, ev eventbus.Event) error {
				c <- true
				return nil
			}
		}

		set messagebus.QueueMessageSet
		p   []byte
	)

	h.prepareRBAC()
	h.prepareQueues(ctx, testQueueDispatched)

	// reinit the messagebus
	h.initMessagebus(ctx)

	// register on message event
	// write to testing channel for us to use later
	// in the payload check
	registerTestEvent(messageProcessedFn(messageProcessed))

	// prepare store
	service.DefaultStore.CreateMessagebusQueuemessage(ctx, &testQueueMessage)

	queue := messagebus.Service().Read("test")

	w := sync.WaitGroup{}
	w.Add(2)

	go func() {
		for {
			select {
			case p = <-queue:
				w.Done()

			case <-messageProcessed:
				set = h.checkPersistedMessages(ctx, messagebus.QueueMessageFilter{Processed: filter.StateInclusive})
				w.Done()
			}
		}
	}()

	w.Wait()

	h.a.Len(set, 1)
	h.a.Equal(testQueueMessage.Payload, p)
}

func registerTestEvent(fn eventbus.HandlerFn) {
	e := event.QueueOnMessage(&expr.String{})

	eventbus.
		Service().
		Register(fn,
			eventbus.On(e.EventType()),
			eventbus.For(e.ResourceType()))
}
