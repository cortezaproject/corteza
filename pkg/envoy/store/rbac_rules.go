package store

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/envoy/node"
	"github.com/cortezaproject/corteza-server/store"
	"time"
)

func storeRbacRules(ctx context.Context, s store.Storer, n *node.Role) error {
	var (
		res = n.Res
	)

	res.ID = nextID()
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return store.CreateRole(ctx, s, n.Res)
}
