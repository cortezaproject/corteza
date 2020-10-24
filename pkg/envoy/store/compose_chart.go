package store

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/envoy/node"
	"github.com/cortezaproject/corteza-server/store"
	"time"
)

func storeComposeChart(ctx context.Context, s store.Storer, n *node.ComposeChart) error {
	var (
		res = n.Res
	)

	res.ID = nextID()
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return store.CreateComposeChart(ctx, s, n.Res)
}
