package rules

type Access int

const (
	Allow   Access = 2
	Deny           = 1
	Inherit        = 0
)

type Rule struct {
	RoleID    uint64   `json:"roleID,string" db:"rel_role"`
	Resource  Resource `json:"resource" db:"resource"`
	Operation string   `json:"operation" db:"operation"`
	Value     Access   `json:"value,string" db:"value"`
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
	var str string

	switch a {
	case Allow:
		str = "allow"
	case Deny:
		str = "deny"
	default:
		str = "inherit"
	}

	return []byte(`"` + str + `"`), nil
}
