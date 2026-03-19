package registry

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/id"
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
