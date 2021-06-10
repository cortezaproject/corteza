package store

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	// genericFilter defines resource identifiers in a generic manner
	genericFilter struct {
		id          uint64
		identifiers []string
	}
)

var (
	// simple uint check.
	// we'll use the pkg/handle to check for handles.
	refy = regexp.MustCompile(`^[1-9](\d*)$`)

	// wrapper around NextID that will aid service testing
	NextID = func() uint64 {
		return id.Next()
	}

	// wrapper around time.Now() that will aid testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}

	// Global expression parser
	exprP = expr.Parser()
)

func makeGenericFilter(ii resource.Identifiers) (f genericFilter) {
	for i := range ii {
		if i == "" {
			continue
		}

		if refy.MatchString(i) && f.id <= 0 {
			id, err := strconv.ParseUint(i, 10, 64)
			if err != nil {
				continue
			}
			f.id = id
		} else {
			f.identifiers = append(f.identifiers, i)
		}
	}

	return f
}

// resolveUserstamps tries to resolve any missing userstamp using the store
//
// @todo rework based on the bellow userIndex
func resolveUserstamps(ctx context.Context, s store.Storer, rr []resource.Interface, us *resource.Userstamps) (*resource.Userstamps, error) {
	if us == nil {
		return nil, nil
	}

	fetch := func(us *resource.Userstamp) (*resource.Userstamp, error) {
		if us == nil || (us.UserID == 0 && us.Ref == "") {
			return nil, nil
		}

		// This one can be considered as valid
		if us.Ref != "" && us.UserID > 0 && us.U != nil {
			return us, nil
		}

		ii := resource.MakeIdentifiers()

		if us.UserID > 0 {
			ii = ii.Add(strconv.FormatUint(us.UserID, 10))
		}
		if us.Ref != "" {
			ii = ii.Add(us.Ref)
		}

		u, err := findUser(ctx, s, rr, ii)
		if err != nil {
			return nil, err
		}
		if u == nil {
			return nil, resource.UserErrUnresolved(ii)
		}

		return resource.MakeUserstamp(u), nil
	}

	var err error
	us.CreatedBy, err = fetch(us.CreatedBy)
	us.UpdatedBy, err = fetch(us.UpdatedBy)
	us.DeletedBy, err = fetch(us.DeletedBy)
	us.OwnedBy, err = fetch(us.OwnedBy)
	us.RunAs, err = fetch(us.RunAs)

	if err != nil {
		return nil, err
	}

	return us, nil
}

// syncUserStamps tries to resolve any missing userstamp using the userIndex
func syncUserStamps(us *resource.Userstamps, ux *userIndex) {
	if us.CreatedBy != nil && us.CreatedBy.UserID > 0 {
		us.CreatedBy.U = ux.users[us.CreatedBy.UserID]
	}
	if us.UpdatedBy != nil && us.UpdatedBy.UserID > 0 {
		us.UpdatedBy.U = ux.users[us.UpdatedBy.UserID]
	}
	if us.DeletedBy != nil && us.DeletedBy.UserID > 0 {
		us.DeletedBy.U = ux.users[us.DeletedBy.UserID]
	}
	if us.OwnedBy != nil && us.OwnedBy.UserID > 0 {
		us.OwnedBy.U = ux.users[us.OwnedBy.UserID]
	}
	if us.RunAs != nil && us.RunAs.UserID > 0 {
		us.RunAs.U = ux.users[us.RunAs.UserID]
	}
}

func mergeConfig(ec *EncoderConfig, rs *resource.EnvoyConfig) *EncoderConfig {
	// Nothing we can do
	if rs == nil {
		return ec
	}

	// Take resource config as base
	rr := &EncoderConfig{
		OnExisting: rs.OnExisting,
		SkipIf:     rs.SkipIf,

		// IgnoreStore is an encoder thing and should not be controlled from a resource
		IgnoreStore: ec.IgnoreStore,
	}

	// Default to store config
	rr.Defer = ec.Defer
	rr.DeferNok = ec.DeferNok
	rr.DeferOk = ec.DeferOk
	if rr.OnExisting == resource.Default {
		rr.OnExisting = ec.OnExisting
	}
	if rr.SkipIf == "" {
		rr.SkipIf = ec.SkipIf
	}

	return rr
}

func basicSkipEval(ctx context.Context, cfg *EncoderConfig, missing bool) (bool, error) {
	if cfg == nil {
		cfg = &EncoderConfig{}
	}

	if cfg.SkipIf != "" {
		evl, err := exprP.NewEvaluable(cfg.SkipIf)
		if err != nil {
			return false, err
		}

		// @todo expand this
		skip, err := evl.EvalBool(ctx, map[string]interface{}{
			"missing": missing,
		})
		if err != nil {
			return false, err
		}

		return skip, nil
	}

	// Don't skip by default
	return false, nil
}
