package service

import (
	"context"

	"github.com/crusttech/crust/sam/repository"
)

type PubSub struct {
	*repository.PubSub
}

func (PubSub) New(ctx context.Context) (*PubSub, error) {
	rpo, err := repository.PubSub{}.New(ctx)
	if err != nil {
		return nil, err
	}
	return &PubSub{rpo}, nil
}
