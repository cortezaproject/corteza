package auth

import (
	"context"
	"encoding/base64"
	"fmt"
)

type (
	ServicerBasic struct {
		user string
		pass string
	}

	BasicParams struct {
		User string
		Pass string
	}
)

func NewBasic(p BasicParams) (s ServicerBasic, err error) {
	if p.User == "" {
		err = fmt.Errorf("invalid param username")
		return
	}

	if p.Pass == "" {
		err = fmt.Errorf("invalid param password")
		return
	}

	s = ServicerBasic{user: p.User, pass: p.Pass}

	return
}

func (s ServicerBasic) Do(ctx context.Context) string {
	auth := s.user + ":" + s.pass
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
