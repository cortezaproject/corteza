package yaml

type (
	RbacRules map[string][]string
	Rbac      struct {
		Allow RbacRules
		Deny  RbacRules
	}
)
