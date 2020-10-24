package store

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/envoy/node"
	"github.com/cortezaproject/corteza-server/store"
)

func storeComposeRecord(ctx context.Context, s store.Storer, n *node.ComposeRecordSet) error {
	//var (
	//	res = n.Res
	//)
	//
	//res.ID = nextID()
	//if res.CreatedAt.IsZero() {
	//	res.CreatedAt = time.Now()
	//}
	//
	//return store.CreateComposeRecord(ctx, s, n.Res)
	return nil
}
