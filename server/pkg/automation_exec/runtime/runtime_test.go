package runtime

// import (
// 	"context"
// 	"testing"

// 	"github.com/crusttech/human/server/pkg/automation_exec/types"
// 	"github.com/crusttech/human/server/pkg/id"
// )

// // ---- ExecutionGate ----

// type NoopGate struct{}

// func (g *NoopGate) Request(ops int) (<-chan struct{}, error) {
// 	ch := make(chan struct{}, 1)
// 	ch <- struct{}{}
// 	return ch, nil
// }

// // ---- StateReporter ----

// type NoopReporter struct{}

// func (r *NoopReporter) ReportStepStart(stepID id.ID) error {
// 	return nil
// }

// func (r *NoopReporter) ReportStepComplete(stepID id.ID, result StepResult) error {
// 	return nil
// }

// func (r *NoopReporter) ReportStepFailed(stepID id.ID, err error) error {
// 	return nil
// }

// func (r *NoopReporter) ReportExecutionComplete() error {
// 	return nil
// }

// func (r *NoopReporter) ReportExecutionFailed(err error) error {
// 	return nil
// }

// // ---- helpers ----

// func BackgroundContext() context.Context {
// 	return context.Background()
// }

// func BenchmarkThing(b *testing.B) {
// 	ctx := context.Background()

// 	exc := DummyExecutable()
// 	g := &NoopGate{}
// 	rep := &NoopReporter{}

// 	scp := map[string]any{}

// 	for n := 0; n < b.N; n++ {
// 		svc := Runtime(
// 			exc,
// 			g,
// 			rep,
// 		)

// 		err := svc.Start(ctx, scp)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// type noopHandler struct{}

// func (h noopHandler) Execute(ctx context.Context, scope map[string]any) (map[string]any, error) {
// 	return map[string]any{}, nil
// }

// func DummyExecutable() types.Executable {
// 	steps := make([]types.Step, 0, 10)
// 	relations := make([]types.Relationship, 0, 9)

// 	for i := 0; i < 10; i++ {
// 		stepID := id.MustNumID(uint64(i + 1))

// 		steps = append(steps, types.Step{
// 			ID:      stepID,
// 			Kind:    "noop",
// 			Config:  map[string]any{},
// 			Handler: noopHandler{},
// 		})

// 		if i > 0 {
// 			relations = append(relations, types.Relationship{
// 				From: steps[i-1].ID,
// 				To:   stepID,
// 			})
// 		}
// 	}

// 	return types.Executable{
// 		ID:            id.MustNumID(1),
// 		Revision:      1,
// 		Label:         "dummy-executable",
// 		Description:   "dummy executable with 10 noop steps",
// 		Limits:        types.ExecutableLimits{},
// 		Steps:         steps,
// 		Relationships: relations,
// 		Metadata:      map[string]any{},
// 	}
// }
