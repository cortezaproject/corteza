package test

import (
	"fmt"
	"runtime"
	"testing"
)

func Assert(t *testing.T, ok bool, format string, args ...interface{}) bool {
	if !ok {
		_, file, line, _ := runtime.Caller(1)
		caller := fmt.Sprintf("\nAsserted at:%s:%d", file, line)

		t.Fatalf(format+caller, args...)
	}
	return ok
}

func NoError(t *testing.T, err error, format string, args ...interface{}) bool {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		caller := fmt.Sprintf("\nAsserted at:%s:%d", file, line)
		t.Fatalf(format+caller, append([]interface{}{err}, args...)...)
		return false
	}
	return true
}
