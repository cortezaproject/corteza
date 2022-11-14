package auth

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	identity struct {
		id       uint64
		memberOf []uint64
	}
)

// Anonymous constructs and returns new anonymous identity with system anonymous roles
func Anonymous() *identity {
	return &identity{
		// memberOf: AnonymousRoles().IDs(),
	}
}

// Authenticated constructs and returns new authenticated identity with assigned roles + system authenticated roles
func Authenticated(id uint64, rr ...uint64) *identity {
	return &identity{
		id: id,
		// memberOf: append(rr, AuthenticatedRoles().IDs()...),
	}
}

func (i identity) Identity() uint64 {
	return i.id
}

func (i identity) Roles() []uint64 {
	return i.memberOf
}

func (i identity) Valid() bool {
	return i.id > 0
}

func (i identity) String() string {
	return fmt.Sprintf("%d", i.Identity())
}

func ExtractFromSubClaim(sub string) (userID uint64, rr []uint64) {
	parts := strings.Split(sub, " ")
	rr = make([]uint64, len(parts)-1)
	for p := range parts {
		id, _ := strconv.ParseUint(parts[p], 10, 64)
		if p == 0 {
			userID = id
		} else {
			rr[p-1] = id
		}
	}

	return
}
