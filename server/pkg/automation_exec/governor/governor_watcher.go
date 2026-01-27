package governor

import (
	"context"
	"time"
)

func (g *governor) watch(ctx context.Context) {
	for {
		g.mu.Lock()
		next := g.nextResetLocked()
		g.mu.Unlock()

		if next.IsZero() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(g.config.WatcherFallbackInterval):
				continue
			}
		}

		wait := time.Until(next)
		if wait < 0 {
			wait = 0
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
			now := g.now()
			g.mu.Lock()
			for execID, ep := range g.exec {
				g.refreshLocked(now, ep, execID)
			}
			g.mu.Unlock()
		}
	}
}

func (g *governor) nextResetLocked() time.Time {
	var next time.Time

	consider := func(t time.Time) {
		if t.IsZero() {
			return
		}
		if next.IsZero() || t.Before(next) {
			next = t
		}
	}

	consider(g.global.budget.reset)
	consider(g.global.rate.reset)

	for _, ep := range g.exec {
		consider(ep.budget.reset)
		consider(ep.rate.reset)
	}

	return next
}
