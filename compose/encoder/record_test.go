package encoder

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/test"
)

func Test_RecordEncoding(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		ff   []field
		rr   []*types.Record

		flatResult   string
		structResult string
	}{
		{
			name: "covering the basics",
			ff:   MakeFields("recordID", "ownedBy", "createdAt", "deletedAt", "some-foo-field", "foo", "fff"),
			rr: []*types.Record{
				&types.Record{
					ID:        12345,
					OwnedBy:   12345,
					CreatedAt: time.Unix(1504976400, 0),
				},
				&types.Record{
					ID:        54321,
					OwnedBy:   12345,
					CreatedAt: time.Unix(12345, 0),
					Values: []*types.RecordValue{
						{
							Name:  "foo",
							Value: "bar",
						},
						{
							Name:  "fff",
							Value: "1",
						},
						{
							Name:  "fff",
							Value: "2",
						},
					},
				},
			},

			flatResult: `recordID,ownedBy,createdAt,deletedAt,some-foo-field,foo,fff` + "\n" +
				`12345,12345,2017-09-09T17:00:00Z,,,,` + "\n" +
				`54321,12345,1970-01-01T03:25:45Z,,,bar,1` + "\n",

			structResult: `{"createdAt":"2017-09-09T17:00:00Z","deletedAt":null,"ownedBy":12345,"recordID":12345}` + "\n" +
				`{"createdAt":"1970-01-01T03:25:45Z","deletedAt":null,"fff":["1","2"],"foo":"bar","ownedBy":12345,"recordID":54321}` + "\n",
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name+" (csv)", func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})
			csvWriter := csv.NewWriter(buf)

			fenc := NewFlatWriter(csvWriter, true, tt.ff...)
			for _, r := range tt.rr {
				if err := fenc.Record(r); err != nil {
					t.Errorf("unexpected error = %v,", err)
				}
			}

			csvWriter.Flush()
			test.Assert(t,
				buf.String() == tt.flatResult,
				"Unexpected result: \n%s\n%s",
				buf.String(),
				tt.flatResult)
		})

		t.Run(tt.name+" (json)", func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})
			jsonEnc := json.NewEncoder(buf)

			senc := NewStructuredEncoder(jsonEnc, tt.ff...)
			for _, r := range tt.rr {
				if err := senc.Record(r); err != nil {
					t.Errorf("unexpected error = %v,", err)
				}
			}

			test.Assert(t,
				buf.String() == tt.structResult,
				"Unexpected result: \n%s\n%s",
				buf.String(),
				tt.structResult)
		})
	}
}
