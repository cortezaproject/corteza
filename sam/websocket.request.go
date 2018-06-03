package sam

import (
	"net/http"

	"github.com/pkg/errors"
)

// Websocket client request parameters
type WebsocketClientRequest struct {
}

func (WebsocketClientRequest) new() *WebsocketClientRequest {
	return &WebsocketClientRequest{}
}

func (w *WebsocketClientRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}
	return errors.New("Not implemented: WebsocketClientRequest.Fill")
}

var _ RequestFiller = WebsocketClientRequest{}.new()
