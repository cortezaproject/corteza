package messagebus

import (
	"context"
	"time"
)

type (
	CortezaChan struct {
		queue   string
		payload []byte
	}

	CortezaMessageStore struct {
		store      chan CortezaChan
		queueChans map[string]chan string
		// messages   chan string
	}

	CortezaQueueHandler struct {
		handle   HandlerType
		messages CortezaMessageStore
		poll     *time.Ticker
	}
)

func (cms CortezaMessageStore) Write(payload []byte) {
	// if _, ok := cms.queueChans[queue]; !ok {
	// 	cms.queueChans[queue] = make(chan string)
	// }

	// cms.queueChans[queue] <- payload.(string)
	cms.store <- CortezaChan{
		payload: payload,
	}
}

func (cms CortezaMessageStore) Read() []byte {
	// fan out based on queue?
	p := <-cms.store
	// p := <-cms.queueChans[queue]
	return p.payload
}

func NewCortezaHandler(settings QueueSettings) *CortezaQueueHandler {
	h := &CortezaQueueHandler{
		handle: HandlerCorteza,
		messages: CortezaMessageStore{
			store:      make(chan CortezaChan),
			queueChans: make(map[string]chan string),
		},
	}

	if settings.Meta.PollDelay != nil {
		h.poll = time.NewTicker(*settings.Meta.PollDelay)
	}

	return h
}

func (cq *CortezaQueueHandler) Write(ctx context.Context, p []byte) (err error) {
	// return cq.messages.Write(p)
	return
}

func (cq *CortezaQueueHandler) Read(ctx context.Context) (set QueueMessageSet, err error) {
	// return cq.messages.Read()
	return
}

// getTicker fetches the ticker channel if it is set-up in
// queue settings (see QueueSettingsMeta.PollDelay)
func (cq *CortezaQueueHandler) Ticker(ctx context.Context) <-chan time.Time {
	if cq.poll != nil {
		return cq.poll.C
	}

	return nil
}

// getNotification fetches the notification channel from the store
// notification mechanism
func (cq *CortezaQueueHandler) Notification(ctx context.Context) <-chan interface{} {
	// @todo - psql, redis (mysql only with a plugin)
	return nil
}

// Process marks a message as processed in store of choice
func (cq *CortezaQueueHandler) Process(ctx context.Context, m QueueMessage) (err error) {
	// return cq.messages.process(ctx, m)
	return
}

func (cq *CortezaQueueHandler) SetStorer(qs QueueStorer) {}
