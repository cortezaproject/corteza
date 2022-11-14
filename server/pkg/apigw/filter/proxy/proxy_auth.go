package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/http/auth"
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
		servicer auth.ServicerBasic

		user string
		pass string
	}

	proxyAuthServicerOauth2 struct {
		c        *http.Client
		servicer auth.ServicerOauth2

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

	ProxyAuthParams struct {
		Type   proxyAuthType          `json:"type"`
		Params map[string]interface{} `json:"params"`
	}

	proxyAuthMetaArg struct {
		Label   string                 `json:"label"`
		Type    string                 `json:"type"`
		Example string                 `json:"example,omitempty"`
		Options map[string]interface{} `json:"options,omitempty"`
	}

	ProxyAuthDefinition struct {
		Type   proxyAuthType       `json:"type"`
		Params []*proxyAuthMetaArg `json:"params"`
	}
)

func NewProxyAuthServicer(c *http.Client, p ProxyAuthParams, s types.SecureStorager) (ProxyAuthServicer, error) {
	switch p.Type {
	case proxyAuthTypeHeader:
		return newProxyAuthHeader(p)
	case proxyAuthTypeQuery:
		return newProxyAuthQuery(p)
	case proxyAuthTypeBasic:
		return newProxyAuthBasic(p)
	case proxyAuthTypeOauth2:
		return newProxyAuthOauth2(p, c, s)
	case proxyAuthTypeJWT:
		return newProxyAuthJWT(p)
	default:
		return proxyAuthServicerNoop{}, nil
	}
}

func newProxyAuthHeader(p ProxyAuthParams) (s proxyAuthServicerHeader, err error) {
	s = proxyAuthServicerHeader{
		params: p.Params,
	}

	return
}

func newProxyAuthQuery(p ProxyAuthParams) (s proxyAuthServicerQuery, err error) {
	s = proxyAuthServicerQuery{
		params: p.Params,
	}

	return
}

func newProxyAuthBasic(p ProxyAuthParams) (s proxyAuthServicerBasic, err error) {
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

	servicer, _ := auth.NewBasic(auth.BasicParams{
		User: user,
		Pass: pass,
	})

	s = proxyAuthServicerBasic{user: user, pass: pass, servicer: servicer}

	return
}

func newProxyAuthOauth2(p ProxyAuthParams, c *http.Client, s types.SecureStorager) (ss proxyAuthServicerOauth2, err error) {
	var (
		u                        *url.URL
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

	if u, err = url.Parse(p.Params["token_url"].(string)); err != nil {
		err = fmt.Errorf("invalid param token url: %s", err)
		return
	}

	servicer, err := auth.NewOauth2(auth.Oauth2Params{
		Client:   client,
		Secret:   secret,
		Scope:    scope,
		TokenUrl: u,
	}, c, s)

	ss = proxyAuthServicerOauth2{
		c:        c,
		client:   client,
		servicer: servicer,
		secret:   secret,
		scope:    scope,
		tokenUrl: tokenUrl,
	}

	return
}

func newProxyAuthJWT(p ProxyAuthParams) (ss proxyAuthServicerJWT, err error) {
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
	r.Header.Set("Authorization", "Basic "+s.servicer.Do(r.Context()))
	return nil
}

func (s proxyAuthServicerOauth2) Do(r *http.Request) error {
	t, err := s.servicer.Do(r.Context())

	if err != nil {
		return err
	}

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

func (s proxyAuthServicerOauth2) Def() *ProxyAuthDefinition {
	return &ProxyAuthDefinition{
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

func (s proxyAuthServicerHeader) Def() *ProxyAuthDefinition {
	return &ProxyAuthDefinition{
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

func (s proxyAuthServicerBasic) Def() *ProxyAuthDefinition {
	return &ProxyAuthDefinition{
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

func (s proxyAuthServicerJWT) Def() *ProxyAuthDefinition {
	return &ProxyAuthDefinition{
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

func ProxyAuthDef() []*ProxyAuthDefinition {
	return []*ProxyAuthDefinition{
		proxyAuthServicerOauth2{}.Def(),
		proxyAuthServicerHeader{}.Def(),
		proxyAuthServicerBasic{}.Def(),
		proxyAuthServicerJWT{}.Def(),
	}
}
