package rest

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/sam/rest/server"
	_ "github.com/crusttech/crust/sam/types"
)

var _ = errors.Wrap

type Websocket struct{}

func (Websocket) New() *Websocket {
	return &Websocket{}
}

func (*Websocket) Client(r *server.WebsocketClientRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Websocket.client")
}
