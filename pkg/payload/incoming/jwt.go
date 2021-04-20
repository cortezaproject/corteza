package incoming

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	Token struct {
		AccessToken *string `json:"access_token"`
	}
)

func (t *Token) ParseWithClaims() (jwt.MapClaims, error) {
	token, err := jwt.Parse(*t.AccessToken, nil)
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
