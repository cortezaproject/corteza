package rbac_test

import (
	"github.com/crusttech/crust/rbac"
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
		t.Errorf(format, args...)
	}
	return ok
}
