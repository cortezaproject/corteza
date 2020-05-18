package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToBulkOperations(t *testing.T) {
	tests := []struct {
		name string
		bb   BulkRecordSet
		size int
	}{
		{
			name: "Return nothing if empty",
			bb:   BulkRecordSet{},
			size: 0,
		},

		{
			name: "Return all sets all records",
			bb: BulkRecordSet{
				&BulkRecord{
					RefField: "f1",
					Set:      RecordSet{&Record{}},
				},
				&BulkRecord{
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
	bb := BulkRecordSet{
		&BulkRecord{
			RefField: "f1",
			Set: RecordSet{&Record{
				ModuleID: 1000,
			}},
		},
		&BulkRecord{
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
	bb := BulkRecordSet{
		&BulkRecord{
			RefField: "f1",
			Set: RecordSet{&Record{
				ID: 1000,
			}},
		},
		&BulkRecord{
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
