package rbac

import (
	"testing"
)

func TestFlags(t *testing.T) {
	c := configuration{}
	mustFail(t, c.validate())
	c.auth = "a"
	mustFail(t, c.validate())
	c.tenant = "a"
	mustFail(t, c.validate())
	c.baseURL = "a"
	must(t, c.validate())
}

/* imported below from main_test.go because of different package namespace */

func assert(t *testing.T, ok bool, format string, args ...interface{}) bool {
	if !ok {
		t.Fatalf(format, args...)
	}
	return ok
}

func must(t *testing.T, err error, message ...string) {
	if len(message) > 0 {
		assert(t, err == nil, message[0]+": %+v", err)
		return
	}
	assert(t, err == nil, "Error: %+v", err)
}

func mustFail(t *testing.T, err error) {
	assert(t, err != nil, "Expected error, got nil")
}
