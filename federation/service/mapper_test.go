package service

import (
	"encoding/json"
	"testing"

	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/stretchr/testify/require"
)

func TestMapper_merge(t *testing.T) {
	var (
		// out = &
		tcc = []struct {
			name   string
			m      string
			in     ct.RecordValueSet
			out    ct.RecordValueSet
			expect interface{}
		}{
			{
				"merge_missing_destination_field",
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				ct.RecordValueSet{&ct.RecordValue{Name: "Facebook", Value: "https://fb.com/user_1"}, &ct.RecordValue{Name: "Phone", Value: "000 111 222"}, &ct.RecordValue{Name: "Twitter", Value: "https://twitter.com/@russian_bot_0"}},
				ct.RecordValueSet{&ct.RecordValue{Name: "Name", Value: ""}, &ct.RecordValue{Name: "Fb", Value: ""}},
				ct.RecordValueSet{&ct.RecordValue{Name: "Name", Value: ""}, &ct.RecordValue{Name: "Fb", Value: "https://fb.com/user_1"}},
			},
			{
				"merge_empty_origin",
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				ct.RecordValueSet{},
				ct.RecordValueSet{&ct.RecordValue{Name: "Name", Value: ""}, &ct.RecordValue{Name: "Fb", Value: ""}},
				ct.RecordValueSet{&ct.RecordValue{Name: "Name", Value: ""}, &ct.RecordValue{Name: "Fb", Value: ""}},
			},
			{
				"merge_empty_destination",
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				ct.RecordValueSet{&ct.RecordValue{Name: "Facebook", Value: "https://fb.com/user_1"}, &ct.RecordValue{Name: "Phone", Value: "000 111 222"}, &ct.RecordValue{Name: "Twitter", Value: "https://twitter.com/@russian_bot_0"}},
				ct.RecordValueSet{},
				ct.RecordValueSet{},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				mapper = &Mapper{}
			)

			// dont catch any helper errors
			mm := &types.ModuleFieldMappingSet{}
			json.Unmarshal([]byte(tc.m), mm)

			mapper.Merge(&tc.in, &tc.out, mm)
			req.Equal(tc.out, tc.expect)
		})
	}
}
