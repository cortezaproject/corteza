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

//type (
//	mockConn struct {
//		close            func() error
//		remoteAddr       func() net.Addr
//		writeMessage     func(messageType int, data []byte) error
//		setWriteDeadline func(t time.Time) error
//		readMessage      func() (messageType int, p []byte, err error)
//		setReadDeadline  func(t time.Time) error
//		setPongHandler   func(h func(appData string) error)
//	}
//)
//
//func MockConn() *mockConn {
//	return &mockConn{
//		close:            func() (err error) { return },
//		remoteAddr:       func() (addr net.Addr) { return &net.IPAddr{IP: net.IPv4(0, 0, 0, 0)} },
//		writeMessage:     func(messageType int, data []byte) (err error) { return },
//		setWriteDeadline: func(t time.Time) (err error) { return },
//		readMessage:      func() (messageType int, p []byte, err error) { return },
//		setReadDeadline:  func(t time.Time) (err error) { return },
//		setPongHandler:   func(h func(appData string) error) {},
//	}
//}
//
//func (c *mockConn) Close() error         { return c.close() }
//func (c *mockConn) RemoteAddr() net.Addr { return c.remoteAddr() }
//func (c *mockConn) WriteMessage(messageType int, data []byte) error {
//	return c.writeMessage(messageType, data)
//}
//func (c *mockConn) SetWriteDeadline(t time.Time) error                  { return c.setWriteDeadline(t) }
//func (c *mockConn) ReadMessage() (messageType int, p []byte, err error) { return c.readMessage() }
//func (c *mockConn) SetReadDeadline(t time.Time) error                   { return c.setReadDeadline(t) }
//func (c *mockConn) SetPongHandler(h func(appData string) error)         { c.setPongHandler(h) }

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
	req.Nil(s.Identity())

	req.EqualError(s.procRawMessage(mockResponse(nil)), "unauthorized: failed to parse token: EOF")
	req.Nil(s.Identity())

	token = []byte("one")
	req.NoError(s.procRawMessage(mockResponse(token)))
	req.NotNil(s.Identity())
	req.Equal(identity1.Identity(), s.Identity().Identity())

	req.EqualError(s.procRawMessage([]byte("{}")), "unknown message type ''")
	req.Equal(identity1.Identity(), s.Identity().Identity())

	token = []byte("one")
	req.NoError(s.procRawMessage(mockResponse(token)))
	req.NotNil(s.Identity())
	req.Equal(identity1.Identity(), s.Identity().Identity())

	token = []byte("two")
	req.EqualError(s.procRawMessage(mockResponse(token)), "unauthorized: identity does not match")
}
