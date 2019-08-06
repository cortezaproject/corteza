package logger

import (
	"context"
	"fmt"
)

type (
	Logger interface {
		Log(ctx context.Context, message string, fields ...Field)
	}

	Field interface {
		Name() string
		Value() interface{}

		fmt.Stringer
	}
)

func New() Logger {
	return Default{}
}
