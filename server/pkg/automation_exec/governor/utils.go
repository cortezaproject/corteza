package governor

import "time"

func newWindow(max int, win time.Duration, now time.Time) windowCounter {
	if max <= 0 {
		return windowCounter{}
	}
	if win <= 0 {
		return windowCounter{max: max}
	}
	return windowCounter{
		max:   max,
		win:   win,
		reset: now.Add(win),
	}
}
