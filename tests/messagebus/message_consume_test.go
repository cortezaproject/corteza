package messagebus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	dispatcher interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Dispatch(ctx context.Context, ev eventbus.Event)
	}
)

var (
	testQueueDispatched = &messagebus.QueueSettings{
		ID:       1,
		Queue:    "test",
		Consumer: string(messagebus.ConsumerStore),
		Meta: messagebus.QueueSettingsMeta{
			PollDelay:      makeDelay(time.Second),
			DispatchEvents: true,
		},
	}

	testQueueEb = &messagebus.QueueSettings{
		ID:       1,
		Queue:    "test_eb",
		Consumer: string(messagebus.ConsumerEventbus),
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

	h.prepareQueues(ctx, testQueueDispatched)

	// reinit the messagebus
	h.initMessagebus(ctx)

	timeout := time.After(time.Second * 5)

	w := sync.WaitGroup{}
	w.Add(1)

	go func() {
		messagebus.Service().Push("test", []byte("this is a test"))
		messagebus.Service().Push("test", []byte("foo bar"))

		for {
			select {
			case <-timeout:
				w.Done()
				t.Fail()
				return

			default:
				set := h.checkPersistedMessages(ctx, messagebus.QueueMessageFilter{})

				// success, will eventually get persisted
				if len(set) >= 2 {
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

func TestMessageWriteEventbus(t *testing.T) {
	var (
		h   = newHelper(t)
		ctx = context.Background()

		messageProcessed   = make(chan bool)
		messageProcessedFn = func(c chan bool) eventbus.HandlerFn {
			return func(ctx context.Context, ev eventbus.Event) error {
				c <- true
				return nil
			}
		}
	)

	h.prepareQueues(ctx, testQueueEb)

	// reinit the messagebus
	h.initMessagebus(ctx)

	// safety net, the test should be instantaneous
	timeout := time.After(time.Second * 3)

	w := sync.WaitGroup{}
	w.Add(1)

	// eventbus is the consumer, update temporary channel
	// so we get the value back lower
	registerTestEvent(messageProcessedFn(messageProcessed))

	go func() {
		messagebus.Service().Push(testQueueEb.Queue, []byte("this is a test"))

		for {
			select {
			case <-timeout:
				w.Done()
				t.Fail()
				return

			case <-messageProcessed:
				w.Done()
				return
			}
		}
	}()

	w.Wait()
}

func registerTestEvent(fn eventbus.HandlerFn) {
	e := event.QueueOnMessage(&types.QueueMessage{})

	eventbus.
		Service().
		Register(fn,
			eventbus.On(e.EventType()),
			eventbus.For(e.ResourceType()))
}
