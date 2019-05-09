package permissions

// General permission stuff, types, constants

type (
	Operation string
	Access    int

	// CheckAccessFunc function.
	CheckAccessFunc func() Access
)

const EveryoneRoleID = 1

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
