package provision

import "time"

func releaseDate(y int, m time.Month) time.Time {
	return time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC)
}
