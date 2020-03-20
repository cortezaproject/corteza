package event

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecordMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &recordBase{
			module:    &types.Module{Handle: "mh1"},
			namespace: &types.Namespace{Slug: "slg1"},
		}

		cMod = eventbus.MustMakeConstraint("module", "eq", "mh1")
		cNms = eventbus.MustMakeConstraint("namespace", "eq", "slg1")
	)

	a.True(res.Match(cMod))
	a.True(res.Match(cNms))
}

func TestRecordMatchValues(t *testing.T) {
	var (
		rec = &types.Record{
			Values: types.RecordValueSet{
				&types.RecordValue{Name: "fld1", Value: "val1"},
				&types.RecordValue{Name: "fld2", Value: "val2"},
			},
		}

		res = &recordBase{record: rec}

		cases = []struct {
			match bool
			name  string
			op    string
			v     string
		}{
			{true, "fld1", "eq", "val1"},
			{false, "fld2", "eq", "val1"},
			{false, "fld1", "!=", "val1"},
		}
	)

	for _, c := range cases {
		t.Run(
			fmt.Sprintf("(%s %s %s) == %v", c.name, c.op, c.v, c.match),
			func(t *testing.T) {
				assert.Equal(t, c.match, res.Match(eventbus.MustMakeConstraint("record.values."+c.name, c.op, c.v)))
			},
		)
	}
}
