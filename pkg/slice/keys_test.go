package slice

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {
	var (
		cases = []struct {
			name string
			in   interface{}
			out  string
		}{
			{
				in: map[string]interface{}{"b": nil, "c": nil, "d": nil, "a": nil},
				// warning! false booleans will not be set and that is ok!
				out: "a,b,c,d",
			},
			{
				in: map[int]interface{}{4: nil, 2: nil, 3: nil, 1: nil},
				// warning! false booleans will not be set and that is ok!
				out: "1,2,3,4",
			},
			{
				in:  nil,
				out: "",
			},
			{
				in:  42,
				out: "",
			},
		}
	)

	for _, c := range cases {
		if c.name == "" && c.in != nil {
			c.name = fmt.Sprintf("%T", c.in)
		}

		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			req.EqualValues(c.out, strings.Join(Keys(c.in), ","))
		})
	}
}
