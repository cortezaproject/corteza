package dev

import (
	"strconv"
	"sync"

	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type (
	// DebugMutex for verbose locking & unlocking
	// use:
	// dev.DebugMutex{Logger: logger.Default()}
	DebugMutex struct {
		sync.RWMutex
		c      atomic.Int64
		Logger *zap.Logger
	}
)

func (m *DebugMutex) log(r, u, done bool) {
	msg := "#" + strconv.FormatInt(m.c.Load(), 10)
	if u {
		msg = "UNLOCK " + msg
	} else {
		msg = "  LOCK " + msg
	}

	if r {
		msg = "R" + msg
	} else {
		msg = " " + msg
	}
	if done {
		msg += " DONE"
	}

	m.Logger.WithOptions(zap.AddCallerSkip(2)).Debug(msg)
}

func (m *DebugMutex) RLock() {
	m.c.Inc()
	m.log(true, false, false)
	m.RWMutex.RLock()
	m.log(true, false, true)
}

func (m *DebugMutex) RUnlock() {
	m.log(true, true, false)
	m.RWMutex.RUnlock()
	m.log(true, true, true)
}

func (m *DebugMutex) Lock() {
	m.c.Inc()
	m.log(true, false, false)
	m.RWMutex.Lock()
	m.log(true, false, true)
}

func (m *DebugMutex) Unlock() {
	m.log(true, true, false)
	m.RWMutex.Unlock()
	m.log(true, true, true)
}
