package j7s

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

type (
	seqNode []interface{}
	mapNode map[string]interface{}
)

// mapTimestamps helper encodes Timestamps into the mapping node
func mapTimestamps(n mapNode, ts *resource.Timestamps) (mapNode, error) {
	if ts == nil {
		return n, nil
	}

	return AddMap(n,
		"createdAt", ts.CreatedAt,
		"updatedAt", ts.UpdatedAt,
		"deletedAt", ts.DeletedAt,
		"archivedAt", ts.ArchivedAt,
		"suspendedAt", ts.SuspendedAt,
	)
}

// mapUserstamps helper encodes Userstamps into the mapping node
func mapUserstamps(n mapNode, us *resource.Userstamps) (mapNode, error) {
	if us == nil {
		return n, nil
	}

	return AddMap(n,
		"createdBy", us.CreatedBy,
		"updatedBy", us.UpdatedBy,
		"deletedBy", us.DeletedBy,
		"ownedBy", us.OwnedBy,
		"runAs", us.RunAs,
	)
}

// CleanMap helper removes any empty k:v nodes from the mapping node
//
// The value is empty when the tag is !!null OR when the value and the content are empty
func CleanMap(n mapNode) mapNode {
	m := make(mapNode)

	for k, v := range n {
		if k != "" && v != nil {
			m[k] = v
		}
	}

	return m
}

// MakeMap creates a new mapping node based on the provided k, v items
//
// pp is a set of k, v items; where k's lie at i, and v's lie at i+1.
// non-string values (required by YAML nodes) are processed further.
// eg.: ["k1", "v1", "k2", "v2"]
func MakeMap(pp ...interface{}) (mapNode, error) {
	return AddMap(make(mapNode), pp...)
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

// AddMap adds a new item to the provided mapping node
//
// pp is a set of k, v items; where k's lie at i, and v's lie at i+1.
// non-string values (required by YAML nodes) are processed further.
// eg.: ["k1", "v1", "k2", "v2"]
func AddMap(n mapNode, pp ...interface{}) (mapNode, error) {
	if len(pp) == 0 {
		return n, nil
	}

	if len(pp)%2 == 1 {
		return nil, fmt.Errorf("uneven number of elements provided (%d): %v", len(pp), pp)
	}

	for i := 0; i < len(pp); i += 2 {
		kRaw := pp[i]
		vRaw := pp[i+1]

		if IsNil(vRaw) {
			continue
		}

		k, ok := kRaw.(string)
		if !ok {
			return nil, fmt.Errorf("keys must be of type string: %v", kRaw)
		}

		switch v := vRaw.(type) {
		case string:
			if v == "" {
				continue
			}
		case bool:
			if !v {
				continue
			}
		case int, uint,
			int32, uint32,
			int64, uint64:
			if v == 0 {
				continue
			}
		}

		if err := setNestedValue(n, k, vRaw); err != nil {
			return nil, err
		}

	}
	return n, nil
}

// GetByPath retrieves a value from a nested map using a simplified JSONPath.
// Example path: "a.b.c"
func GetByPath(root mapNode, path string) (any, bool, error) {
	if path == "" {
		return nil, false, fmt.Errorf("empty path")
	}

	parts := strings.Split(path, ".")
	var current any = root

	for _, part := range parts {
		if part == "" {
			return nil, false, fmt.Errorf("invalid path %q", path)
		}

		m, ok := current.(map[string]any)
		if !ok {
			return nil, false, nil
		}

		val, exists := m[part]
		if !exists {
			return nil, false, nil
		}

		current = val
	}

	return current, true, nil
}

// MakeSeq creates a new sequence node based on the provided items
func MakeSeq(vv ...interface{}) (seqNode, error) {
	return AddSeq(make(seqNode, 0, len(vv)), vv...)
}

// AddSeq adds new items to the sequence node
func AddSeq(n seqNode, vv ...interface{}) (seqNode, error) {
	for _, vRaw := range vv {
		if vRaw != "" && vRaw != nil {
			n = append(n, vRaw)
		}
	}

	return n, nil
}

func setNestedValue(root mapNode, key string, value interface{}) error {
	parts := strings.Split(key, ".")
	current := root

	for i, part := range parts {
		if part == "" {
			return fmt.Errorf("invalid key %q", key)
		}

		isLast := i == len(parts)-1

		if isLast {
			current[part] = value
			return nil
		}

		// Walk or create next level
		next, exists := current[part]
		if !exists {
			child, _ := MakeMap()
			current[part] = child
			current = child
			continue
		}

		// Ensure existing value is a map
		child, ok := next.(map[string]interface{})
		if !ok {
			return fmt.Errorf(
				"key conflict at %q: expected map, found %T",
				part,
				next,
			)
		}

		current = child
	}

	return nil
}
