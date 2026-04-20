package registry

import (
	"context"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/cli"
	"github.com/crusttech/human/server/pkg/id"
	"go.uber.org/zap"
)

func init() {
	id.Init(cli.Context())
}

var ctx = context.Background()

type mockUsageChecker struct {
	inUse bool
	err   error
}

func (m *mockUsageChecker) IsExecutableInUse(_ context.Context, _ id.ID, _ int) (bool, error) {
	return m.inUse, m.err
}

func newRegistry(checker usageChecker) *registry {
	return Registry(zap.NewNop(), checker)
}

func nextID() id.ID { return id.MustNumID(id.Next()) }

func makeExec(execID id.ID, revision int) types.Executable {
	return types.Executable{ID: execID, Revision: revision}
}
