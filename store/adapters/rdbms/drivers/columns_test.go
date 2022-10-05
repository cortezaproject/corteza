package drivers

import (
	"database/sql"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/stretchr/testify/require"
	"testing"
)

// test deep ident expression generator
func Test_SimpleJsonDocColumn_Decode(t *testing.T) {
	var (
		cc = []struct {
			name string
			json string
			rvs  types.RecordValueSet
			attr []*dal.Attribute
		}{
			{
				name: "empty object",
				json: `{}`,
				rvs:  nil,
			},
			{
				name: "empty string",
				json: ``,
				rvs:  nil,
			},
			{
				name: "mixed",
				json: `{"foo":["bar"],"num":[11,"22"],"bool":[true, false]}`,
				rvs: types.RecordValueSet{
					{RecordID: 0, Name: "foo", Value: "bar"},
					{RecordID: 0, Name: "num", Value: "11"},
					{RecordID: 0, Name: "num", Value: "22", Place: 1},
					{RecordID: 0, Name: "bool", Value: "1"},
					{RecordID: 0, Name: "bool", Value: "", Place: 1},
				},

				attr: []*dal.Attribute{
					{Ident: "foo"},
					{Ident: "num", MultiValue: true},
					{Ident: "bool", MultiValue: true, Type: &dal.TypeBoolean{}},
				},
			},
		}
	)

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			var (
				sjdc = &SimpleJsonDocColumn{attributes: c.attr}
				raw  = sql.RawBytes(c.json)
				rec  = &types.Record{}
			)

			req.NoError(sjdc.Decode(&raw, rec))
			req.Equal(c.rvs.String(), rec.Values.String())
		})
	}
}
