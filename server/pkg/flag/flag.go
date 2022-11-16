package flag

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/flag/types"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	FlaggedResource interface {
		GetFlags() []string
		SetFlags([]string)
		FlagResourceKind() string
		FlagResourceID() uint64
	}
)

// Search returns a slice of IDs corresponding to the filtered flags
func Search(ctx context.Context, s store.Storer, owner uint64, kind string, flags ...string) ([]uint64, error) {
	rr := make([]uint64, 0, 100)

	ff, _, err := store.SearchFlags(ctx, s, types.FlagFilter{
		Kind:    kind,
		OwnedBy: []uint64{0, owner},
		Name:    flags,
	})

	if err != nil {
		return nil, err
	}

	// Little helper to generate a map index for the fetched label resource
	mix := func(resID uint64) string {
		return fmt.Sprintf("%s:%d", kind, resID)
	}

	// Firstly get all of the flags for the given user.
	// Take note of inactive flags so we can filter them out of the global set
	out := make(map[string]bool)
	for _, f := range ff {
		if f.OwnedBy != 0 {
			if f.Active {
				rr = append(rr, f.ResourceID)
			} else {
				out[mix(f.ResourceID)] = true
			}
		}
	}

	// Go over global flags, exclude any ignored flags
	for _, f := range ff {
		if f.OwnedBy == 0 {
			if f.Active && !out[mix(f.ResourceID)] {
				rr = append(rr, f.ResourceID)
			}
		}
	}

	return rr, nil
}

// Create creates a new flag for the given resource
//
// If that flag for that owner for this resource already exists, it's skipped.
// Access control and any other validations should be performed by the caller.
func Create(ctx context.Context, s store.Storer, r FlaggedResource, ownedBy uint64, flag string) error {
	// Try to preload existing own flag
	own, err := store.LookupFlagByKindResourceIDOwnedByName(ctx, s, r.FlagResourceKind(), r.FlagResourceID(), ownedBy, flag)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	if own != nil && own.Active {
		return fmt.Errorf("flag %s for resource %s %d already exists", flag, r.FlagResourceKind(), r.FlagResourceID())
	}

	// If we have an inactive flag, mark it as active
	if own != nil && !own.Active {
		own.Active = true
		return store.UpdateFlag(ctx, s, own)
	}

	own = &types.Flag{
		Kind:       r.FlagResourceKind(),
		ResourceID: r.FlagResourceID(),
		OwnedBy:    ownedBy,
		Name:       flag,
		Active:     true,
	}

	return store.CreateFlag(ctx, s, own)
}

// Delete removes the flag from this resource
//
// Access control and any other validations should be performed by the caller.
//
// This operation has two outcomes:
// 	* if we are removing a flag defined for a specifc user, it is deleted
//  * if we are removing a flag defined globally (no owner), we create a new inactive flag
func Delete(ctx context.Context, s store.Storer, r FlaggedResource, ownedBy uint64, flag string) error {
	var (
		own    *types.Flag
		global *types.Flag
		err    error
	)

	// Try to find the global flag
	global, err = store.LookupFlagByKindResourceIDOwnedByName(ctx, s, r.FlagResourceKind(), r.FlagResourceID(), 0, flag)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	// Try to find own flag
	own, err = store.LookupFlagByKindResourceIDOwnedByName(ctx, s, r.FlagResourceKind(), r.FlagResourceID(), ownedBy, flag)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	if own == nil && global == nil {
		return fmt.Errorf("flag not found for resource %s %d", r.FlagResourceKind(), r.FlagResourceID())
	}

	// If we're deleting global flag, do it
	if ownedBy == 0 {
		if global == nil {
			return fmt.Errorf("global flag not found for %s %d", r.FlagResourceKind(), r.FlagResourceID())
		}

		return store.DeleteFlag(ctx, s, global)
	}

	// If we're deleting own flag and there is no global flag, delete own flag
	if own != nil && global == nil {
		return store.DeleteFlag(ctx, s, own)
	}

	// If we're deleting own flag and there is a global flag, mark own flag as inactive
	if own != nil && global != nil {
		own.Active = false
		return store.UpdateFlag(ctx, s, own)
	}

	// This can't happen, but just to be safe
	return fmt.Errorf("invalid flag removal state")
}

