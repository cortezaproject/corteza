package types

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetTimeRecStructField(t *testing.T) {
	var (
		req = require.New(t)
		r   = &Record{}
		now = time.Now()
	)

	req.NoError(setTimeRecStructField(r, "createdAt", now))
	req.Equal(r.CreatedAt, now)
	req.NoError(setTimeRecStructField(r, "createdAt", nil))
	req.Equal(r.CreatedAt, now)
	req.NoError(setTimeRecStructField(r, "updatedAt", now))
	req.Equal(r.UpdatedAt, &now)
	req.NoError(setTimeRecStructField(r, "deletedAt", now))
	req.Equal(r.DeletedAt, &now)
	req.NoError(setTimeRecStructField(r, "deletedAt", nil))
	req.Nil(r.DeletedAt)

}

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
			}

			if !reflect.DeepEqual(tt.preloaded, tt.final) {
				t.Errorf("preloaded:\n%v\n\nfinal\n%v", tt.preloaded.Values, tt.final.Values)
			}
		})
	}
}

func TestToBulkOperations(t *testing.T) {
	tests := []struct {
		name string
		bb   RecordBulkSet
		size int
	}{
		{
			name: "Return nothing if empty",
			bb:   RecordBulkSet{},
			size: 0,
		},

		{
			name: "Return all sets all records",
			bb: RecordBulkSet{
				&RecordBulk{
					RefField: "f1",
					Set:      RecordSet{&Record{}},
				},
				&RecordBulk{
					RefField: "f2",
					Set:      RecordSet{&Record{}},
				},
			},
			size: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr, err := tt.bb.ToBulkOperations(0, 0)
			if err != nil {
				t.Errorf("unexpected error = %v,", err)
			}

			require.Equal(t,
				tt.size,
				len(rr))
		})
	}
}

func TestToBulkOperationsDefaultModule(t *testing.T) {
	bb := RecordBulkSet{
		&RecordBulk{
			RefField: "f1",
			Set: RecordSet{&Record{
				ModuleID: 1000,
			}},
		},
		&RecordBulk{
			RefField: "f2",
			Set:      RecordSet{&Record{}},
		},
	}

	rr, err := bb.ToBulkOperations(2000, 0)
	if err != nil {
		t.Errorf("unexpected error = %v,", err)
	}

	require.Equal(t,
		uint64(1000),
		rr[0].Record.ModuleID,
	)

	require.Equal(t,
		uint64(2000),
		rr[1].Record.ModuleID,
		"Expected default value of \n%d got \n%d",
		2000,
		rr[1].Record.ModuleID,
	)
}

func TestToBulkOperationsDetermineOperation(t *testing.T) {
	bb := RecordBulkSet{
		&RecordBulk{
			RefField: "f1",
			Set: RecordSet{&Record{
				ID: 1000,
			}},
		},
		&RecordBulk{
			RefField: "f2",
			Set:      RecordSet{&Record{}},
		},
	}

	rr, err := bb.ToBulkOperations(0, 0)
	if err != nil {
		t.Errorf("unexpected error = %v,", err)
	}

	require.Equal(t,
		OperationTypeUpdate,
		rr[0].Operation,
	)

	require.Equal(t,
		OperationTypeCreate,
		rr[1].Operation,
	)
}

func TestRecord_ValueGetter(t *testing.T) {
	var (
		r = &Record{
			module: &Module{
				Fields: []*ModuleField{
					{Name: "svf"},
					{Name: "mvf", Multi: true},
				},
			},
			Values: RecordValueSet{
				{Name: "svf", Value: "foo"},
				{Name: "mvf", Value: "foo0", Place: 0},
				{Name: "mvf", Value: "foo1", Place: 1},
				{Name: "mvf", Value: "foo2", Place: 2},
			},
		}

		tests = []struct {
			name    string
			field   string
			pos     uint
			count   uint
			want    any
			wantErr bool
		}{
			{
				name:    "Get single value",
				field:   "svf",
				want:    "foo",
				wantErr: false,
			},
			{
				name:  "Get first multi value",
				field: "mvf",
				pos:   0,
				count: 1,
				want:  "foo0",
			},
			{
				name:  "Get second multi value",
				field: "mvf",
				pos:   1,
				count: 1,
				want:  "foo1",
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.GetValue(tt.field, tt.pos)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
