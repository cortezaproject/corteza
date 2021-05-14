package websocket

import (
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestSession_procRawMessage(t *testing.T) {
	var (
		req             = require.New(t)
		s               = session{server: Server(nil, options.WebsocketOpt{})}
		jwtHandler, err = auth.JWT("secret", time.Minute)

		userID uint64 = 123
	)

	if testing.Verbose() {
		s.logger = logger.MakeDebugLogger()
	} else {
		s.logger = zap.NewNop()
	}

	req.NoError(err)
	s.server.accessToken = jwtHandler

	jwt := jwtHandler.Encode(auth.NewIdentity(userID, 456, 789))

	req.EqualError(s.procRawMessage([]byte("{}")), "unauthenticated session")
	req.Nil(s.identity)

	req.EqualError(s.procRawMessage([]byte(`{"@type": "credentials", "@value": {"accessToken": ""}}`)), "unauthorized: token contains an invalid number of segments")
	req.Nil(s.identity)

	req.NoError(s.procRawMessage([]byte(`{"@type": "credentials", "@value": {"accessToken": "` + jwt + `"}}`)))
	req.NotNil(s.identity)
	req.Equal(userID, s.identity.Identity())

	req.EqualError(s.procRawMessage([]byte("{}")), "unknown message type ''")
	req.Equal(userID, s.identity.Identity())

	// Repeat with the same user
	jwt = jwtHandler.Encode(auth.NewIdentity(userID, 456, 789))

	req.NoError(s.procRawMessage([]byte(`{"@type": "credentials", "@value": {"accessToken": "` + jwt + `"}}`)))
	req.NotNil(s.identity)
	req.Equal(userID, s.identity.Identity())

	// Try to authenticate on an existing authenticated session as a different user
	jwt = jwtHandler.Encode(auth.NewIdentity(userID+1, 456, 789))

	req.EqualError(s.procRawMessage([]byte(`{"@type": "credentials", "@value": {"accessToken": "`+jwt+`"}}`)), "unauthorized: identity does not match")

}
