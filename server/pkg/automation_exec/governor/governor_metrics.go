package governor

import "sync/atomic"

type (
	Metrics struct {
		TotalRequests         atomic.Int64
		GrantedImmediately    atomic.Int64
		BlockedByPause        atomic.Int64
		BlockedByExecBudget   atomic.Int64
		BlockedByExecRate     atomic.Int64
		BlockedByGlobalBudget atomic.Int64
		BlockedByGlobalRate   atomic.Int64
	}

	MetricsSnapshot struct {
		TotalRequests         int64
		GrantedImmediately    int64
		BlockedByPause        int64
		BlockedByExecBudget   int64
		BlockedByExecRate     int64
		BlockedByGlobalBudget int64
		BlockedByGlobalRate   int64
	}
)

func (g *governor) GetMetrics() MetricsSnapshot {
	return MetricsSnapshot{
		TotalRequests:         g.metrics.TotalRequests.Load(),
		GrantedImmediately:    g.metrics.GrantedImmediately.Load(),
		BlockedByPause:        g.metrics.BlockedByPause.Load(),
		BlockedByExecBudget:   g.metrics.BlockedByExecBudget.Load(),
		BlockedByExecRate:     g.metrics.BlockedByExecRate.Load(),
		BlockedByGlobalBudget: g.metrics.BlockedByGlobalBudget.Load(),
		BlockedByGlobalRate:   g.metrics.BlockedByGlobalRate.Load(),
	}
}

func (g *governor) ResetMetrics() {
	g.metrics = Metrics{}
}
