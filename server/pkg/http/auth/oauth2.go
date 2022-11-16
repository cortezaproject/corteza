package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type (
	ServicerOauth2 struct {
		c *http.Client

		client   string
		secret   string
		scope    []string
		tokenUrl *url.URL
	}

	Oauth2Params struct {
		Client   string
		Secret   string
		Scope    []string
		TokenUrl *url.URL
	}
)

func NewOauth2(p Oauth2Params, c *http.Client, s types.SecureStorager) (ss ServicerOauth2, err error) {
	if p.Client == "" {
		err = fmt.Errorf("invalid param client")
		return
	}

	if p.Secret == "" {
		err = fmt.Errorf("invalid param secret")
		return
	}

	if p.TokenUrl == nil || p.TokenUrl.String() == "" {
		err = fmt.Errorf("invalid param token url")
		return
	}

	ss = ServicerOauth2{
		c:        c,
		client:   p.Client,
		secret:   p.Secret,
		scope:    p.Scope,
		tokenUrl: p.TokenUrl,
	}

	return
}

func (s ServicerOauth2) Do(ctx context.Context) (t *oauth2.Token, err error) {
	c := &clientcredentials.Config{
		ClientID:     s.client,
		ClientSecret: s.secret,
		Scopes:       s.scope,
		TokenURL:     s.tokenUrl.String(),
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.c)

	return c.Token(ctx)
}
