package types

import (
	"fmt"
	"runtime"
	"testing"
)

func assert(t *testing.T, ok bool, format string, args ...interface{}) bool {
	if !ok {
		_, file, line, _ := runtime.Caller(1)
		caller := fmt.Sprintf("\nAsserted at:%s:%d", file, line)

		t.Fatalf(format+caller, args...)
	}
	return ok
}
