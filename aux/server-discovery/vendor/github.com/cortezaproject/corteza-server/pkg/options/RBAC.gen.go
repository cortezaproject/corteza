package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/RBAC.yaml

type (
	RBACOpt struct {
		Log                bool   `env:"RBAC_LOG"`
		ServiceUser        string `env:"RBAC_SERVICE_USER"`
		BypassRoles        string `env:"RBAC_BYPASS_ROLES"`
		AuthenticatedRoles string `env:"RBAC_AUTHENTICATED_ROLES"`
		AnonymousRoles     string `env:"RBAC_ANONYMOUS_ROLES"`
	}
)

// RBAC initializes and returns a RBACOpt with default values
func RBAC() (o *RBACOpt) {
	o = &RBACOpt{
		BypassRoles:        "super-admin",
		AuthenticatedRoles: "authenticated",
		AnonymousRoles:     "anonymous",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *RBAC) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
