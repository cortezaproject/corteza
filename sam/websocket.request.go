package sam

import (
	"net/http"
)

// Websocket client request parameters
type websocketClientRequest struct {
}

func (websocketClientRequest) new() *websocketClientRequest {
	return &websocketClientRequest{}
}

func (w *websocketClientRequest) Fill(r *http.Request) error {
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
	return nil
}

var _ RequestFiller = websocketClientRequest{}.new()
