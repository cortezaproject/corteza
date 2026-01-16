package label

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/pkg/str"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	Labels map[string]types.LabelValue

	LabeledResource interface {
		GetLabels() map[string]types.LabelValue
		SetLabel(key string, value types.LabelValue)
		LabelResourceKind() string
		LabelResourceID() uint64
	}
)

// Changed checks if label maps are same or different
func Changed(old, new map[string]types.LabelValue) bool {
	for k := range old {
		if _, has := new[k]; !has {
			return true
		} else if !new[k].Equal(old[k]) {
			return true
		}
	}

	for k := range new {
		if _, has := old[k]; !has {
			return true
		} else if !new[k].Equal(old[k]) {
			return true
		}
	}

	return false
}

// ParseStrings converts slice of strings with "key=val" format into
func ParseStrings(ss []string) (m map[string]types.LabelValue, err error) {
	h, err := str.ParseStrings(ss)
	if err != nil {
		return nil, err
	}

	m = make(map[string]types.LabelValue, len(h))
	for k, v := range h {
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			var arr []string
			if err := json.Unmarshal([]byte(v), &arr); err == nil {
				m[k] = types.LabelValue{Values: arr}
				continue
			}
		}
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") && len(v) > 1 {
			v = v[1 : len(v)-1]
		}
		m[k] = types.LabelValue{Val: v}
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
func Search(ctx context.Context, s store.Labels, kind string, f map[string]types.LabelValue, base ...uint64) ([]uint64, error) {
	// label filter not set,
	// return base resource IDs as-is
	if len(f) == 0 {
		return base, nil
	}
	var result []uint64 = base
	// search for filters
	for k, v := range f {
		var values []string
		if v.Val != "" {
			values = []string{v.Val}
		} else if len(v.Values) > 0 {
			values = v.Values
		} else {
			continue
		}

		set, _, err := store.SearchLabels(ctx, s, types.LabelFilter{Kind: kind, Filter: map[string][]string{k: values}, ResourceID: result})
		if err != nil {
			return nil, err
		}

		// If we have slice with base IDs, calculate intersection between it and fetched resourceIDs
		// from the labels to ensure we only return results that satisfy BOTH conditions
		result = set.ResourceIDs()

		if len(result) == 0 {
			return []uint64{}, nil
		}
	}

	return result, nil
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
