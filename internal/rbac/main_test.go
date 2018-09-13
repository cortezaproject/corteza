package rbac_test

import (
	"github.com/crusttech/crust/internal/rbac"
	"github.com/namsral/flag"
	"testing"
)

var loaded bool

func getClient() (*rbac.Client, error) {
	if !loaded {
		rbac.Flags()
		flag.Parse()
		loaded = true
	}
	return rbac.New()
}

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
