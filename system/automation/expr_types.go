package automation

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/system/types"
)

func CastToUser(val interface{}) (out *types.User, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.User{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToUser(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.User:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToRole(val interface{}) (out *types.Role, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.Role{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToRole(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.Role:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}
