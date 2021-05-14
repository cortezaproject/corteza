package websocket

import (
	"bytes"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestWebsocketSend_NoSessions(t *testing.T) {
	var (
		req = require.New(t)
		ws  = Server(zap.NewNop(), options.WebsocketOpt{})
	)

	req.NoError(ws.Send("msg", "msg"))
	req.NoError(ws.Send("msg", "msg", 1))
	req.NoError(ws.Send("msg", "msg", 1, 2))
	req.NoError(ws.Send("msg", "msg", 1, 2, 3))
}

func TestWebsocketSend_ExistingSessions(t *testing.T) {
	var (
		req = require.New(t)
		ws  = Server(zap.NewNop(), options.WebsocketOpt{})

		s1User uint64 = 100
		s1ID   uint64 = 101

		s2User uint64 = 200
		s2ID   uint64 = 201

		s1 = &bytes.Buffer{}
		s2 = &bytes.Buffer{}
	)

	ws.storeSession(s1, s1User, s1ID)
	ws.storeSession(s2, s2User, s2ID)

	req.Empty(s1)
	req.Empty(s2)

	req.NoError(ws.Send("msg", "msg", 0))
	req.Empty(s1)
	req.Empty(s2)

	req.NoError(ws.Send("msg", "msg1", s1User))
	req.Equal(s1.String(), `{"@type":"msg","@value":"msg1"}`)
	req.Equal(s2.String(), "")

	req.NoError(ws.Send("msg", "msg2", s2User))
	req.Equal(s1.String(), `{"@type":"msg","@value":"msg1"}`)
	req.Equal(s2.String(), `{"@type":"msg","@value":"msg2"}`)

	req.NoError(ws.Send("both", "msg3", s1User, s2User))
	req.Equal(s1.String(), `{"@type":"msg","@value":"msg1"}{"@type":"both","@value":"msg3"}`)
	req.Equal(s2.String(), `{"@type":"msg","@value":"msg2"}{"@type":"both","@value":"msg3"}`)

	req.NoError(ws.Send("all", "msg4"))
	req.Equal(s1.String(), `{"@type":"msg","@value":"msg1"}{"@type":"both","@value":"msg3"}{"@type":"all","@value":"msg4"}`)
	req.Equal(s2.String(), `{"@type":"msg","@value":"msg2"}{"@type":"both","@value":"msg3"}{"@type":"all","@value":"msg4"}`)
}
