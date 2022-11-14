package websocket

import (
	"encoding/json"
)

type (
	// Auth is JWT token provided by client as first message,
	// and will be passed whenever it changes
	payloadAuth struct {
		AccessToken string `json:"accessToken"`
	}

	payloadWrap struct {
		Type  string          `json:"@type"`
		Value json.RawMessage `json:"@value"`
	}
)

const (
	payloadTypeCredentials = "credentials"
)

var (
	closingUnidentifiedConn, _ = MarshalPayload("error", "closing unidentified connection")
	ok, _                      = MarshalPayload("message", "authenticated")
)

func (p payloadWrap) UnmarshalValue(m interface{}) error {
	return json.Unmarshal(p.Value, m)
}

func MarshalPayload(t string, m interface{}) ([]byte, error) {
	var (
		err error
		w   = payloadWrap{Type: t}
	)

	if w.Value, err = json.Marshal(m); err != nil {
		return nil, err
	}

	return json.Marshal(w)
}
