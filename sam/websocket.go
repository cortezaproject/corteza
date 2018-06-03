package sam

import (
	"github.com/pkg/errors"
)

func (w *Websocket) Client(r *websocketClientRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Websocket.client")
}
