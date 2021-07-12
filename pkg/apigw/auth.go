package apigw

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

const (
	authTypeOauth2 authType = "oauth2"
	authTypeHeader authType = "header"
	authTypeQuery  authType = "query"
	authTypeBasic  authType = "basic"
	authTypeJwt    authType = "jwt"
	authTypeNoop   authType = "noop"
)

type (
	AuthServicer interface {
		Do(*http.Request) error
	}

	authServicerNoop struct{}

	authServicerHeader struct {
		params map[string]interface{}
	}

	authServicerQuery struct {
		params map[string]interface{}
	}

	authServicerBasic struct {
		user string
		pass string
	}

	authType string

	authParams struct {
		Type   authType               `json:"type"`
		Params map[string]interface{} `json:"params"`
	}
)

func NewAuthHeader(p authParams) (s authServicerHeader, err error) {
	s = authServicerHeader{
		params: p.Params,
	}

	return
}

func NewAuthQuery(p authParams) (s authServicerQuery, err error) {
	s = authServicerQuery{
		params: p.Params,
	}

	return
}

func NewAuthBasic(p authParams) (s authServicerBasic, err error) {
	var (
		ok         bool
		user, pass string
	)

	if user, ok = p.Params["username"].(string); !ok {
		err = fmt.Errorf("invalid param username")
		return
	}

	if pass, ok = p.Params["password"].(string); !ok {
		err = fmt.Errorf("invalid param password")
		return
	}

	s = authServicerBasic{user: user, pass: pass}

	return
}

func NewAuthServicer(c *http.Client, p authParams) (AuthServicer, error) {
	switch p.Type {
	case authTypeHeader:
		return NewAuthHeader(p)
	case authTypeQuery:
		return NewAuthQuery(p)
	case authTypeBasic:
		return NewAuthBasic(p)
	default:
		return authServicerNoop{}, nil
	}
}

func (s authServicerHeader) Do(r *http.Request) error {
	for k, v := range s.params {
		r.Header.Add(k, v.(string))
	}
	return nil
}

func (s authServicerQuery) Do(r *http.Request) error {
	if len(s.params) == 0 {
		return nil
	}

	q := r.URL.Query()

	for k, v := range s.params {
		q.Set(k, v.(string))
	}

	r.URL.RawQuery = q.Encode()

	return nil
}

func (s authServicerBasic) Do(r *http.Request) error {
	bs := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.user, s.pass)))

	r.Header.Set("Authorization", fmt.Sprintf("Basic %s", bs))
	return nil
}

func (s authServicerNoop) Do(r *http.Request) error {
	return nil
}
