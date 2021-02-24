package json

import (
	"fmt"
	"reflect"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
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

	return addMap(n,
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

	return addMap(n,
		"createdBy", us.CreatedBy,
		"updatedBy", us.UpdatedBy,
		"deletedBy", us.DeletedBy,
		"ownedBy", us.OwnedBy,
	)
}

// cleanMap helper removes any empty k:v nodes from the mapping node
//
// The value is empty when the tag is !!null OR when the value and the content are empty
func cleanMap(n mapNode) mapNode {
	m := make(mapNode)

	for k, v := range n {
		if k != "" && v != nil {
			m[k] = v
		}
	}

	return m
}

// makeMap creates a new mapping node based on the provided k, v items
//
// pp is a set of k, v items; where k's lie at i, and v's lie at i+1.
// non-string values (required by YAML nodes) are processed further.
// eg.: ["k1", "v1", "k2", "v2"]
func makeMap(pp ...interface{}) (mapNode, error) {
	return addMap(make(mapNode), pp...)
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

// addMap adds a new item to the provided mapping node
//
// pp is a set of k, v items; where k's lie at i, and v's lie at i+1.
// non-string values (required by YAML nodes) are processed further.
// eg.: ["k1", "v1", "k2", "v2"]
func addMap(n mapNode, pp ...interface{}) (mapNode, error) {
	if len(pp) == 0 {
		return n, nil
	}

	if len(pp)%2 == 1 {
		return nil, fmt.Errorf("uneven number of elements provided (%d): %v", len(pp), pp)
	}

	for i := 0; i < len(pp); i += 2 {
		kRaw := pp[i]
		vRaw := pp[i+1]

		if isNil(vRaw) {
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

		n[k] = vRaw
	}
	return n, nil
}

// makeSeq creates a new sequence node based on the provided items
func makeSeq(vv ...interface{}) (seqNode, error) {
	return addSeq(make(seqNode, 0, len(vv)), vv...)
}

// addSeq adds new items to the sequence node
func addSeq(n seqNode, vv ...interface{}) (seqNode, error) {
	for _, vRaw := range vv {
		if vRaw != "" && vRaw != nil {
			n = append(n, vRaw)
		}
	}

	return n, nil
}
