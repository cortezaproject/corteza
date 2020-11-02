package tmp

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func encodeComposeRecordSet(ctx context.Context, s store.Storer, rec *resource.ComposeRecordSet, rm resMap) (uint64, error) {
	// @todo...
	return 0, nil
}
