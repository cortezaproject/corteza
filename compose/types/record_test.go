package types

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"
	"time"
)

func TestRecordUnmarshal(t *testing.T) {
	tests := []struct {
		name string

		// clean
		preloaded *Record

		// update with
		scripted *Record

		// final version
		final *Record
	}{
		{
			"first step",
			&Record{
				ID: 42,
				Values: RecordValueSet{
					&RecordValue{Name: "foo", Value: "foo"},
					&RecordValue{Name: "bar", Value: "foo", Updated: true, DeletedAt: &time.Time{}},
					&RecordValue{Name: "baz", Value: "1"},
				},
			},
			&Record{
				ID: 82,
				Values: RecordValueSet{
					&RecordValue{Name: "foo", Value: "foo"},
					&RecordValue{Name: "baz", Value: "1"},
					&RecordValue{Name: "baz", Value: "2"},
					&RecordValue{Name: "baz", Value: "3"},
				},
			},
			&Record{
				ID: 82,
				Values: RecordValueSet{
					&RecordValue{Name: "foo", Value: "foo"},
					&RecordValue{Name: "baz", Value: "1"},
					&RecordValue{Name: "baz", Value: "2"},
					&RecordValue{Name: "baz", Value: "3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if j, err := json.Marshal(tt.scripted); err != nil {
				t.Errorf("failed to marshal record: %v", err)
			} else if err = json.Unmarshal(j, tt.preloaded); err != nil {
				t.Errorf("failed to unmarshal record: %v", err)
			} else {
				spew.Dump(string(j))
			}

			if !reflect.DeepEqual(tt.preloaded, tt.final) {
				t.Errorf("preloaded:\n%v\n\nfinal\n%v", tt.preloaded.Values, tt.final.Values)
			}
		})
	}
}
