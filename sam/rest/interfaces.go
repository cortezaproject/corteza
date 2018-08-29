package rest

import (
	"context"
)

type (
	suspender interface {
		Suspend(ctx context.Context, ID uint64) error
		Unsuspend(ctx context.Context, ID uint64) error
	}

	archiver interface {
		Archive(ctx context.Context, ID uint64) error
		Unarchive(ctx context.Context, ID uint64) error
	}

	deleter interface {
		Delete(ctx context.Context, ID uint64) error
	}
)
