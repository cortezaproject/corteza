package db

import (
	"log"
	"time"

	"github.com/titpetric/factory"
)

func TryToConnect(dsn, profiler string) (db *factory.DB, err error) {
	factory.Database.Add("default", dsn)
	var try = 0

connect:
	try++
	db, err = factory.Database.Get()

	if err != nil {
		delay := time.Second * 5
		log.Printf("Failed to connect, try %d, error: %v, retry in %.0fs", try, err, delay.Seconds())
		time.Sleep(delay)
		goto connect
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
