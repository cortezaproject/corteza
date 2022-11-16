package monitor

import (
	"context"
	"runtime"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza/server/pkg/options"
)

var (
	// Holds options for monitor
	opt options.MonitorOpt

	log *zap.Logger
)

func Setup(logger *zap.Logger, o options.MonitorOpt) {
	log = logger.Named("monitor")
	opt = o
}

func Watcher(ctx context.Context) {
	if opt.Interval == 0 {
		return
	}

	if opt.Interval < time.Second {
		log.Warn("monitoring interval less than 1 second, disabling")
		return
	}

	go NewMonitor(ctx, opt.Interval)
	log.Debug("watcher initialized", zap.Duration("interval", opt.Interval))
}

func NewMonitor(ctx context.Context, interval time.Duration) {
	var (
		ticker = time.NewTicker(interval)
	)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Read full mem stats
			m := new(runtime.MemStats)
			runtime.ReadMemStats(m)

			log.With(
				zap.Uint64("alloc", m.HeapAlloc),
				zap.Uint64("totalAlloc", m.TotalAlloc),
				zap.Uint64("sys", m.Sys),
				zap.Uint64("mallocs", m.Mallocs),
				zap.Uint64("frees", m.Frees),

				zap.Uint64("liveObjects", m.Mallocs-m.Frees),
				zap.Uint64("pauseTotalNs", m.PauseTotalNs),
				zap.Uint32("numGC", m.NumGC),
				zap.Int("numGoRoutines", runtime.NumGoroutine()),
			).Info("tick")

		case <-ctx.Done():
			return
		}

	}
}
