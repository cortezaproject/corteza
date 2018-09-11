package service

import (
	"github.com/crusttech/crust/sam/repository"
)

type PubSub struct {
	repository.PubSubClient
}

func (PubSub) New() (*PubSub, error) {
	rpo, err := repository.PubSub{}.New()
	if err != nil {
		return nil, err
	}
	return &PubSub{rpo}, nil
}
