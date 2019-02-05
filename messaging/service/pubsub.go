package service

import (
	"github.com/crusttech/crust/messaging/repository"
)

type pubSub struct {
	repository.PubSubClient
}

func PubSub() *pubSub {
	return &pubSub{repository.PubSub{}.New()}
}
