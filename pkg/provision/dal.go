package provision

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	DefaultComposeRecordTable    = "compose_record"
	DefaultComposeRecordValueCol = "values"
	DefaultPartitionFormat       = "compose_record"
)

// Injects primary connection
func defaultDalConnection(ctx context.Context, s store.DalConnections) (err error) {
	cc, err := store.LookupDalConnectionByHandle(ctx, s, types.DalPrimaryConnectionHandle)
	if err != nil && err != store.ErrNotFound {
		return
	}

	// Already exists
	if cc != nil {
		return
	}

	// Create it
	var conn = &types.DalConnection{
		// Using id.Next since we dropped "special" ids a while ago.
		// If needed, use the handle
		ID:     id.Next(),
		Name:   "Primary Database",
		Handle: types.DalPrimaryConnectionHandle,
		Type:   types.DalPrimaryConnectionResourceType,

		Config: types.ConnectionConfig{
			DefaultModelIdent:      DefaultComposeRecordTable,
			DefaultAttributeIdent:  DefaultComposeRecordValueCol,
			DefaultPartitionFormat: DefaultPartitionFormat,
		},
		Capabilities: types.ConnectionCapabilities{
			Supported: capabilities.FullCapabilities(),
		},
		CreatedAt: *now(),
		CreatedBy: auth.ServiceUser().ID,
	}

	return store.CreateDalConnection(ctx, s, conn)
}
