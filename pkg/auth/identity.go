package auth

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	Identity struct {
		id       uint64
		memberOf []uint64
	}
)

const (
	superUserID uint64 = 10000000000000000
)

func NewIdentity(id uint64, rr ...uint64) *Identity {
	return &Identity{
		id:       id,
		memberOf: rr,
	}
}

func (i Identity) Identity() uint64 {
	return i.id
}

func (i Identity) Roles() []uint64 {
	return i.memberOf
}

func (i Identity) Valid() bool {
	return i.id > 0
}

func (i Identity) String() string {
	return fmt.Sprintf("%d", i.Identity())
}

func NewSuperUserIdentity() *Identity {
	return NewIdentity(superUserID)
}

func IsSuperUser(i Identifiable) bool {
	return i != nil && superUserID == i.Identity()
}

func ExtractUserIDFromSubClaim(sub string) uint64 {
	userID, _ := ExtractFromSubClaim(sub)
	return userID
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
