package apigw

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	proxyAuthTypeOauth2 proxyAuthType = "oauth2"
	proxyAuthTypeHeader proxyAuthType = "header"
	proxyAuthTypeQuery  proxyAuthType = "query"
	proxyAuthTypeBasic  proxyAuthType = "basic"
	proxyAuthTypeJWT    proxyAuthType = "jwt"
	proxyAuthTypeNoop   proxyAuthType = "noop"
)

type (
	ProxyAuthServicer interface {
		Do(*http.Request) error
	}

	proxyAuthServicerNoop struct{}

	proxyAuthServicerHeader struct {
		params map[string]interface{}
	}

	proxyAuthServicerQuery struct {
		params map[string]interface{}
	}

	proxyAuthServicerBasic struct {
		user string
		pass string
	}

	proxyAuthServicerOauth2 struct {
		c *http.Client

		client   string
		secret   string
		scope    []string
		tokenUrl string

		Args []*proxyAuthMetaArg `json:"params"`
	}

	proxyAuthServicerJWT struct {
		JWT string
	}

	proxyAuthType string

	proxyAuthParams struct {
		Type   proxyAuthType          `json:"type"`
		Params map[string]interface{} `json:"params"`
	}

	proxyAuthMetaArg struct {
		Label   string                 `json:"label"`
		Type    string                 `json:"type"`
		Example string                 `json:"example,omitempty"`
		Options map[string]interface{} `json:"options,omitempty"`
	}

	proxyAuthDefinition struct {
		Type   proxyAuthType       `json:"type"`
		Params []*proxyAuthMetaArg `json:"params"`
	}
)

func NewProxyAuthHeader(p proxyAuthParams) (s proxyAuthServicerHeader, err error) {
	s = proxyAuthServicerHeader{
		params: p.Params,
	}

	return
}

func NewProxyAuthQuery(p proxyAuthParams) (s proxyAuthServicerQuery, err error) {
	s = proxyAuthServicerQuery{
		params: p.Params,
	}

	return
}

func NewProxyAuthBasic(p proxyAuthParams) (s proxyAuthServicerBasic, err error) {
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

	s = proxyAuthServicerBasic{user: user, pass: pass}

	return
}

func NewProxyAuthOauth2(p proxyAuthParams, c *http.Client, s SecureStorager) (ss proxyAuthServicerOauth2, err error) {
	var (
		ok                       bool
		client, secret, tokenUrl string
		scope                    []string = []string{}
	)

	if client, ok = p.Params["client"].(string); !ok {
		err = fmt.Errorf("invalid param client")
		return
	}

	if secret, ok = p.Params["secret"].(string); !ok {
		err = fmt.Errorf("invalid param secret")
		return
	}

	if scopes, ok := p.Params["scope"].(string); ok {
		scope = strings.Fields(scopes)
	}

	if tokenUrl, ok = p.Params["token_url"].(string); !ok {
		err = fmt.Errorf("invalid param token url")
		return
	}

	ss = proxyAuthServicerOauth2{
		c:        c,
		client:   client,
		secret:   secret,
		scope:    scope,
		tokenUrl: tokenUrl,
	}

	return
}

func NewProxyAuthJWT(p proxyAuthParams) (ss proxyAuthServicerJWT, err error) {
	var (
		ok  bool
		jwt string
	)

	if jwt, ok = p.Params["jwt"].(string); !ok {
		err = fmt.Errorf("invalid param jwt")
		return
	}

	ss = proxyAuthServicerJWT{
		JWT: jwt,
	}

	return
}

func NewProxyAuthServicer(c *http.Client, p proxyAuthParams, s SecureStorager) (ProxyAuthServicer, error) {
	switch p.Type {
	case proxyAuthTypeHeader:
		return NewProxyAuthHeader(p)
	case proxyAuthTypeQuery:
		return NewProxyAuthQuery(p)
	case proxyAuthTypeBasic:
		return NewProxyAuthBasic(p)
	case proxyAuthTypeOauth2:
		return NewProxyAuthOauth2(p, c, s)
	case proxyAuthTypeJWT:
		return NewProxyAuthJWT(p)
	default:
		return proxyAuthServicerNoop{}, nil
	}
}

func (s proxyAuthServicerHeader) Do(r *http.Request) error {
	for k, v := range s.params {
		r.Header.Add(k, v.(string))
	}
	return nil
}

func (s proxyAuthServicerQuery) Do(r *http.Request) error {
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

func (s proxyAuthServicerBasic) Do(r *http.Request) error {
	r.SetBasicAuth(s.user, s.pass)
	return nil
}

func (s proxyAuthServicerOauth2) Do(r *http.Request) error {
	c := &clientcredentials.Config{
		ClientID:     s.client,
		ClientSecret: s.secret,
		Scopes:       s.scope,
		TokenURL:     s.tokenUrl,
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.c)

	t, err := c.Token(ctx)

	if err != nil {
		return err
	}

	// todo - store this to secure storager
	r.Header.Set("Authorization", "Bearer "+t.AccessToken)

	return nil
}

func (s proxyAuthServicerJWT) Do(r *http.Request) (err error) {
	// set up the request and get the jwt
	r.Header.Set("Authorization", "Bearer "+s.JWT)
	return
}

func (s proxyAuthServicerNoop) Do(r *http.Request) error {
	return nil
}

func (s proxyAuthServicerOauth2) Def() *proxyAuthDefinition {
	return &proxyAuthDefinition{
		Type: proxyAuthTypeOauth2,
		Params: []*proxyAuthMetaArg{
			{
				Label:   "client",
				Type:    "text",
				Example: "",
			},
			{
				Label:   "secret",
				Type:    "text",
				Example: "",
			},
			{
				Label:   "scope",
				Type:    "text",
				Example: "Keep the scopes separated with spaces ie.: 'scope1 scope2'",
			},
			{
				Label:   "token_url",
				Type:    "text",
				Example: "",
			},
		},
	}
}

func (s proxyAuthServicerHeader) Def() *proxyAuthDefinition {
	return &proxyAuthDefinition{
		Type: proxyAuthTypeHeader,
		Params: []*proxyAuthMetaArg{
			{
				Label:   "header",
				Type:    "text",
				Example: "",
			},
			{
				Label:   "key",
				Type:    "text",
				Example: "",
			},
		},
	}
}

func (s proxyAuthServicerBasic) Def() *proxyAuthDefinition {
	return &proxyAuthDefinition{
		Type: proxyAuthTypeBasic,
		Params: []*proxyAuthMetaArg{
			{
				Label:   "username",
				Type:    "text",
				Example: "",
			},
			{
				Label:   "password",
				Type:    "text",
				Example: "",
			},
		},
	}
}

func (s proxyAuthServicerJWT) Def() *proxyAuthDefinition {
	return &proxyAuthDefinition{
		Type: proxyAuthTypeJWT,
		Params: []*proxyAuthMetaArg{
			{
				Label:   "jwt",
				Type:    "text",
				Example: "",
			},
		},
	}
}

func ProxyAuthDef() []*proxyAuthDefinition {
	return []*proxyAuthDefinition{
		proxyAuthServicerOauth2{}.Def(),
		proxyAuthServicerHeader{}.Def(),
		proxyAuthServicerBasic{}.Def(),
		proxyAuthServicerJWT{}.Def(),
	}
}
