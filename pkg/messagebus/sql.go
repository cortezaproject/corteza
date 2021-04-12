package messagebus

import (
	"context"
	"time"
)

type (
	SqlQueueHandler struct {
		queue  string
		handle HandlerType
		client SqlClient
		poll   *time.Ticker
	}
)

func NewSqlHandler(settings QueueSettings) *SqlQueueHandler {
	h := &SqlQueueHandler{
		queue:  settings.Queue,
		handle: HandlerSql,
		client: &sClient{},
	}

	if settings.Meta.PollDelay != nil {
		h.poll = time.NewTicker(*settings.Meta.PollDelay)
	}

	return h
}

// getTicker fetches the ticker channel if it is set-up in
// queue settings (see QueueSettingsMeta.PollDelay)
func (cq *SqlQueueHandler) Ticker(ctx context.Context) <-chan time.Time {
	if cq.poll != nil {
		return cq.poll.C
	}

	return nil
}

// getNotification fetches the notification channel from the store
// notification mechanism
func (cq *SqlQueueHandler) Notification(ctx context.Context) <-chan interface{} {
	// @todo - psql, redis (mysql only with a plugin)
	return nil
}

func (cq *SqlQueueHandler) Read(ctx context.Context) (QueueMessageSet, error) {
	return cq.client.get(ctx, cq.queue)
}

func (cq *SqlQueueHandler) Write(ctx context.Context, p []byte) error {
	return cq.client.add(ctx, cq.queue, p)
}

func (cq *SqlQueueHandler) SetStorer(s QueueStorer) {
	cq.client.SetStorer(s)
}

// Process marks a message as processed in store of choice
func (cq *SqlQueueHandler) Process(ctx context.Context, m QueueMessage) error {
	return cq.client.process(ctx, m)
}
