package auth

import (
	"sort"

	"github.com/cortezaproject/corteza/server/pkg/slice"
)

// ApplyRoleSecurity takes role security params (set of permitted, prohibited and forced roles)
// and applies these rules to the set of given roles
//
// Filtered set of roles is returned
//
// String slices are used intentionally, because of the data source used
func ApplyRoleSecurity(permitted, prohibited, forced []uint64, rr ...uint64) (out []uint64) {
	var (
		mPermitted  = slice.ToUint64BoolMap(permitted)
		mProhibited = slice.ToUint64BoolMap(prohibited)
		mForced     = slice.ToUint64BoolMap(forced)
	)

	// iterate over user's roles and just append them (obeying allow&deny rules)
	// to list of mForced roles
	for _, r := range rr {
		if (len(mPermitted) == 0 || mPermitted[r]) && !mProhibited[r] {
			mForced[r] = true
		}
	}

	out = make([]uint64, 0, len(mForced))
	for forcedRoleID := range mForced {
		out = append(out, forcedRoleID)
	}

	// for stable output
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})

	return
}
