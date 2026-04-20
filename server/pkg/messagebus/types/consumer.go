package types

import (
	"context"
)

const (
	ConsumerHuman  ConsumerType = "human"
	ConsumerNoop     ConsumerType = "noop"
	ConsumerRedis    ConsumerType = "redis"
	ConsumerStore    ConsumerType = "store"
	ConsumerEventbus ConsumerType = "eventbus"
)

type (
	ConsumerType string

	Consumer interface {
		Write(ctx context.Context, p []byte) error
	}
)

func ConsumerTypes() []ConsumerType {
	return []ConsumerType{
		ConsumerHuman,
		ConsumerEventbus,
		ConsumerRedis,
		ConsumerStore,
	}
}
