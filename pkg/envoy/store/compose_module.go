package store

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/envoy/node"
	"github.com/cortezaproject/corteza-server/store"
	"time"
)

func storeComposeModule(ctx context.Context, s store.Storer, n *node.ComposeModule) error {
	var (
		res = n.Res
	)

	res.ID = nextID()
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return store.CreateComposeModule(ctx, s, n.Res)
}
