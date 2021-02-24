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

	exprP = expr.Parser()
)

// makeGenericFilter is a helper to determine the base resource filter.
//
// It attempts to determine an identifier, handle, and name.
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

func resolveUserstamps(ctx context.Context, s store.Storer, rr []resource.Interface, us *resource.Userstamps) (*resource.Userstamps, error) {
	if us == nil {
		return nil, nil
	}

	fetch := func(us *resource.Userstamp) (*resource.Userstamp, error) {
		if us == nil {
			return nil, nil
		}
		ii := resource.MakeIdentifiers()

		if us.UserID > 0 {
			ii = ii.Add(strconv.FormatUint(us.UserID, 10))
		}
		if us.Ref != "" {
			ii = ii.Add(us.Ref)
		}
		if us.U != nil {
			if us.U.Handle != "" {
				ii = ii.Add(us.U.Handle)
			}
			if us.U.Username != "" {
				ii = ii.Add(us.U.Username)
			}
			if us.U.Email != "" {
				ii = ii.Add(us.U.Email)
			}
		}

		u, err := findUserRS(ctx, s, rr, ii)
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

	if err != nil {
		return nil, err
	}

	return us, nil
}

func resolveUserRefs(ctx context.Context, s store.Storer, pr []resource.Interface, refs resource.RefSet, dst map[string]uint64) (err error) {
	for _, uRef := range refs {
		u := resource.FindUser(pr, uRef.Identifiers)
		if u == nil {
			u, err = findUserS(ctx, s, makeGenericFilter(uRef.Identifiers))
			if err != nil {
				return err
			}
		}
		if u == nil {
			return resource.UserErrUnresolved(uRef.Identifiers)
		}

		// Unexisting users will have ID 0, but that's ok, as long as they exist
		for i := range uRef.Identifiers {
			dst[i] = u.ID
		}
	}
	return nil
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
