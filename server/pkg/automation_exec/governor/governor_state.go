package governor

import "time"

type (
	Budget struct {
		MaxOps int
		Window time.Duration // 0 = hard cap
	}

	RateLimit struct {
		MaxOps int
		Window time.Duration
	}

	windowCounter struct {
		max   int
		used  int
		win   time.Duration
		reset time.Time
	}

	globalPolicy struct {
		budget windowCounter
		rate   windowCounter
	}

	execPolicy struct {
		budget windowCounter
		rate   windowCounter
	}

	execGates struct {
		budget *gate
		rate   *gate
	}

	Config struct {
		WatcherFallbackInterval time.Duration
	}
)

func (g *governor) SetGlobalBudget(max int, window time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.global.budget = newWindow(max, window, g.now())
}

func (g *governor) SetGlobalRate(max int, window time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.global.rate = newWindow(max, window, g.now())
}
