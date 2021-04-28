package websocket

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	// Auth is JWT token provided by client as first message,
	// and will be passed whenever it changes
	Auth struct {
		AccessToken *string `json:"access_token"`
	}

	// payload for incoming messages from user
	payload struct {
		*Auth `json:"auth"`
	}

	// response for sending messages to user
	response struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}
)

func (a *Auth) ParseWithClaims() (jwt.MapClaims, error) {
	token, err := jwt.Parse(*a.AccessToken, nil)
	if token == nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	} else {
		return nil, errors.New("Invalid token")
	}
}

func Unmarshal(raw []byte) (*payload, error) {
	var p payload
	return &p, json.Unmarshal(raw, &p)
}

func Response(status string, data interface{}) *response {
	return &response{
		Status: status,
		Data:   data,
	}
}

func (m response) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
