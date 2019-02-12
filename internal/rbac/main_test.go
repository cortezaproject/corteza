package rbac_test

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"

	"github.com/crusttech/crust/internal/rbac"
)

var loaded bool

func getClient() (*rbac.Client, error) {
	if !loaded {
		godotenv.Load("../../.env")

		rbac.Flags()
		flag.Parse()
		loaded = true
	}
	rbac.Debug()
	client, err := rbac.New()
	if err != nil {
		return nil, err
	}
	client.Debug("debug")
	return client, nil
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
