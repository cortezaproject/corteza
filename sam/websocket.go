package sam

import (
	"github.com/pkg/errors"
)

func (w *Websocket) Client(r *WebsocketClientRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Websocket.client")
}
