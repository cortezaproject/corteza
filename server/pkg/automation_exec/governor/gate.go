package governor

import "github.com/cortezaproject/corteza/server/pkg/id"

type (
	gates struct {
		globalPause      *gate
		globalBudgetHold *gate
		globalRateHold   *gate
		exec             map[id.ID]*execGates
	}

	gate struct {
		ch chan struct{}
	}
)

func newGate() *gate { return &gate{ch: make(chan struct{})} }
func (g *gate) open() {
	select {
	case <-g.ch:
	default:
		close(g.ch)
	}
}
