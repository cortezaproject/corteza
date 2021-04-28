package websocket

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/options"
	gWebsocket "github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// WebsocketServer provide websocket server for testing
func WebsocketTestServer(t *testing.T) (*httptest.Server, *gWebsocket.Conn) {
	// Create test server with the websocket handler.
	s := httptest.NewServer(http.HandlerFunc(wsOpen))

	// Convert http://.. to ws://..
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	conn, _, err := gWebsocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("WebsocketServer() error while creating connection = %v", err)
	}

	return s, conn
}

// wsOpen opens websocket connection
func wsOpen(w http.ResponseWriter, r *http.Request) {
	var gUpgrader = gWebsocket.Upgrader{}
	c, err := gUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer func(c *gWebsocket.Conn) {
		_ = c.Close()
	}(c)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestSendingMessageToUser(t *testing.T) {
	tests := []struct {
		name      string
		kind      string
		payload   interface{}
		expectedP string
	}{
		{
			name: "send json",
			kind: "Json",
			payload: struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			}{Title: "Websocket", Description: "Testing connection.."},
			expectedP: `{"status":"ok","data":{"title":"Websocket","description":"Testing connection.."}}`,
		},
		{
			name:      "send text",
			kind:      "Text",
			payload:   "testing connectivity",
			expectedP: "testing connectivity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logger *zap.Logger
			var config options.WebsocketOpt
			ws := Websocket(logger, config)
			s, conn := WebsocketTestServer(t)
			defer s.Close()
			defer func(ws *gWebsocket.Conn) {
				err := ws.Close()
				if err != nil {
					t.Fatalf("TestSendingMessageToUser() error closing connection = %v", err)
				}
			}(conn)

			// Open a session using ws connection
			wsSession := Session(context.Background(), ws.logger, ws.config, conn)

			var messageType int
			var data []byte
			switch tt.kind {
			case "Text":
				messageType = gWebsocket.TextMessage
				data = []byte(fmt.Sprintf("%v", tt.payload))
			case "Json":
				res := Response("ok", tt.payload)
				messageType = gWebsocket.BinaryMessage
				var err error
				data, err = res.Marshal()
				if err != nil {
					t.Fatalf("TestSendingMessageToUser() error while marshaling payload = %v", err)
				}
			}

			// Send message to server, read response and check to see if it's what we expect.
			if err := wsSession.conn.WriteMessage(messageType, data); err != nil {
				t.Fatalf("TestSendingMessageToUser() error while sending message = %v", err)
			}

			_, p, err := wsSession.conn.ReadMessage()
			if err != nil {
				t.Fatalf("TestSendingMessageToUser() error while reading message =%v", err)
			}

			if string(p) != tt.expectedP {
				t.Fatalf("TestSendingMessageToUser() gotP = %v, want = %v", string(p), tt.expectedP)
			}
		})
	}
}
