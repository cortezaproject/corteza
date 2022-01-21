package websocket

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSession_procRawMessage(t *testing.T) {
	var (
		req = require.New(t)

		identity1 = auth.Authenticated(123, 456, 789)
		identity2 = auth.Authenticated(321, 456, 789)

		s = session{
			server: Server(
				nil,
				options.WebsocketOpt{},
				func(ctx context.Context, accessToken string) (auth.Identifiable, error) {
					//token, err := jwt.Parse([]byte(accessToken))
					//if err != nil {
					//	return nil, err
					//}
					//return auth.IdentityFromToken(token), nil
					switch accessToken {
					case "one":
						return identity1, nil
					case "two":
						return identity2, nil
					case "":
						return nil, fmt.Errorf("failed to parse token: EOF")
					}

					return nil, fmt.Errorf("something else went wrong")
				}),
		}

		token []byte

		mockResponse = func(token []byte) (out []byte) {
			out = []byte(`{"@type": "credentials", "@value": {"accessToken": "`)
			out = append(out, token...)
			out = append(out, []byte(`"}}`)...)
			return
		}
	)

	if testing.Verbose() {
		s.logger = logger.MakeDebugLogger()
	} else {
		s.logger = zap.NewNop()
	}

	req.EqualError(s.procRawMessage([]byte("{}")), "unauthenticated session")
	req.Nil(s.identity)

	req.EqualError(s.procRawMessage(mockResponse(nil)), "unauthorized: failed to parse token: EOF")
	req.Nil(s.identity)

	token = []byte("one")
	req.NoError(s.procRawMessage(mockResponse(token)))
	req.NotNil(s.identity)
	req.Equal(identity1.Identity(), s.identity.Identity())

	req.EqualError(s.procRawMessage([]byte("{}")), "unknown message type ''")
	req.Equal(identity1.Identity(), s.identity.Identity())

	token = []byte("one")
	req.NoError(s.procRawMessage(mockResponse(token)))
	req.NotNil(s.identity)
	req.Equal(identity1.Identity(), s.identity.Identity())

	token = []byte("two")
	req.EqualError(s.procRawMessage(mockResponse(token)), "unauthorized: identity does not match")
}
