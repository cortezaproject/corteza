package monitor

import (
	"context"
	"expvar"
	"runtime"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/options"
)

type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
}

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
	if opt.Interval > 0 {
		go NewMonitor(int(opt.Interval / time.Second))
		log.Debug("watcher initialized")
	}
}

func NewMonitor(duration int) {
	var (
		m          = Monitor{}
		rtm        runtime.MemStats
		goroutines = expvar.NewInt("num_goroutine")
	)
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()
		goroutines.Set(int64(m.NumGoroutine))

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		log.With(
			zap.Uint64("alloc", m.Alloc),
			zap.Uint64("totalAlloc", m.TotalAlloc),
			zap.Uint64("sys", m.Sys),
			zap.Uint64("mallocs", m.Mallocs),
			zap.Uint64("frees", m.Frees),
			zap.Uint64("liveObjects", m.LiveObjects),
			zap.Uint64("pauseTotalNs", m.PauseTotalNs),
			zap.Uint32("numGC", m.NumGC),
			zap.Int("numGoRoutines", m.NumGoroutine),
		).
			Info("tick")
	}
}
