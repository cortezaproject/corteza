package messagebus

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type (
	RedisQueueHandler struct {
		queue  string
		handle handler
		client *redis.Client
		poll   *time.Ticker
	}

	redisSettings struct {
		DSN      string `json:"dsn"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	}
)

func NewRedisHandler(settings QueueSettings) *RedisQueueHandler {
	rdb := redis.NewClient(&redis.Options{
		// Addr:     settings.Meta.(redisSettings).DSN,
		// Password: settings.Meta.(redisSettings).Password,
		// DB:       settings.Meta.(redisSettings).Db,
	})

	h := &RedisQueueHandler{
		queue:  settings.Queue,
		handle: HandlerRedis,
		client: rdb,
	}

	if settings.Meta.PollDelay != nil {
		h.poll = time.NewTicker(*settings.Meta.PollDelay)
	}

	return h
}

// getTicker fetches the ticker channel if it is set-up in
// queue settings (see QueueSettingsMeta.PollDelay)
func (cq *RedisQueueHandler) Ticker(ctx context.Context) <-chan time.Time {
	if cq.poll != nil {
		return cq.poll.C
	}

	return nil
}

// getNotification fetches the notification channel from the store
// notification mechanism
func (cq *RedisQueueHandler) Notification(ctx context.Context) <-chan interface{} {
	// @todo - psql, redis (mysql only with a plugin)
	return nil
}

func (cq *RedisQueueHandler) Read(ctx context.Context) (set QueueMessageSet, err error) {
	var (
		vals []string
	)

	vals, err = cq.client.LRange(ctx, cq.queue, 0, -1).Result()

	if err != nil {
		return
	}

	for _, v := range vals {
		set = append(set, &QueueMessage{
			Queue:   cq.queue,
			Payload: []byte(v),
		})
	}

	return
}

func (cq *RedisQueueHandler) Write(ctx context.Context, p []byte) error {
	return cq.client.LPush(ctx, cq.queue, p).Err()
}

// func (cq *RedisQueueHandler) Write(ctx context.Context, in chan []byte, e chan error) {
// 	for p := range in {

// 		err :=

// 		if err != nil {
// 			// do nothing for now
// 		}
// 	}
// 	return
// }

// func (cq *RedisQueueHandler) Read(ctx context.Context, out chan []byte, e chan error) {
// 	vals, err := cq.client.LRange(ctx, cq.queue, 0, -1).Result()
// 	// val, err := cq.client.RPop(ctx, "queue").Result()

// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, v := range vals {
// 		out <- []byte(v)
// 	}

// 	close(out)
// }
