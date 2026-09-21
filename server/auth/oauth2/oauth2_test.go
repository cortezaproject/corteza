package oauth2

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/options"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	testClientStore struct{ client oauth2.ClientInfo }
	testTokenStore  struct{ tokens map[string]oauth2.TokenInfo }
)

func (s testClientStore) GetByID(context.Context, string) (oauth2.ClientInfo, error) {
	return s.client, nil
}

func (s *testTokenStore) Create(_ context.Context, ti oauth2.TokenInfo) error {
	s.tokens[ti.GetCode()] = ti
	return nil
}

func (s *testTokenStore) GetByCode(_ context.Context, code string) (oauth2.TokenInfo, error) {
	return s.tokens[code], nil
}

func (s *testTokenStore) RemoveByCode(_ context.Context, code string) error {
	delete(s.tokens, code)
	return nil
}

func (*testTokenStore) RemoveByAccess(context.Context, string) error  { return nil }
func (*testTokenStore) RemoveByRefresh(context.Context, string) error { return nil }
func (*testTokenStore) GetByAccess(context.Context, string) (oauth2.TokenInfo, error) {
	return nil, nil
}
func (*testTokenStore) GetByRefresh(context.Context, string) (oauth2.TokenInfo, error) {
	return nil, nil
}

func TestManager_redirectURI(t *testing.T) {
	var (
		ctx = context.Background()
		opt = options.AuthOpt{BaseURL: "http://localhost:8090/auth"}

		// default client, provisioned without redirect URI
		cs = testClientStore{client: &models.Client{ID: "42", Secret: "secret", Domain: ""}}
		m  = NewManager(opt, zap.NewNop(), cs, &testTokenStore{tokens: make(map[string]oauth2.TokenInfo)})

		authorize = func(redirectURI string) (oauth2.TokenInfo, error) {
			return m.GenerateAuthToken(ctx, oauth2.Code, &oauth2.TokenGenerateRequest{
				ClientID:    "42",
				UserID:      "1",
				RedirectURI: redirectURI,
				Scope:       "api",
			})
		}
	)

	t.Run("authorization code is issued for own origin", func(t *testing.T) {
		ti, err := authorize("http://localhost:8090/compose/auth")
		require.NoError(t, err)
		require.NotEmpty(t, ti.GetCode())
	})

	for _, redirectURI := range []string{
		"http://evil.com/steal",
		"http://localhost:8090@evil.com/steal",
		"http://localhost:8090.evil.com/steal",
	} {
		t.Run("authorization code is not issued for "+redirectURI, func(t *testing.T) {
			_, err := authorize(redirectURI)
			require.ErrorIs(t, err, errors.ErrInvalidRedirectURI)
		})

		t.Run("code is not exchanged for "+redirectURI, func(t *testing.T) {
			_, err := m.GenerateAccessToken(ctx, oauth2.AuthorizationCode, &oauth2.TokenGenerateRequest{
				ClientID:     "42",
				ClientSecret: "secret",
				Code:         "stolen",
				RedirectURI:  redirectURI,
			})
			require.ErrorIs(t, err, errors.ErrInvalidRedirectURI)
		})
	}
}
