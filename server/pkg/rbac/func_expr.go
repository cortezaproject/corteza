package rbac

import (
	"reflect"

	"github.com/PaesslerAG/gval"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/spf13/cast"
)

var (
	isAboveChecker = func(owner id.ID, user id.ID, paths ...string) bool {
		return gRBAC.orgTree.IsAbove(owner, user, paths...)
	}
)

func AllFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("isDescendantOf", isDescendantOf),

		gval.Function("isDescendantOfW", isDescendantOfW),
		gval.Function("isDescendantOfC", isDescendantOfC),
		gval.Function("isDescendantOfR", isDescendantOfR),
		gval.Function("isDescendantOfU", isDescendantOfU),
		gval.Function("isDescendantOfD", isDescendantOfD),
	}
}

func isDescendantOf(userID any, resourceOwner any, paths ...string) bool {
	if gRBAC == nil {
		return false
	}

	owners := toUint64Slice(resourceOwner)
	if len(owners) == 0 {
		return false
	}

	user := id.MustNumID(cast.ToUint64(cast.ToString(userID)))
	if user.IsZero() {
		return false
	}

	for _, owner := range owners {
		if isAboveChecker(id.MustNumID(owner), user, paths...) {
			return true
		}
	}

	return false
}

func isDescendantOfR(userID any, resourceOwner any) bool {
	return isDescendantOf(userID, resourceOwner, "read")
}

func isDescendantOfC(userID any, resourceOwner any) bool {
	return isDescendantOf(userID, resourceOwner, "create")
}

func isDescendantOfW(userID any, resourceOwner any) bool {
	return isDescendantOf(userID, resourceOwner, "create", "update", "delete")
}

func isDescendantOfU(userID any, resourceOwner any) bool {
	return isDescendantOf(userID, resourceOwner, "update")
}

func isDescendantOfD(userID any, resourceOwner any) bool {
	return isDescendantOf(userID, resourceOwner, "delete")
}

func toUint64Slice(v any) []uint64 {
	switch x := v.(type) {
	case uint64:
		return []uint64{x}
	case []uint64:
		return x
	case [1]uint64:
		return x[:]
	case [2]uint64:
		return x[:]
	case [3]uint64:
		return x[:]
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Array && rv.Type().Elem().Kind() == reflect.Uint64 {
			out := make([]uint64, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				out[i] = rv.Index(i).Uint()
			}
			return out
		}
	}
	return nil
}
