package ledger

import (
	"context"

	"github.com/crusttech/human/server/pkg/cli"
	"github.com/crusttech/human/server/pkg/id"
	"go.uber.org/zap"
)

func init() {
	id.Init(cli.Context())
}

func newLedger() *ledger { return Ledger(zap.NewNop()) }

func nextID() id.ID { return id.MustNumID(id.Next()) }

var ctx = context.Background()
