package websocket

import (
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
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

	req.NoError(err)
	s.server.accessToken = jwtHandler

	jwt := jwtHandler.Encode(auth.NewIdentity(userID, 456, 789))

	req.EqualError(s.procRawMessage([]byte("{}")), "empty payload")
	req.Nil(s.identity)

	req.EqualError(s.procRawMessage([]byte(`{"auth":{}}`)), "unauthorized: token contains an invalid number of segments")
	req.Nil(s.identity)

	req.EqualError(s.procRawMessage([]byte(`{"auth":{"access_token": ""}}`)), "unauthorized: token contains an invalid number of segments")
	req.Nil(s.identity)

	req.NoError(s.procRawMessage([]byte(`{"auth":{"access_token": "` + jwt + `"}}`)))
	req.NotNil(s.identity)
	req.Equal(userID, s.identity.Identity())

	req.EqualError(s.procRawMessage([]byte("{}")), "empty payload")
	req.Equal(userID, s.identity.Identity())

	// Repeat with the same user
	jwt = jwtHandler.Encode(auth.NewIdentity(userID, 456, 789))

	req.NoError(s.procRawMessage([]byte(`{"auth":{"access_token": "` + jwt + `"}}`)))
	req.NotNil(s.identity)
	req.Equal(userID, s.identity.Identity())

	// Try to authenticate on an existing authenticated session as a different user
	jwt = jwtHandler.Encode(auth.NewIdentity(userID+1, 456, 789))

	req.EqualError(s.procRawMessage([]byte(`{"auth":{"access_token": "`+jwt+`"}}`)), "unauthorized: identity does not match")

}
