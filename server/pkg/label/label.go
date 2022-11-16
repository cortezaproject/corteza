package label

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/store"
	"strings"
)

type (
	Labels map[string]string

	LabeledResource interface {
		GetLabels() map[string]string
		SetLabel(key string, value string)
		LabelResourceKind() string
		LabelResourceID() uint64
	}
)

// Changed checks if label maps are same or different
func Changed(old, new map[string]string) bool {
	for k := range old {
		if _, has := new[k]; !has {
			return true
		} else if new[k] != old[k] {
			return true
		}
	}

	for k := range new {
		if _, has := old[k]; !has {
			return true
		} else if new[k] != old[k] {
			return true
		}
	}

	return false
}

// ParseStrings converts slice of strings with "key=val" format into
func ParseStrings(ss []string) (m map[string]string, err error) {
	if len(ss) == 0 {
		return nil, nil
	}

	m = make(map[string]string)

	for _, s := range ss {
		if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
			// assume json
			if err = json.Unmarshal([]byte(s), &m); err != nil {
				return nil, err
			}

			continue
		}

		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid label format")
		}

		if !handle.IsValid(kv[0]) {
			return nil, fmt.Errorf("invalid label key format")
		}

		m[kv[0]] = kv[1]
	}

	return m, nil
}

// Search queries all matching (by kind and key-value filter) labels
//
// In case when list of (base) resources is given, labels are also filtered by resource IDs
// to ensure only matching subset is returned
//
// 3 scenarios:
// - empty filter
// - filter set
// - filter & base set
func Search(ctx context.Context, s store.Labels, kind string, f map[string]string, base ...uint64) ([]uint64, error) {
	// label filter not set,
	// return base resource IDs as-is
	if len(f) == 0 {
		return base, nil
	}

	// search for filters
	set, _, err := store.SearchLabels(ctx, s, types.LabelFilter{Kind: kind, Filter: f, ResourceID: base})
	if err != nil {
		return nil, err
	}

	// If we have slice with base IDs, calculate intersection between it and fetched resourceIDs
	// from the labels to ensure we only return results that satisfy BOTH conditions
	return set.ResourceIDs(), nil
}

// Load searches labels for all labeled resources
func Load(ctx context.Context, s store.Labels, rr ...LabeledResource) error {
	if len(rr) == 0 {
		return nil
	}

	var (
		f = types.LabelFilter{ResourceID: make([]uint64, 0, len(rr))}
	)

	for _, r := range rr {
		if f.Kind == "" {
			f.Kind = r.LabelResourceKind()
		} else if f.Kind != r.LabelResourceKind() {
			return fmt.Errorf("expecting one label type, got two: %q, %q", f.Kind, r.LabelResourceKind())
		}

		f.ResourceID = append(f.ResourceID, r.LabelResourceID())
	}

	set, _, err := store.SearchLabels(ctx, s, f)
	if err != nil {
		return err
	}

	for _, r := range rr {
		for k, v := range set.FilterByResource(r.LabelResourceKind(), r.LabelResourceID()) {
			r.SetLabel(k, v)
		}
	}

	return nil
}

// Update updates/creates all labels on labeled resource and removes that are explicitly passed
func Create(ctx context.Context, s store.Labels, r LabeledResource) error {
	var (
		err error
		l   = &types.Label{
			Kind:       r.LabelResourceKind(),
			ResourceID: r.LabelResourceID(),
		}
	)

	for l.Name, l.Value = range r.GetLabels() {
		if err = store.CreateLabel(ctx, s, l); err != nil {
			return err
		}
	}

	return nil
}

// Update updates or creates all labels on labeled resource and removes all non explicitly defined
func Update(ctx context.Context, s store.Labels, r LabeledResource) error {
	var (
		err    error
		labels = r.GetLabels()
		keys   = make([]string, 0, len(labels))
		key    string

		l = &types.Label{
			Kind:       r.LabelResourceKind(),
			ResourceID: r.LabelResourceID(),
		}
	)

	for key = range labels {
		keys = append(keys, key)
	}

	if err = store.DeleteExtraLabels(ctx, s, r.LabelResourceKind(), r.LabelResourceID(), keys...); err != nil {
		return err
	}

	for l.Name, l.Value = range r.GetLabels() {
		if err = store.UpsertLabel(ctx, s, l); err != nil {
			return err
		}
	}

	return nil
}
