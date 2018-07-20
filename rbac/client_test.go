package rbac_test

import (
	"github.com/crusttech/crust/rbac"
	"github.com/namsral/flag"
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
