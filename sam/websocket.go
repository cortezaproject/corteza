package sam

import (
	"github.com/pkg/errors"
)

func (*Websocket) Client(r *websocketClientRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Websocket.client")
}
