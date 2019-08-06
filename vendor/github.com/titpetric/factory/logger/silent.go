package logger

import (
	"context"
)

type (
	Silent struct{}
)

func (Silent) Log(ctx context.Context, msg string, fields ...Field) {
}
