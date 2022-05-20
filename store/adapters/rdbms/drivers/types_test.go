package drivers

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeDecoders(t *testing.T) {
	var (
		cc = []struct {
			name string
			inp  any
			typ  Type
			err  error
			ok   bool
			dec  any
		}{
			{
				name: "time -tz -prec",
				inp:  &sql.NullString{String: "06:06:06", Valid: true},
				typ:  TypeWrap(&dal.TypeTime{}),
				ok:   true,
				dec:  "06:06:06",
			},
			{
				name: "time -tz +prec",
				inp:  &sql.NullString{String: "06:06:06.321", Valid: true},
				typ:  TypeWrap(&dal.TypeTime{Precision: 3}),
				ok:   true,
				dec:  "06:06:06.321",
			},
			{
				name: "time +tz +prec",
				inp:  &sql.NullString{String: "06:06:06.321+04:00", Valid: true},
				typ:  TypeWrap(&dal.TypeTime{Precision: 3, Timezone: true}),
				ok:   true,
				dec:  "06:06:06.321+04:00",
			},
			//{
			//	name: "time ?tz ?prec",
			//	inp:  &sql.NullString{String: "06:06:06", Valid: true},
			//	typ:  TypeWrap(&data.TypeTime{Precision: 3, Timezone: true}),
			//	ok:   true,
			//	dec:  "06:06:06.321+04:00",
			//},
		}
	)

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			dec, ok, err := c.typ.Decode(c.inp)
			if c.err == nil {
				req.NoError(err)
			} else {
				req.ErrorIs(err, c.err)
			}

			req.Equal(c.ok, ok)
			req.EqualValues(c.dec, dec)
		})
	}
}
