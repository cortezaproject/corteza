package revisions

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/crusttech/human/server/pkg/dal"
)

// DiffAny computes a field-level delta old→new by JSON-encoding both resources,
// flattening nested objects/arrays into dot-paths (e.g. "config.dal.ident",
// "fields.0.name"), and reusing Revision.CollectChanges.
//
// skip entries drop whole top-level subtrees (audit stamps, runtime metadata)
// before comparison. Numbers are kept as json.Number so large uint64 ids keep
// their precision through the round-trip.
func DiffAny(newState, oldState any, skip ...string) ([]*Change, error) {
	newFlat, err := flattenResource(newState, skip)
	if err != nil {
		return nil, err
	}

	oldFlat, err := flattenResource(oldState, skip)
	if err != nil {
		return nil, err
	}

	// CollectChanges only walks keys present in the "new" getter; inject keys
	// that exist solely in old as nil so removals register as changes too.
	for key := range oldFlat {
		if _, ok := newFlat[key]; !ok {
			newFlat[key] = nil
		}
	}

	rev := &Revision{}
	if err = rev.CollectChanges(flatMapGetter(newFlat), flatMapGetter(oldFlat)); err != nil {
		return nil, err
	}

	return rev.Changes, nil
}

func flattenResource(v any, skip []string) (map[string]any, error) {
	if v == nil {
		return map[string]any{}, nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()

	var raw map[string]any
	if err = dec.Decode(&raw); err != nil {
		return nil, err
	}

	for _, key := range skip {
		delete(raw, key)
	}

	out := make(map[string]any)
	flatten("", raw, out)
	return out, nil
}

func flatten(prefix string, v any, out map[string]any) {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			flatten(joinKey(prefix, k), val, out)
		}
	case []any:
		for i, val := range t {
			flatten(joinKey(prefix, strconv.Itoa(i)), val, out)
		}
	default:
		out[prefix] = v
	}
}

func joinKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// flatMapGetter adapts a flattened resource map to dal.ValueGetter so the diff
// can reuse Revision.CollectChanges; every key carries exactly one value.
type flatMapGetter map[string]any

var _ dal.ValueGetter = flatMapGetter(nil)

func (m flatMapGetter) CountValues() map[string]uint {
	out := make(map[string]uint, len(m))
	for k := range m {
		out[k] = 1
	}
	return out
}

func (m flatMapGetter) GetValue(key string, _ uint) (any, error) {
	return m[key], nil
}
