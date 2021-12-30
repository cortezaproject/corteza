package gig

import (
	"context"
	"time"
)

type (
	WorkerStatus struct {
		Error       error         `json:"error,omitempty"`
		StartedAt   *time.Time    `json:"startedAt"`
		CompletedAt *time.Time    `json:"completedAt"`
		Elapsed     time.Duration `json:"elapsed"`
		Meta        WorkMeta      `json:"meta,omitempty"`
	}

	WorkMeta    map[string]interface{}
	WorkerState interface{}

	Worker interface {
		Prepare(context.Context, ...Source) error
		Preprocess(context.Context, ...Preprocessor) error
		Exec(context.Context) (SourceSet, WorkMeta, error)
		State(context.Context) (WorkerState, error)
		Cleanup(context.Context) error
		Ref() string
	}
)
