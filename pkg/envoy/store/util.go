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

// Taken from compose/service/values/sanitizer.go
func toTime(v string) *time.Time {
	ff := []string{
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC850,
		time.RFC822Z,
		time.RFC822,
		time.RubyDate,
		time.UnixDate,
		time.ANSIC,
		"2006/_1/_2 15:04:05",
		"2006/_1/_2 15:04",
	}

	for _, f := range ff {
		parsed, err := time.Parse(f, v)
		if err == nil {
			return &parsed
		}
	}

	return nil
}

func resolveUserRefs(ctx context.Context, s store.Storer, pr []resource.Interface, refs resource.RefSet, dst map[string]uint64) (err error) {
	for _, uRef := range refs {
		u := findUserR(ctx, pr, uRef.Identifiers)
		if u == nil {
			u, err = findUserS(ctx, s, makeGenericFilter(uRef.Identifiers))
			if err != nil {
				return err
			}
		}
		if u == nil {
			return userErrUnresolved(uRef.Identifiers)
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
