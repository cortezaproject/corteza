package service

import (
	"context"
)

type (
	suspender interface {
		Suspend(context.Context, uint64) error
		Unsuspend(context.Context, uint64) error
	}
	archiver interface {
		Archive(context.Context, uint64) error
		Unarchive(context.Context, uint64) error
	}

	deleter interface {
		Delete(context.Context, uint64) error
	}
)
