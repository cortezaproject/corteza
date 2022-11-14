package service

import (
	"context"
	"time"
)

type (
	worker interface {
		PrepareForNodes(ctx context.Context, sync *Sync, urls chan Url)
		Watch(ctx context.Context, delay time.Duration, limit int)
	}
)

const (
	defaultDelay = time.Second * 10
	defaultPage  = 10
)
