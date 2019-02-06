package types

import (
	"github.com/crusttech/crust/internal/rules"
)

type (
	ResourceProvider interface {
		Resource() rules.Resource
	}
)

// These entities create resources in RBAC
var _ ResourceProvider = &Organisation{}
var _ ResourceProvider = &Team{}
var _ ResourceProvider = &Channel{}
