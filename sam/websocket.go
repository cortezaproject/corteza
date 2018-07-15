package sam

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/sam/rest"
	_ "github.com/crusttech/crust/sam/types"
)

var _ = errors.Wrap

type Websocket struct{}

func (Websocket) New() *Websocket {
	return &Websocket{}
}

func (*Websocket) Client(r *rest.WebsocketClientRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Websocket.client")
}
