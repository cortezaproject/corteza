package websocket

import (
	"net/http"

	"context"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/outgoing"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"
)

type (
	Websocket struct {
		svc struct {
			userFinder wsUserFinder
		}
	}

	wsUserFinder interface {
		FindByID(ctx context.Context, userID uint64) (*types.User, error)
	}
)

func (Websocket) New(svcUser wsUserFinder) *Websocket {
	ws := &Websocket{}
	ws.svc.userFinder = svcUser
	return ws
}

// Handles websocket requests from peers
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Allow connections from any Origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (ws Websocket) Open(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Disallow all unauthorized!
	identity := auth.GetIdentityFromContext(ctx)
	if !identity.Valid() {
		resputil.JSON(w, errors.New("Unauthorized"))
		return
	}

	// @todo validate user (ws.svc.userFinder) here...

	conn, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		resputil.JSON(w, errors.Wrap(err, "ws: need a websocket handshake"))
		return
	} else if err != nil {
		resputil.JSON(w, errors.Wrap(err, "ws: failed to upgrade connection"))
		return
	}

	session := store.Save(Session{}.New(r.Context(), conn))

	session.sendToAll(&outgoing.Connected{UserID: uint64toa(identity.GetID())})

	if err := session.Handle(); err != nil {
		// @todo: log error, because at this point we can't really write it to w
	}

	session.sendToAll(&outgoing.Disconnected{UserID: uint64toa(identity.GetID())})
}