// Load updates the provided resources with storreed flags
//
// 1. All global flags for this resource are fetched
// 2. All user-specifc flags for this resource are fetched, overwriting global flags
func Load(ctx context.Context, s store.Storer, incFlags uint, userID uint64, rr ...FlaggedResource) error {
	for _, r := range rr {

		var (
			flags []string
			err   error
		)

		if incFlags == 0 {
			flags, err = loadCalculated(ctx, s, userID, r)
		} else if incFlags == 1 {
			flags, err = loadGlobal(ctx, s, userID, r)
		} else if incFlags == 2 {
			flags, err = loadOwn(ctx, s, userID, r)
		} else {
			return fmt.Errorf("unknown flag inclusion: %d", incFlags)
		}
		if err != nil {
			return err
		}

		r.SetFlags(flags)
	}
	return nil
}

func loadCalculated(ctx context.Context, s store.Storer, userID uint64, r FlaggedResource) ([]string, error) {
	var (
		ff  types.FlagSet
		err error

		fMap = make(map[string]bool)
	)

	// Get flags for all users
	ff, _, err = store.SearchFlags(ctx, s, types.FlagFilter{
		Kind:       r.FlagResourceKind(),
		ResourceID: []uint64{r.FlagResourceID()},
		OwnedBy:    []uint64{0},
	})
	if err != nil {
		return nil, err
	}
	for _, f := range ff {
		fMap[f.Name] = f.Active
	}

	// Get flags for the given user & merge with general flags
	ff, _, err = store.SearchFlags(ctx, s, types.FlagFilter{
		Kind:       r.FlagResourceKind(),
		ResourceID: []uint64{r.FlagResourceID()},
		OwnedBy:    []uint64{userID},
	})
	if err != nil {
		return nil, err
	}

	for _, f := range ff {
		fMap[f.Name] = f.Active
	}

	// convert to a slice
	rr := make([]string, 0, len(fMap))
	for k, v := range fMap {
		if v {
			rr = append(rr, k)
		}
	}

	return rr, nil
}

func loadGlobal(ctx context.Context, s store.Storer, userID uint64, r FlaggedResource) ([]string, error) {
	var (
		ff  types.FlagSet
		err error

		fMap = make(map[string]bool)
	)

	// Get flags for all users
	ff, _, err = store.SearchFlags(ctx, s, types.FlagFilter{
		Kind:       r.FlagResourceKind(),
		ResourceID: []uint64{r.FlagResourceID()},
		OwnedBy:    []uint64{0},
	})
	if err != nil {
		return nil, err
	}
	for _, f := range ff {
		fMap[f.Name] = f.Active
	}

	// convert to a slice
	rr := make([]string, 0, len(fMap))
	for _, f := range ff {
		if f.Active {
			rr = append(rr, f.Name)
		}
	}

	return rr, nil
}

func loadOwn(ctx context.Context, s store.Storer, userID uint64, r FlaggedResource) ([]string, error) {
	var (
		ff  types.FlagSet
		err error

		fMap = make(map[string]bool)
	)

	// Get flags for all users
	ff, _, err = store.SearchFlags(ctx, s, types.FlagFilter{
		Kind:       r.FlagResourceKind(),
		ResourceID: []uint64{r.FlagResourceID()},
		OwnedBy:    []uint64{userID},
	})
	if err != nil {
		return nil, err
	}
	for _, f := range ff {
		fMap[f.Name] = f.Active
	}

	// convert to a slice
	rr := make([]string, 0, len(fMap))
	for _, f := range ff {
		if f.Active {
			rr = append(rr, f.Name)
		}
	}

	return rr, nil
}
