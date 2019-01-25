package types

import (
	"github.com/crusttech/crust/internal/rules"
)

type (
	ResourceProvider interface {
		Scope() string
		Resource() string
		Operation(name string) string
		Permissions() []rules.OperationGroup
	}
)

// These entities create resources in RBAC
var _ ResourceProvider = &Organisation{}
var _ ResourceProvider = &Team{}
var _ ResourceProvider = &Channel{}
