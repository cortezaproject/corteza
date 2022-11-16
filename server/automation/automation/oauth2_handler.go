package automation

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/http/auth"
	"golang.org/x/oauth2"
)

type (
	secureStoragerTodo struct{}

	oauth2Handler struct {
		reg oauth2HandlerRegistry
	}
)

func Oauth2Handler(reg oauth2HandlerRegistry) *oauth2Handler {
	h := &oauth2Handler{
		reg: reg,
	}

	h.register()
	return h
}

func (h oauth2Handler) authenticate(ctx context.Context, args *oauth2AuthenticateArgs) (res *oauth2AuthenticateResults, err error) {
	var (
		u     *url.URL
		token *oauth2.Token
	)

	if !args.hasClient {
		err = fmt.Errorf("could not init auth, client key missing")
		return
	}

	if !args.hasSecret {
		err = fmt.Errorf("could not init auth, secret key missing")
		return
	}

	if !args.hasScope {
		err = fmt.Errorf("could not init auth, scope missing")
		return
	}

	if !args.hasTokenUrl {
		err = fmt.Errorf("could not init auth, token url missing")
		return
	}

	if u, err = url.Parse(args.TokenUrl); err != nil {
		err = fmt.Errorf("could not init auth, token url invalid: %s", err)
		return
	}

	params := auth.Oauth2Params{
		Client:   args.Client,
		Secret:   args.Secret,
		Scope:    strings.Fields(args.Scope),
		TokenUrl: u,
	}

	servicer, err := auth.NewOauth2(params, http.DefaultClient, secureStoragerTodo{})

	if err != nil {
		err = fmt.Errorf("could not init servicer: %s", err)
		return
	}

	// call servicer, fetch the token
	if token, err = servicer.Do(ctx); err != nil {
		return
	}

	res = &oauth2AuthenticateResults{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Token:        token,
	}

	return
}
