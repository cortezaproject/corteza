package types

import (
	"context"
)

const (
	ConsumerCorteza  ConsumerType = "corteza"
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
		ConsumerCorteza,
		ConsumerEventbus,
		ConsumerRedis,
		ConsumerStore,
	}
}
