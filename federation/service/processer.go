package service

import "context"

type (
	Processer interface {
		Process(ctx context.Context, payload []byte) (ProcesserResponse, error)
	}

	ProcesserResponse interface{}
)
