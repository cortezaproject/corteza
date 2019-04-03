package db

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

func TryToConnect(ctx context.Context, name, dsn, profiler string) (db *factory.DB, err error) {
	factory.Database.Add(name, dsn)
	var try = 0

	timeout := time.After(time.Minute)
	delay := 5 * time.Second

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("cancelled")
		case <-timeout:
			return nil, errors.New("connection timeout")
		case <-time.After(delay):
			try++
			db, err = factory.Database.Get(name)
			if err != nil {
				log.Printf("Failed to connect, try %d, error: %v, retry in %.0fs", try, err, delay.Seconds())
				continue
			}
		}
		break
	}

	// @todo: profiling as an external service?
	switch profiler {
	case "stdout":
		db.Profiler = &factory.Database.ProfilerStdout
	default:
		log.Println("No database query profiler selected")
	}

	return
}
