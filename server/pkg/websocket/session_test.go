package websocket

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type (
	mockConn struct {
		close            func() error
		remoteAddr       func() net.Addr
		writeMessage     func(messageType int, data []byte) error
		setWriteDeadline func(t time.Time) error
		readMessage      func() (messageType int, p []byte, err error)
		setReadDeadline  func(t time.Time) error
		setPongHandler   func(h func(appData string) error)
	}
)

func MockConn() *mockConn {
	return &mockConn{
		close:            func() (err error) { return },
		remoteAddr:       func() (addr net.Addr) { return &net.IPAddr{IP: net.IPv4(0, 0, 0, 0)} },
		writeMessage:     func(messageType int, data []byte) (err error) { return },
		setWriteDeadline: func(t time.Time) (err error) { return },
		readMessage:      func() (messageType int, p []byte, err error) { return },
		setReadDeadline:  func(t time.Time) (err error) { return },
		setPongHandler:   func(h func(appData string) error) {},
	}
}

func (c *mockConn) Close() error         { return c.close() }
func (c *mockConn) RemoteAddr() net.Addr { return c.remoteAddr() }
func (c *mockConn) WriteMessage(messageType int, data []byte) error {
	return c.writeMessage(messageType, data)
}
func (c *mockConn) SetWriteDeadline(t time.Time) error                  { return c.setWriteDeadline(t) }
func (c *mockConn) ReadMessage() (messageType int, p []byte, err error) { return c.readMessage() }
func (c *mockConn) SetReadDeadline(t time.Time) error                   { return c.setReadDeadline(t) }
func (c *mockConn) SetPongHandler(h func(appData string) error)         { c.setPongHandler(h) }

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

func TestSession_writeDuringDisconnect(t *testing.T) {
	var (
		req = require.New(t)

		core, logs = observer.New(zap.DebugLevel)

		s = session{
			conn:   MockConn(),
			logger: zap.New(core),
			config: options.WebsocketOpt{},
			send:   make(chan []byte, 512),
			stop:   make(chan []byte, 1),
		}

		wg   sync.WaitGroup
		halt = make(chan struct{})
	)

	s.ctx, s.ctxCancel = context.WithCancel(context.Background())

	// hammer Write() from several goroutines while the session is torn
	// down underneath them
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-halt:
					return
				default:
					_, _ = s.Write([]byte("foo"))
				}
			}
		}()
	}

	s.disconnect()
	close(halt)
	wg.Wait()

	// disconnect must leave the channels open — closing them races with
	// Write() and panics with "send on closed channel". stop is never
	// written to, so a receive here only succeeds if it was closed.
	select {
	case _, ok := <-s.stop:
		req.True(ok, "stop channel must not be closed by disconnect")
	default:
	}

	req.Zero(logs.FilterMessageSnippet("recovering from websocket").Len(), "expected no recovered panics")
}

func TestSession_disconnected(t *testing.T) {
	var (
		req = require.New(t)

		core, logs = observer.New(zap.DebugLevel)

		s = session{
			conn:   MockConn(),
			logger: zap.New(core),
			config: options.WebsocketOpt{},
			send:   make(chan []byte, 1),
			stop:   make(chan []byte, 1),
		}
	)

	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	s.disconnect()

	// reading, writing and queueing on a disconnected session must report a
	// closed connection instead of panicking on a nil connection or on a
	// closed send channel
	raw, err := s.read()
	req.Nil(raw)
	req.ErrorIs(err, net.ErrClosed)

	req.ErrorIs(s.write(websocket.TextMessage, []byte("foo")), net.ErrClosed)

	n, err := s.Write([]byte("foo"))
	req.Zero(n)
	req.ErrorIs(err, net.ErrClosed)

	req.Zero(logs.FilterMessageSnippet("recovering from websocket").Len(), "expected no recovered panics")
}
