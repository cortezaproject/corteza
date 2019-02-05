package service

import (
	"fmt"
	"runtime"
	"testing"
)

type mockDB struct{}

func (mockDB) Transaction(callback func() error) error { return callback() }

func assert(t *testing.T, ok bool, format string, args ...interface{}) bool {
	if !ok {
		_, file, line, _ := runtime.Caller(1)
		caller := fmt.Sprintf("\nAsserted at:%s:%d", file, line)

		t.Fatalf(format+caller, args...)
	}
	return ok
}
