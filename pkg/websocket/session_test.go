package websocket

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSession_procRawMessage(t *testing.T) {
	var (
		req             = require.New(t)
		s               = session{server: Server(nil, options.WebsocketOpt{})}
		jwtHandler, err = auth.TokenManager("secret", time.Minute)

		userID uint64 = 123
		token  []byte

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

	req.NoError(err)

	token, err = jwtHandler.Encode(auth.Authenticated(userID, 456, 789), 0, "api")
	req.NoError(err)

	req.EqualError(s.procRawMessage([]byte("{}")), "unauthenticated session")
	req.Nil(s.identity)

	req.EqualError(s.procRawMessage(mockResponse(nil)), "unauthorized: failed to parse token: EOF")
	req.Nil(s.identity)

	req.NoError(s.procRawMessage(mockResponse(token)))
	req.NotNil(s.identity)
	req.Equal(userID, s.identity.Identity())

	req.EqualError(s.procRawMessage([]byte("{}")), "unknown message type ''")
	req.Equal(userID, s.identity.Identity())

	// Repeat with the same user
	token, err = jwtHandler.Encode(auth.Authenticated(userID, 456, 789), 0, "api")
	req.NoError(err)

	req.NoError(s.procRawMessage(mockResponse(token)))
	req.NotNil(s.identity)
	req.Equal(userID, s.identity.Identity())

	// Try to authenticate on an existing authenticated session as a different user
	token, err = jwtHandler.Encode(auth.Authenticated(userID+1, 456, 789), 0, "api")
	req.NoError(err)

	req.EqualError(s.procRawMessage(mockResponse(token)), "unauthorized: identity does not match")

}
