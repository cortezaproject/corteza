package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Monitor struct {
	runtime.MemStats
	NumGoroutine int
}

func (m Monitor) MarshalJSON() ([]byte, error) {
	result := struct {
		Alloc,
		TotalAlloc,
		Sys,
		Mallocs,
		Frees uint64

		NumGC        uint32
		NumGoroutine int
	}{
		m.MemStats.Alloc,
		m.MemStats.TotalAlloc,
		m.MemStats.Sys,
		m.MemStats.Mallocs,
		m.MemStats.Frees,

		m.MemStats.NumGC,
		m.NumGoroutine,
	}
	return json.Marshal(result)
}

func NewMonitor(duration int) {
	var m Monitor
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)
		runtime.ReadMemStats(&m.MemStats)
		m.NumGoroutine = runtime.NumGoroutine()
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
	}
}
