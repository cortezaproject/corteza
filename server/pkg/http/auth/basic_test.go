package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_basic(t *testing.T) {
	type (
		tf struct {
			name   string
			err    string
			exp    string
			params BasicParams
		}
	)

	var (
		tcc = []tf{
			{
				name:   "match basic headers fail username validation",
				err:    "invalid param username",
				params: BasicParams{},
			},
			{
				name: "match basic headers fail password validation",
				err:  "invalid param password",
				params: BasicParams{
					User: "user",
				},
			},
			{
				name: "match basic headers success",
				params: BasicParams{
					User: "thou",
					Pass: "shallnotpass",
				},
				exp: "dGhvdTpzaGFsbG5vdHBhc3M=",
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				s, err = NewBasic(tc.params)
			)

			if tc.err != "" {
				req.EqualError(err, tc.err)
				return
			}

			req.Equal(tc.exp, s.Do(context.Background()))
		})
	}
}
