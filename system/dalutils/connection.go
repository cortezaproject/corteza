package dalutils

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	connectionCreator interface {
		CreateConnection(ctx context.Context, connectionID uint64, cp dal.ConnectionParams, cm dal.ConnectionMeta, capabilities ...capabilities.Capability) (err error)
	}

	connectionDeleter interface {
		DeleteConnection(ctx context.Context, connectionID uint64) (err error)
	}

	connectionUpdater interface {
		UpdateConnection(ctx context.Context, connectionID uint64, cp dal.ConnectionParams, cm dal.ConnectionMeta, capabilities ...capabilities.Capability) (err error)
	}

	connectionDeleteCreator interface {
		connectionDeleter
		connectionCreator
	}
)

func DalConnectionReload(ctx context.Context, s store.Storer, dc connectionDeleteCreator) (err error) {
	// Get all available connections
	cc, _, err := store.SearchDalConnections(ctx, s, types.DalConnectionFilter{
		Type: types.DalConnectionResourceType,
	})
	if err != nil {
		return
	}

	for _, c := range cc {
		var cm dal.ConnectionMeta
		cm, err = ConnectionMeta(ctx, c)
		if err != nil {
			return
		}
		if err = dc.CreateConnection(ctx, c.ID, c.Config.Connection, cm, c.ActiveCapabilities()...); err != nil {
			return
		}
	}

	return
}

func DalConnectionCreate(ctx context.Context, c connectionCreator, connections ...*types.DalConnection) (err error) {
	var cm dal.ConnectionMeta
	for _, connection := range connections {
		cm, err = ConnectionMeta(ctx, connection)
		if err != nil {
			return
		}

		if err = c.CreateConnection(ctx, connection.ID, connection.Config.Connection, cm, connection.ActiveCapabilities()...); err != nil {
			return err
		}
	}

	return
}

func DalConnectionUpdate(ctx context.Context, u connectionUpdater, connections ...*types.DalConnection) (err error) {
	var cm dal.ConnectionMeta
	for _, connection := range connections {
		cm, err = ConnectionMeta(ctx, connection)
		if err != nil {
			return
		}

		if err = u.UpdateConnection(ctx, connection.ID, connection.Config.Connection, cm, connection.ActiveCapabilities()...); err != nil {
			return err
		}
	}

	return
}

func DalConnectionDelete(ctx context.Context, d connectionDeleter, connections ...*types.DalConnection) (err error) {
	for _, connection := range connections {
		if err = d.DeleteConnection(ctx, connection.ID); err != nil {
			return err
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utils

func ConnectionMeta(ctx context.Context, c *types.DalConnection) (cm dal.ConnectionMeta, err error) {
	// @todo we could probably utilize connection params more here
	cm = dal.ConnectionMeta{
		DefaultModelIdent:      c.Config.DefaultModelIdent,
		DefaultAttributeIdent:  c.Config.DefaultAttributeIdent,
		DefaultPartitionFormat: c.Config.DefaultPartitionFormat,
		SensitivityLevel:       c.SensitivityLevel,
		Label:                  c.Handle,
	}

	return
}
