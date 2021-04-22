package websocket

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	Payload struct {
		// Auth is JWT token provided by client as first message,
		// and will be passed whenever it changes
		*Auth `json:"auth"`
	}

	Auth struct {
		AccessToken *string `json:"access_token"`
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

func Unmarshal(raw []byte) (*Payload, error) {
	var p Payload
	return &p, json.Unmarshal(raw, &p)
}
