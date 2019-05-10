package permissions

// General permission stuff, types, constants

type (
	Operation string
	Access    int

	// CheckAccessFunc function.
	CheckAccessFunc func() Access

	Whitelist map[Resource]map[Operation]bool
)

const (
	// Hardcoded Role ID for everyone
	EveryoneRoleID = 1

	// Hardcoded ID for Admin role
	AdminRoleID = 2
)

func (a Access) String() string {
	switch a {
	case Allow:
		return "allow"
	case Deny:
		return "deny"
	default:
		return "inherit"
	}
}

func (a *Access) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "allow":
		*a = Allow
	case "deny":
		*a = Deny
	default:
		*a = Inherit
	}
	return nil
}

func (a Access) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.String() + `"`), nil
}

func Allowed() Access {
	return Allow
}

func Denied() Access {
	return Deny
}

func (wl *Whitelist) Set(r Resource, oo ...Operation) {
	(*wl)[r] = map[Operation]bool{}

	for _, o := range oo {
		(*wl)[r][o] = true
	}
}

func (wl Whitelist) Check(rule *Rule) bool {
	if rule == nil {
		return false
	}

	res := rule.Resource.TrimID()

	if _, ok := wl[res]; !ok {
		return false
	}

	return wl[res][rule.Operation]
}
