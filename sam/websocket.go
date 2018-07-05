package sam

import (
	"github.com/pkg/errors"
)

var _ = errors.Wrap

func (*Websocket) Client(r *websocketClientRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Websocket.client")
}
