package app

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	"time"
)

func (app *CortezaApp) initDAL(ctx context.Context, log *zap.Logger) (err error) {
	// no-op - if DAL is already initialized
	if dal.Initialized() {
		return
	}

	// Verify that primary store is connected
	// or return error
	if app.Store == nil {
		return fmt.Errorf("primary store not connected")
	}

	var (
		conn *types.DalConnection
		cw   *dal.ConnectionWrap
	)

	// load/create primary connection
	if conn, err = provisionPrimaryDalConnection(ctx, app.Store); err != nil {
		return fmt.Errorf("could not provision primary connection: %w", err)
	}

	// Convert connection to dal.ConnectionWrap
	if cw, err = service.MakeDalConnection(conn, app.Store.ToDalConn()); err != nil {
		return fmt.Errorf("could not convert connection: %w", err)
	}

	// Init DAL and prepare default connection
	dal.SetGlobal(dal.New(log.Named("dal"), app.Opt.Environment.IsDevelopment()))
	if err = dal.Service().ReplaceConnection(ctx, cw, true); err != nil {
		return fmt.Errorf("could not set primary connection: %w", err)
	}

	return
}

// Creates entry for primary connection in the store
func provisionPrimaryDalConnection(ctx context.Context, s store.DalConnections) (conn *types.DalConnection, err error) {
	conn, err = store.LookupDalConnectionByHandle(ctx, s, types.DalPrimaryConnectionHandle)
	if err != nil && err != store.ErrNotFound {
		return
	}

	// Already exists
	if conn != nil {
		return
	}

	conn = &types.DalConnection{
		// Using id.Next since we dropped "special" ids a while ago.
		// If needed, use the handle
		ID:     id.Next(),
		Handle: types.DalPrimaryConnectionHandle,
		Type:   types.DalPrimaryConnectionResourceType,

		Meta: types.ConnectionMeta{
			Name: "Primary Database",
		},

		Config: types.ConnectionConfig{
			DAL: &types.ConnectionConfigDAL{
				ModelIdent: "compose_record",
				Operations: dal.FullOperations(),
			},
		},

		CreatedAt: time.Now(),
	}

	return conn, store.CreateDalConnection(ctx, s, conn)
}
