package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/sam/repository"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	Websocket struct {
		svc struct {
			user systemService.UserService
		}
		config *repository.Flags
	}
)

func (Websocket) New(config *repository.Flags) *Websocket {
	ws := &Websocket{
		config: config,
	}
	ws.svc.user = systemService.DefaultUser
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

	user, err := ws.svc.user.With(ctx).FindByID(identity.Identity())
	if err != nil {
		resputil.JSON(w, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		resputil.JSON(w, errors.Wrap(err, "ws: need a websocket handshake"))
		return
	} else if err != nil {
		resputil.JSON(w, errors.Wrap(err, "ws: failed to upgrade connection"))
		return
	}

	session := store.Save((&Session{}).New(ctx, ws.config, conn))
	session.user = user

	if err := session.Handle(); err != nil {
		log.Printf("Session handler returned an error: %v", err)
	}

}
