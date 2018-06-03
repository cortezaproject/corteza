package sam

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (wh *WebsocketHandlers) Client(w http.ResponseWriter, r *http.Request) {
	params := websocketClientRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return wh.Websocket.Client(params) })
}
