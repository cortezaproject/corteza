package permissions

// General permission stuff, types, constants

type (
	Operation string
	Access    int

	// CheckAccessFunc function.
	CheckAccessFunc func() Access

	Whitelist struct {
		// Index is used for fast lookups
		index map[Resource]map[Operation]bool

		// we need this to maintain a stable order of res/ops
		rules RuleSet
	}

	whitelistFlatten struct {
		Resource  `json:"resource"`
		Operation `json:"operation"`
	}
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
	if wl.index == nil {
		wl.index = map[Resource]map[Operation]bool{}
	}

	wl.index[r] = map[Operation]bool{}

	for _, o := range oo {
		wl.index[r][o] = true
		wl.rules = append(wl.rules, &Rule{Resource: r, Operation: o})
	}
}

func (wl Whitelist) Check(rule *Rule) bool {
	if rule == nil {
		return false
	}

	res := rule.Resource.TrimID()

	if _, ok := wl.index[res]; !ok {
		return false
	}

	return wl.index[res][rule.Operation]
}

// Flatten casts list of operations for each resource from map to slice and creates more output friendly format
func (wl Whitelist) Flatten() []whitelistFlatten {
	var (
		wlf = []whitelistFlatten{}
	)
	for _, r := range wl.rules {
		wlf = append(wlf, whitelistFlatten{r.Resource, r.Operation})
	}

	return wlf
}
